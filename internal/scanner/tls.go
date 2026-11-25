package scanner

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

const scanTimeout = 3 * time.Second

// ScanResult is the observed result of a TLS connection to domain:443.
type ScanResult struct {
	Domain      string
	IP          string
	TLSVersion  string
	CipherSuite string
	NotBefore   time.Time
	NotAfter    time.Time
	Issuer      string
	Subject     string
	Error       string
	Logs        []string
}

func (r ScanResult) OK() bool {
	return r.Error == "" && r.TLSVersion != ""
}

func (r *ScanResult) logf(format string, args ...any) {
	r.Logs = append(r.Logs, fmt.Sprintf(format, args...))
}

func tlsVersionName(v uint16) string {
	switch v {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("unknown (0x%04x)", v)
	}
}

func pickIP(ips []net.IP) net.IP {
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip
		}
	}
	return ips[0]
}

// ScanDomainTLS dials the domain on port 443 and records the certificate that the server presents.
func ScanDomainTLS(domain, lang string) ScanResult {
	en := lang == "en"
	result := ScanResult{Domain: domain}
	deadline := time.Now().Add(scanTimeout)

	if en {
		result.logf("[DNS] Resolving %s ...", domain)
	} else {
		result.logf("[DNS] Резолв имени %s ...", domain)
	}
	ips, err := net.LookupIP(domain)
	if err != nil {
		result.Error = err.Error()
		if en {
			result.logf("[DNS] Error: %v", err)
		} else {
			result.logf("[DNS] Ошибка: %v", err)
		}
		return result
	}
	ip := pickIP(ips)
	result.IP = ip.String()
	if en {
		result.logf("[DNS] OK → %s", result.IP)
	} else {
		result.logf("[DNS] OK → %s", result.IP)
	}

	tcpAddr := net.JoinHostPort(result.IP, "443")
	if en {
		result.logf("[TCP] Connecting to %s (timeout %s) ...", tcpAddr, scanTimeout)
	} else {
		result.logf("[TCP] Подключение к %s (таймаут %s) ...", tcpAddr, scanTimeout)
	}
	dialer := &net.Dialer{Timeout: time.Until(deadline)}
	if dialer.Timeout <= 0 {
		if en {
			result.Error = "timeout before TCP connect"
			result.logf("[TCP] Timed out before the connection started")
		} else {
			result.Error = "timeout before TCP connect"
			result.logf("[TCP] Таймаут до начала подключения")
		}
		return result
	}
	rawConn, err := dialer.Dial("tcp", tcpAddr)
	if err != nil {
		result.Error = err.Error()
		if en {
			result.logf("[TCP] Error: %v", err)
		} else {
			result.logf("[TCP] Ошибка: %v", err)
		}
		return result
	}
	defer rawConn.Close()
	if en {
		result.logf("[TCP] Connected")
	} else {
		result.logf("[TCP] Соединение установлено")
	}

	if err := rawConn.SetDeadline(deadline); err != nil {
		result.Error = err.Error()
		if en {
			result.logf("[TCP] Could not set a deadline: %v", err)
		} else {
			result.logf("[TCP] Не удалось выставить deadline: %v", err)
		}
		return result
	}

	if en {
		result.logf("[TLS] Handshake (SNI=%s) ...", domain)
	} else {
		result.logf("[TLS] Handshake (SNI=%s) ...", domain)
	}
	tlsConn := tls.Client(rawConn, &tls.Config{
		ServerName: domain,
		MinVersion: tls.VersionTLS10,
	})
	if err := tlsConn.Handshake(); err != nil {
		result.Error = err.Error()
		if en {
			result.logf("[TLS] Handshake failed: %v", err)
		} else {
			result.logf("[TLS] Handshake не удался: %v", err)
		}
		return result
	}
	state := tlsConn.ConnectionState()
	result.TLSVersion = tlsVersionName(state.Version)
	result.CipherSuite = tls.CipherSuiteName(state.CipherSuite)
	result.logf("[TLS] Handshake OK: %s / %s", result.TLSVersion, result.CipherSuite)

	if len(state.PeerCertificates) == 0 {
		if en {
			result.Error = "peer certificate missing"
			result.logf("[CERT] The server did not present a certificate")
		} else {
			result.Error = "peer certificate missing"
			result.logf("[CERT] Сервер не предъявил сертификат")
		}
		return result
	}
	cert := state.PeerCertificates[0]
	result.NotBefore = cert.NotBefore
	result.NotAfter = cert.NotAfter
	result.Issuer = cert.Issuer.String()
	result.Subject = cert.Subject.CommonName
	if result.Subject == "" {
		result.Subject = cert.Subject.String()
	}
	result.logf("[CERT] Subject=%s", result.Subject)
	if en {
		result.logf("[CERT] Validity: %s — %s",
			cert.NotBefore.UTC().Format(time.RFC3339),
			cert.NotAfter.UTC().Format(time.RFC3339))
	} else {
		result.logf("[CERT] Срок действия: %s — %s",
			cert.NotBefore.UTC().Format(time.RFC3339),
			cert.NotAfter.UTC().Format(time.RFC3339))
	}
	return result
}
