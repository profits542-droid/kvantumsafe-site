package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

type DomainCheck struct {
	ID        int64
	Domain    string
	CreatedAt time.Time
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(4)
	db.SetConnMaxLifetime(0)

	store := &Store{db: db}
	if err := store.configure(); err != nil {
		db.Close()
		return nil, err
	}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := store.seedIfEmpty(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) configure() error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA busy_timeout = 5000",
		"PRAGMA journal_mode = WAL",
	}
	for _, q := range pragmas {
		if _, err := s.db.Exec(q); err != nil {
			return fmt.Errorf("%s: %w", q, err)
		}
	}
	return nil
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS sessions (
			token_hash TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			expires_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS domain_checks (
			id INTEGER PRIMARY KEY,
			domain TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	return s.migrateContent()
}

func (s *Store) CountUsers(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&n)
	return n, err
}

func (s *Store) CreateUser(ctx context.Context, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)`,
		username, string(hash), time.Now().UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (s *Store) SetPassword(ctx context.Context, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	res, err := s.db.ExecContext(ctx, `UPDATE users SET password_hash = ? WHERE username = ?`, string(hash), username)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return s.CreateUser(ctx, username, password)
	}
	return nil
}

func (s *Store) Authenticate(ctx context.Context, username, password string, dummyHash []byte) (int64, error) {
	var id int64
	var hash string
	err := s.db.QueryRowContext(ctx,
		`SELECT id, password_hash FROM users WHERE username = ?`, username,
	).Scan(&id, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
		return 0, errInvalidLogin
	}
	if err != nil {
		return 0, err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return 0, errInvalidLogin
	}
	return id, nil
}

var errInvalidLogin = errors.New("invalid login")

func (s *Store) InvalidLogin(err error) bool {
	return errors.Is(err, errInvalidLogin)
}

func (s *Store) CreateSession(ctx context.Context, userID int64, token string, expires time.Time) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)`,
		hashToken(token), userID, expires.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (s *Store) UsernameBySession(ctx context.Context, token string) (string, error) {
	var username string
	var expiresRaw string
	err := s.db.QueryRowContext(ctx, `
		SELECT users.username, sessions.expires_at
		FROM sessions
		JOIN users ON users.id = sessions.user_id
		WHERE sessions.token_hash = ?`,
		hashToken(token),
	).Scan(&username, &expiresRaw)
	if errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	if err != nil {
		return "", err
	}
	expires, err := time.Parse(time.RFC3339Nano, expiresRaw)
	if err != nil || !expires.After(time.Now()) {
		_, _ = s.db.ExecContext(ctx, `DELETE FROM sessions WHERE token_hash = ?`, hashToken(token))
		return "", sql.ErrNoRows
	}
	return username, nil
}

func (s *Store) DeleteSession(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE token_hash = ?`, hashToken(token))
	return err
}

func (s *Store) AddDomainCheck(ctx context.Context, domain string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO domain_checks (domain, created_at) VALUES (?, ?)`,
		domain, time.Now().UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (s *Store) RecentDomainChecks(ctx context.Context, limit int) ([]DomainCheck, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, domain, created_at FROM domain_checks ORDER BY id DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []DomainCheck
	for rows.Next() {
		var item DomainCheck
		var raw string
		if err := rows.Scan(&item.ID, &item.Domain, &raw); err != nil {
			return nil, err
		}
		item.CreatedAt, err = time.Parse(time.RFC3339Nano, raw)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
