package axcessms

import (
	"fmt"
	"os"
)

// TokenProvider is a function interface type which attmpts to get an API token or returns an error
type TokenProvider = func() (string, error)

// EnvironmentTokenProvider gets a token from the environment
func EnvironmentTokenProvider(varname string) TokenProvider {
	return func() (string, error) {
		token := os.Getenv(varname)

		if token == "" {
			return token, fmt.Errorf("environment variable '%s' not set", varname)
		}

		return token, nil
	}
}

// StaticToken accepts a plaintext auth token and returns it back as a provider
func StaticToken(token string) TokenProvider {
	return func() (string, error) {
		return token, nil
	}
}
