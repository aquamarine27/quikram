package worker

import (
	"io"
	"strings"
)

func extractText(reader io.Reader, mimeType string) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", nil
	}

	if mimeType == "application/pdf" {
		return extractPDFText(data)
	}

	return extractDOCXText(data)
}

func extractPDFText(data []byte) (string, error) {
	content := string(data)

	start := strings.Index(content, "<<"+"/Type")
	if start > 0 {
		content = content[start:]
	}

	var result strings.Builder
	inText := false
	for i := 0; i < len(content); i++ {
		if i+4 < len(content) && content[i:i+4] == "Tj\n" {
			inText = false
		}
		if inText {
			result.WriteByte(content[i])
		}
	}

	if result.Len() == 0 {
		result.WriteString("(Placeholder: PDF text extraction will be implemented with a proper library)\n")
	}

	return result.String(), nil
}

func extractDOCXText(data []byte) (string, error) {
	return string(data), nil
}
