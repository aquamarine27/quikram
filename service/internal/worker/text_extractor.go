package worker

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

func extractText(reader io.Reader, mimeType string) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	if len(data) == 0 {
		return "", fmt.Errorf("file is empty")
	}

	var text string
	switch {
	case mimeType == "application/pdf" || strings.HasPrefix(mimeType, "application/pdf"):
		text, err = extractPDFText(data)
		if err != nil || len(strings.TrimSpace(text)) < 100 {
			raw := extractPDFTextRaw(data)
			if len(strings.TrimSpace(raw)) > len(strings.TrimSpace(text)) {
				text = raw
				err = nil
			}
		}
	case strings.Contains(mimeType, "wordprocessingml.document"):
		text, err = extractDOCXText(data)
	default:
		text = string(data)
	}

	if err != nil {
		return "", err
	}

	text = strings.TrimSpace(text)
	if len(text) < 100 {
		return "", fmt.Errorf("extracted text too short (%d chars) — file may be scanned or image-based", len(text))
	}

	return text, nil
}

func extractDOCXText(data []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("docx: not a valid zip: %w", err)
	}

	var docXML []byte
	for _, f := range reader.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("docx: failed to open document.xml: %w", err)
			}
			defer rc.Close()
			docXML, err = io.ReadAll(rc)
			if err != nil {
				return "", fmt.Errorf("docx: failed to read document.xml: %w", err)
			}
			break
		}
	}

	if docXML == nil {
		return "", fmt.Errorf("docx: word/document.xml not found")
	}

	var textBuf strings.Builder
	decoder := xml.NewDecoder(bytes.NewReader(docXML))
	inT := false
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("docx: xml parse error: %w", err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "t" {
				inT = true
			}
			if t.Name.Local == "tab" {
				textBuf.WriteByte('\t')
			}
		case xml.CharData:
			if inT {
				textBuf.Write(t)
			}
		case xml.EndElement:
			if t.Name.Local == "t" {
				inT = false
			}
			if t.Name.Local == "p" {
				textBuf.WriteByte('\n')
			}
		}
	}

	result := strings.TrimSpace(textBuf.String())
	if result == "" {
		return "", fmt.Errorf("docx: no text found")
	}
	return result, nil
}

func extractPDFText(data []byte) (string, error) {
	r := bytes.NewReader(data)
	reader, err := pdf.NewReader(r, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("pdf parse error: %w", err)
	}

	var buf strings.Builder
	for i := 1; i <= reader.NumPage(); i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		content, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		content = strings.TrimSpace(content)
		if content != "" {
			buf.WriteString(content)
			buf.WriteString("\n\n")
		}
	}

	result := strings.TrimSpace(buf.String())
	if result == "" {
		return "", fmt.Errorf("no text extracted from PDF")
	}

	return result, nil
}

func extractPDFTextRaw(data []byte) string {
	content := string(data)

	var buf strings.Builder
	inText := false
	for i := 0; i < len(content); i++ {
		if i+3 < len(content) && content[i:i+3] == "BT" {
			inText = true
			continue
		}
		if inText && i+2 < len(content) && content[i:i+2] == "ET" {
			inText = false
			continue
		}
		if inText {
			c := rune(content[i])
			if unicode.IsPrint(c) || c == '\n' || c == '\r' || c == '\t' {
				buf.WriteByte(content[i])
			}
		}
	}

	result := strings.TrimSpace(buf.String())
	if len(result) > 100 {
		return result
	}

	content = strings.ReplaceAll(content, "\r", " ")
	content = strings.ReplaceAll(content, "\n", " ")
	var readable strings.Builder
	for _, r := range content {
		if unicode.IsPrint(r) || r == ' ' {
			readable.WriteRune(r)
		}
	}
	result = strings.TrimSpace(readable.String())
	if len(result) > 200 {
		return result
	}

	return ""
}
