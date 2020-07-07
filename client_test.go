package axcessms

import (
	"context"
	"net/http"
	"os"
	"testing"
)

var demoToken = StaticToken("OGFjN2E0Yzg3MTBkMjY1YzAxNzEwZDJjZWQ3MDAwMTV8WnBIalo0R215eg==")

func getTestClient(t *testing.T) *Client {

	t.Helper()

	c := New(context.TODO(), demoToken)
	c.SetTestMode(true)
	c.DebugWriter = os.Stdout

	return c
}

func TestClientDo(t *testing.T) {

	c := getTestClient(t)

	req, _ := http.NewRequest(http.MethodGet, TestHost, nil)

	_, err := c.Do(context.TODO(), req)

	if err != nil {
		t.Error(err)
	}
}
