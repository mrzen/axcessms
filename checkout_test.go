package axcessms

import (
	"context"
	"testing"
)

func TestCreateCheckout(t *testing.T) {
	t.Parallel()

	client := getTestClient(t)

	req := &CreateCheckoutRequest{
		EntityID:    "8ac7a4c8710d265c01710d30b6e60023",
		Currency:    "USD",
		Amount:      6969,
		PaymentType: "DB",
	}

	res, err := client.CreateCheckout(context.TODO(), req)

	if err != nil {
		t.Errorf("Expected not to get an error but got '%s' instead", err)
	}

	t.Logf("Response: %#+v", res)
}
