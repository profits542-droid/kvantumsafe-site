package studio

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMaskHidesClientData(t *testing.T) {
	in := "Клиент: Айгуль Нурланова, телефон +996555123456, счёт 1234567890123456, ИИН 10203198912345, паспорт AB1234567"
	out := Mask(in)
	for _, secret := range []string{"Айгуль", "996555123456", "1234567890123456", "10203198912345", "AB1234567"} {
		if contains(out, secret) {
			t.Fatalf("secret %q remained in %q", secret, out)
		}
	}
	plain := Mask("имя стороны платежа нужно проверить")
	if contains(plain, "[CLIENT_NAME]") {
		t.Fatalf("ordinary phrase was masked: %q", plain)
	}
	for _, tag := range []string{"[CLIENT_NAME]", "[PHONE]", "[ACCOUNT_ID]"} {
		if !contains(out, tag) {
			t.Fatalf("tag %q missing in %q", tag, out)
		}
	}
}

func TestGeneratedGoBuilds(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go is not on PATH")
	}
	for _, mode := range []string{"ollama", "xai"} {
		code, err := Generate(mode, "grok-3")
		if err != nil {
			t.Fatal(err)
		}
		dir := t.TempDir()
		path := filepath.Join(dir, "main.go")
		if err := os.WriteFile(path, []byte(code), 0o600); err != nil {
			t.Fatal(err)
		}
		initCmd := exec.Command("go", "mod", "init", "block")
		initCmd.Dir = dir
		if out, err := initCmd.CombinedOutput(); err != nil {
			t.Fatalf("mod init: %v\n%s", err, out)
		}
		cmd := exec.Command("go", "build", "-o", filepath.Join(dir, "block.exe"), ".")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s: %v\n%s", mode, err, out)
		}
	}
}

func contains(s, part string) bool {
	return len(s) >= len(part) && (s == part || len(part) == 0 || indexOf(s, part) >= 0)
}

func indexOf(s, part string) int {
	for i := 0; i+len(part) <= len(s); i++ {
		if s[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
