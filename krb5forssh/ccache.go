package krb5forssh

import (
	"fmt"
	"os"
	"strings"
)

type InvalidCCacheTypeErr struct {
	supportedTypes []string
	gotType        string
}

func (i InvalidCCacheTypeErr) Error() string {
	return fmt.Sprintf("only %s credentials cache types are supported, got %s", strings.Join(i.supportedTypes, ","), i.gotType)
}

func GetKrb5CCacheFilename() (string, error) {
	env := os.Getenv("KRB5CCNAME")
	if env == "" {
		return "", fmt.Errorf("empty value for KRB5CCNAME environment variable")
	}

	if !strings.HasPrefix(env, "FILE:") {
		return "", InvalidCCacheTypeErr{
			supportedTypes: []string{"FILE"},
			gotType:        env,
		}
	}

	return env[len("FILE:"):], nil
}
