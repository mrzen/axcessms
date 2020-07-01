package axcessms

import "context"

// CreateCheckoutRequest represents the parameters for creating a checkout
type CreateCheckoutRequest struct {
	EntityID    string `schema:"entityId"`
	Amount      int    `schema:"amount"`
	Currency    string `schema:"currency"`
	PaymentType string `schema:"paymentType"`
}

// CreateCheckoutResponse represents the data for a checout creation request
type CreateCheckoutResponse struct {
	APIResponse
	ID string `json:"id"`
}

// CreateCheckout creates a new Checkout
func (c Client) CreateCheckout(ctx context.Context, checkout *CreateCheckoutRequest) (*CreateCheckoutResponse, error) {
	response := &CreateCheckoutResponse{}
	err := c.PostForm(ctx, "/v1/checkouts", checkout, response)

	return response, err
}
