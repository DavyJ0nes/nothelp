package templates

import (
	"bufio"
	"bytes"
	"embed"
	"os"
	"strings"
	"text/template"
)

//go:embed template.md
var tmplFS embed.FS

func Parse(date string) ([]byte, error) {
	tmpl, err := template.ParseFS(tmplFS, "template.md")
	if err != nil {
		return nil, err
	}

	data := struct {
		Date string
	}{
		Date: date,
	}

	out := bytes.NewBuffer(nil)
	if err := tmpl.Execute(out, data); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func GetLineNumber(filePath, text string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			return lineNumber, nil
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return -1, err
	}

	return -1, nil // Return 0 if the text is not found

}
