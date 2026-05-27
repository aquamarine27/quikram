package validator

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func Email(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func Password(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

func Name(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if len(name) > 100 {
		return errors.New("name must be at most 100 characters")
	}
	return nil
}

func SubjectTitle(title string) error {
	if len(strings.TrimSpace(title)) == 0 {
		return errors.New("title cannot be empty")
	}
	if len(title) > 200 {
		return errors.New("title must be at most 200 characters")
	}
	return nil
}

func MimeType(mime string) error {
	allowed := map[string]bool{
		"application/pdf":                     true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}
	if !allowed[mime] {
		return errors.New("only PDF and DOCX files are allowed")
	}
	return nil
}

func Difficulty(d string) error {
	allowed := map[string]bool{"easy": true, "medium": true, "hard": true}
	if !allowed[d] {
		return errors.New("difficulty must be easy, medium, or hard")
	}
	return nil
}

func CompressionLevel(level string) error {
	allowed := map[string]bool{"short": true, "medium": true, "detailed": true}
	if !allowed[level] {
		return errors.New("compression level must be short, medium, or detailed")
	}
	return nil
}
