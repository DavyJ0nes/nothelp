package templates

import (
	"bufio"
	"bytes"
	"embed"
	"os"
	"strings"
	"text/template"
)

//go:embed daily_template.md weekly_template.md
var tmplFS embed.FS

func Parse(date string) ([]byte, error) {
	tmpl, err := template.New("daily_template.md").Funcs(template.FuncMap{
		"date": func() string { return date },
	}).ParseFS(tmplFS, "daily_template.md")
	if err != nil {
		return nil, err
	}

	out := bytes.NewBuffer(nil)
	if err := tmpl.Execute(out, nil); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func ParseWeekly(week string) ([]byte, error) {
	tmpl, err := template.New("weekly_template.md").Funcs(template.FuncMap{
		"week": func() string { return week },
	}).ParseFS(tmplFS, "weekly_template.md")
	if err != nil {
		return nil, err
	}

	out := bytes.NewBuffer(nil)
	if err := tmpl.Execute(out, nil); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func GetLineNumber(filePath, text string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return -1, err
	}
	defer func() { _ = file.Close() }()

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
