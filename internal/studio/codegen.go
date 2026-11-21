package studio

import (
	"fmt"
	"strings"
)

func Generate(mode, model string) (string, error) {
	switch mode {
	case "ollama":
		return ollamaSource(), nil
	case "xai":
		if model != "grok-3" && model != "grok-beta" {
			return "", fmt.Errorf("unknown model")
		}
		return xaiSource(model), nil
	default:
		return "", fmt.Errorf("unknown mode")
	}
}

func ollamaSource() string {
	return `package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
)

func main() {
	prompt := ""
	if len(os.Args) > 1 {
		prompt = os.Args[1]
	}
	fmtBody, err := askLocal(maskSensitive(prompt))
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	os.Stdout.WriteString(fmtBody)
}

func askLocal(prompt string) (string, error) {
	payload, _ := json.Marshal(map[string]any{
		"model": "llama3.2",
		"messages": []map[string]string{{"role": "user", "content": prompt}},
		"stream": false,
	})
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:11434/api/chat", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

` + maskFunc()
}

func xaiSource(model string) string {
	return fmt.Sprintf(`package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
)

func main() {
	prompt := ""
	if len(os.Args) > 1 {
		prompt = os.Args[1]
	}
	answer, err := askGrok(maskSensitive(prompt))
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	os.Stdout.WriteString(answer)
}

func askGrok(prompt string) (string, error) {
	key := os.Getenv("XAI_API_KEY")
	if key == "" {
		return "", fmtMissingKey()
	}
	payload, _ := json.Marshal(map[string]any{
		"model": %q,
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	})
	req, err := http.NewRequest(http.MethodPost, "https://api.x.ai/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func fmtMissingKey() error { return errString("set XAI_API_KEY in the environment; the key is not written into this file") }

type errString string

func (e errString) Error() string { return string(e) }

`, model) + maskFunc()
}

func maskFunc() string {
	return strings.TrimSpace(`
var (
	labeledName = regexp.MustCompile(`+"`"+`(?i:(клиент|client|имя|name|фио))\s*[:：]\s*[А-ЯЁA-Z][А-Яа-яЁёA-Za-z]+(?:[\s-][А-ЯЁA-Z][А-Яа-яЁёA-Za-z]+){0,2}`+"`"+`)
	labeledDoc  = regexp.MustCompile(`+"`"+`(?i)(паспорт|passport|иин|iin|инн)\s*[:：]?\s*[A-ZА-Я0-9]{6,14}`+"`"+`)
	phone       = regexp.MustCompile(`+"`"+`\+\d[\d\s()\-]{8,}\d`+"`"+`)
	account     = regexp.MustCompile(`+"`"+`\b\d{13,20}\b`+"`"+`)
	iin         = regexp.MustCompile(`+"`"+`\b\d{12}\b`+"`"+`)
)

func maskSensitive(text string) string {
	text = labeledName.ReplaceAllString(text, "$1: [CLIENT_NAME]")
	text = labeledDoc.ReplaceAllString(text, "$1: [ACCOUNT_ID]")
	text = phone.ReplaceAllString(text, "[PHONE]")
	text = account.ReplaceAllString(text, "[ACCOUNT_ID]")
	text = iin.ReplaceAllString(text, "[ACCOUNT_ID]")
	return text
}
`) + "\n"
}
