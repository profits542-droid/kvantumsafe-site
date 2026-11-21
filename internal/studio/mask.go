package studio

import "regexp"

var (
	labeledName = regexp.MustCompile(`(?i:(клиент|client|имя|name|фио))\s*[:：]\s*[А-ЯЁA-Z][А-Яа-яЁёA-Za-z]+(?:[\s-][А-ЯЁA-Z][А-Яа-яЁёA-Za-z]+){0,2}`)
	labeledDoc  = regexp.MustCompile(`(?i)(паспорт|passport|иин|iin|инн)\s*[:：]?\s*[A-ZА-Я0-9]{6,14}`)
	phone       = regexp.MustCompile(`\+\d[\d\s()\-]{8,}\d`)
	account     = regexp.MustCompile(`\b\d{13,20}\b`)
	iin         = regexp.MustCompile(`\b\d{12}\b`)
)

func Mask(text string) string {
	text = labeledName.ReplaceAllString(text, "$1: [CLIENT_NAME]")
	text = labeledDoc.ReplaceAllString(text, "$1: [ACCOUNT_ID]")
	text = phone.ReplaceAllString(text, "[PHONE]")
	text = account.ReplaceAllString(text, "[ACCOUNT_ID]")
	text = iin.ReplaceAllString(text, "[ACCOUNT_ID]")
	return text
}
