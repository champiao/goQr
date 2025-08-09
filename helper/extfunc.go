package helper

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func IsRowEmpty(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}
func GetCellByHeader(row []string, hIndex map[string]int, headerName string) string {
	if idx, ok := hIndex[strings.ToLower(headerName)]; ok {
		if idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
	}
	return ""
}
func BuildMailtoURI(to, subject, body string) string {

	mailto := "mailto:"
	if to != "" {
		mailto += to
	}
	body = RemoverAcentosECaracteresEspeciais(body)
	body = strings.ReplaceAll(body, "\n", "%0A") // Encode newlines for URL
	return mailto + "?subject=" + subject + "&body=" + body

}

func RemoverAcentosECaracteresEspeciais(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)

	// Regex para remover caracteres não alfanuméricos
	reg, _ := regexp.Compile(`[^a-zA-Z0-9\s]+:`)
	result = reg.ReplaceAllString(result, "")

	return result
}
