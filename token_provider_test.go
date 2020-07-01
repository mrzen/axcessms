package axcessms

import (
	"os"
	"testing"
)

func TestStaticToken(t *testing.T) {
	t.Parallel()

	val := "TEST_TOKEN"

	provider := StaticToken(val)

	if ret, err := provider(); err != nil || ret != val {
		t.Errorf("expected token to be %s but it was %s instead", val, ret)
	}
}

func TestEnvToken(t *testing.T) {
	t.Parallel()

	key := "API_TOKEN"
	val := "TEST_TOKEN"

	provider := EnvironmentTokenProvider(key)

	if tok, err := provider(); err == nil || tok != "" {
		if err == nil {
			t.Error("Expected an error but didn't get one")
		}

		if tok != "" {
			t.Errorf("Expected not to get a token but got %s", tok)
		}
	}

	os.Setenv(key, val)

	if tok, err := provider(); err != nil || tok != val {
		if err != nil {
			t.Errorf("Expected no error but got: %s instead", err)
		}

		if tok != val {
			t.Errorf("Expected %s but got %s instead", val, tok)
		}
	}
}
