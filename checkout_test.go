package axcessms

import (
	"context"
	"testing"
)

func TestCreateCheckout(t *testing.T) {

	client := getTestClient(t)

	req := &CreateCheckoutRequest{
		EntityID:            "8ac7a4c8710d265c01710d30b6e60023",
		Currency:            "USD",
		Amount:              6969,
		PaymentType:         PaymentTypeDebit,
		TransactionCategory: TransactionCategoryECommerce,
	}

	res, err := client.CreateCheckout(context.TODO(), req)

	if err != nil {
		t.Errorf("Expected not to get an error but got '%s' instead", err)
	}

	t.Logf("Response: %#+v", res)
}

func TestGetCheckout(t *testing.T) {

	client := getTestClient(t)

	// Create a checkout to test with
	req := &CreateCheckoutRequest{
		EntityID:            "8ac7a4c8710d265c01710d30b6e60023",
		Currency:            "USD",
		Amount:              6969,
		PaymentType:         PaymentTypeDebit,
		TransactionCategory: TransactionCategoryECommerce,
		Customer: &Customer{
			ID:          "TEST",
			FirstName:   "John",
			MiddleNames: "James",
			LastName:    "Smith",
			Browser: &UserAgent{
				IsJavaEnabled: false,
				ScreenWidth:   640,
				ScreenHeight:  480,
				Accept:        "application/json",
				UserAgent:     "go/1.14 +testUA",
				Language:      "en-US",
				UTCOffset:     1440,
			},
		},
	}

	cko, err := client.CreateCheckout(context.TODO(), req)

	if err != nil {
		t.Errorf("Failed to create a checkout to test against: %s", err)
		return
	}

	resp, err := client.GetCheckout(context.TODO(), *cko)

	if err != nil {
		t.Errorf("Expected not to get an error but insted got: %s", err)
		return
	}

	t.Logf("Checkout: %#+v\n", resp)
}
