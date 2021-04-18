package helper

import (
	"os"
	"strings"
)

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func GetEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	return env
}
