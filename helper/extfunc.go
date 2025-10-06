package helper

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mozillazg/go-unidecode"
	"golang.org/x/text/encoding/charmap"
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
	// Normaliza para decompor acentos e remove marcas
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)

	// Remove caracteres especiais (mas preserva ':', '-', '_', etc se desejar)
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s\-\_:]+`)
	result = reg.ReplaceAllString(result, "")

	// Remove espaços duplicados
	result = strings.Join(strings.Fields(result), " ")

	return result

}

func SplitLines(content string) []string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	return lines

}

func ToUTF8(s string) string {
	// Primeiro, tenta decodificar como Latin1 (ISO-8859-1)
	reader := transform.NewReader(strings.NewReader(s), charmap.ISO8859_1.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err == nil {
		return string(decoded)
	}

	// Se não der certo, tenta como Windows-1252 (comuns em Excel do Windows)
	reader = transform.NewReader(strings.NewReader(s), charmap.Windows1252.NewDecoder())
	decoded, _ = io.ReadAll(reader)

	return string(decoded)
}

func NormalizeText(s string) string {
	s = CleanString(s)
	s = ToUTF8(s)
	s = DecodeToUTF8(s)
	s = RemoverAcentosECaracteresEspeciais(s)
	return s
}

func DecodeToUTF8(s string) string {
	buf := bytes.Buffer{}
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 {
			// substitui byte inválido por espaço (ou ignora)
			s = s[1:]
			continue
		}
		buf.WriteRune(r)
		s = s[size:]
	}
	return buf.String()
}

func CleanString(s string) string {
	// Usa unidecode, que converte qualquer caractere unicode para o equivalente ASCII
	ascii := unidecode.Unidecode(s)

	// Remove caracteres não alfanuméricos, mantendo espaço
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	clean := reg.ReplaceAllString(ascii, "")

	// Remove espaços duplos e normaliza
	clean = strings.Join(strings.Fields(clean), " ")

	return strings.TrimSpace(clean)

}
