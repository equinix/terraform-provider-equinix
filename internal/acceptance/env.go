package acceptance

import (
	"fmt"
	"os"
)

func GetFromEnvDefault(varName string, defaultValue string) string {
	if v := os.Getenv(varName); v != "" {
		return v
	}
	return defaultValue
}

func GetFromEnv(varName string) (string, error) {
	if v := os.Getenv(varName); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("environmental variable '%s' is not set", varName)
}
