package infobotdb

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

var QUERY_LIMIT = 5

// генерация секретного ключа указанной длины
func GenerateSecretKey(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func EscapeMarkdownV2(text string) string {
	// Экранирование специальных символов для MarkdownV2
	replacements := map[string]string{
		"_": "\\_",
		"*": "\\*",
		"[": "\\[",
		"]": "\\]",
		"(": "\\(",
		")": "\\)",
		"~": "\\~",
		"`": "\\`",
		">": "\\>",
		"#": "\\#",
		"+": "\\+",
		"-": "\\-",
		"=": "\\=",
		"|": "\\|",
		"{": "\\{",
		"}": "\\}",
		".": "\\.",
		"!": "\\!",
	}

	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}

	return text
}

// генерация sql скрипта
func Template(name, sqlt string, data *OptionsInfoBot) (string, error) {

	var re = regexp.MustCompile(`[ ]{2,}|[\t\n]+`)
	var sqlBuf bytes.Buffer

	tmp, err := template.New(name).Funcs(template.FuncMap{
		"isnnil": isnnil,
	}).Parse(sqlt)
	if err != nil {
		return "", err
	}
	err = tmp.Execute(&sqlBuf, data)
	if err != nil {
		return "", err
	}
	s := re.ReplaceAllString(sqlBuf.String(), ` `)
	return s, nil
}

// проверка аргумента на nil
func isnnil(obj ...any) bool {
	for _, c := range obj {
		if !(c == nil || (reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil())) {
			return true
		}
	}
	return false
}
