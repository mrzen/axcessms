package axcessms

import (
	"context"
	"fmt"
	"net/http"
)

// PaymentType is an enum representing the different types of supported payments
type PaymentType string

// TransactionCategory is an enum representing the different types of transaction
type TransactionCategory string

const (
	// PaymentTypePreAuthorization is the payment type to pre-authorize a future payment
	PaymentTypePreAuthorization PaymentType = "PA"

	// PaymentTypeDebit is the payment type to debit funds form an account or card
	PaymentTypeDebit PaymentType = "DB"

	// PaymentTypeCredit is the payment type to credit fund to an account or card
	PaymentTypeCredit PaymentType = "CD"

	// PaymentTypeCapture is the payment type to capture a previously pre-authorized payment
	PaymentTypeCapture PaymentType = "CP"

	// PaymentTypeReversal is the payment type to recverse a previous pre-authorization, debit or credit before the transaction is settled
	PaymentTypeReversal PaymentType = "RV"

	// PaymentTypeRefund is the payment type to credit an account based on a previous account debit or credit.
	PaymentTypeRefund PaymentType = "RF"

	// TransactionCategoryECommerce is the transaction category for an e-commerce transaction
	TransactionCategoryECommerce TransactionCategory = "EC"

	// TransactionCategoryMailOrder is the transaction category for a mail-order transaction
	TransactionCategoryMailOrder TransactionCategory = "MO"

	// TransactionCategoryTelephoneOrder is the transaction category for a telephone order transaction
	TransactionCategoryTelephoneOrder TransactionCategory = "TO"

	// TransactionCategoryRecurring is the transaction category for a recurring subscription payment
	TransactionCategoryRecurring TransactionCategory = "RC"

	// TransactionCategoryInstallment is the transaction category for an installment payment
	TransactionCategoryInstallment TransactionCategory = "IN"

	// TransactionCategoryPointOfSale is the transaction category for a payment made at a fixed Point-of-Sale
	TransactionCategoryPointOfSale TransactionCategory = "PO"

	// TransactionCategoryMobilePointOfSale is the transaction category for a payment made at a mobile point-of-sale
	TransactionCategoryMobilePointOfSale TransactionCategory = "PM"
)

// CreateCheckoutRequest represents the parameters for creating a checkout
type CreateCheckoutRequest struct {
	EntityID            string              `schema:"entityId"`
	Amount              int                 `schema:"amount"`
	Currency            string              `schema:"currency"`
	PaymentType         PaymentType         `schema:"paymentType"`
	TransactionCategory TransactionCategory `schema:"transactionCategory"`
	Customer            *Customer           `schema:"customer,omitempty"`
	BillingAddress      *BillingAddress     `schema:"billing,omitempty"`
}

// Checkout represents the data for a checkout creation request
type Checkout struct {
	APIResponse
	ID       string `json:"id"`
	EntityID string `json:"-"` // Not set by the API, added by client.
}

// Payment represents the data for a checkout show request
type Payment struct {
	APIResponse

	ID           string      `json:"id"`
	PaymentType  PaymentType `json:"paymentType"`
	PaymentBrand string      `json:"paymentBrand"`
	Amount       float64     `json:"amount,string"`
	Currency     string      `json:"currency"`
	Descriptor   string      `json:"descriptor"`

	Risk struct {
		Score uint `json:"score,string"`
	} `json:"risk"`

	Custom map[string]string `json:"customParameters"`
}

// Customer represents data known about the customer
type Customer struct {
	ID          string     `schema:"customer.merchantCustomerId" json:"id"`
	FirstName   string     `schema:"customer.givenName" json:"first_name"`
	MiddleNames string     `schema:"customer.middleName" json:"middle_names"`
	LastName    string     `schema:"customer.surname" json:"last_name"`
	Browser     *UserAgent `schema:"customer.browser" json:"browser"`
}

// UserAgent represents confuration and identification information about the customer's user agent
type UserAgent struct {
	Accept        string `schema:"customer.browser.acceptHeader" json:"accept"`
	Language      string `schema:"customer.browser.language" json:"language"`
	ScreenHeight  uint   `schema:"customer.browser.screenHeight" json:"screen_height"`
	ScreenWidth   uint   `schema:"customer.browser.screenWidth" json:"screen_width"`
	UTCOffset     int    `schema:"customer.browser.timezone" json:"utc_offset"`
	UserAgent     string `schema:"customer.browser.userAgent" json:"user_agent"`
	IsJavaEnabled bool   `schema:"customer.browser.javaEnabled" json:"java_enabled"`
}

// BillingAddress represents the content of the customer's billing address. Used for AVS.
type BillingAddress struct {
	Street1      string  `schema:"billing.street1" json:"line_1"`
	Street2      *string `schema:"billing.street2,omitempty" json:"line_2"`
	HouseNumber1 *string `schema:"billing.houseNumber1,omitempty" json:"house_number_1"`
	HouseNumber2 *string `schema:"billing.houseNumber2,omitempty" json:"house_number_2"`
	City         string  `schema:"billing.city" json:"city"`
	State        *string `schema:"billing.state,omitempty" json:"state"`
	PostalCode   string  `schema:"billing.postcode,omitempty" json:"postal_code"`
	CountryCode  string  `schema:"billing.country" json:"country"`
}

// CreateCheckout creates a new Checkout
func (c Client) CreateCheckout(ctx context.Context, checkout *CreateCheckoutRequest) (*Checkout, error) {
	response := &Checkout{}
	err := c.Run(ctx, http.MethodPost, "/v1/checkouts", checkout, response)

	if err != nil {
		return nil, err
	}

	response.EntityID = checkout.EntityID

	return response, err
}

// GetCheckout gets details of a checkout by its ID
func (c Client) GetCheckout(ctx context.Context, checkout Checkout) (*Payment, error) {

	response := &Payment{}

	err := c.Run(ctx, http.MethodGet, fmt.Sprintf("/v1/checkouts/%s/payment?entityId=%s", checkout.ID, checkout.EntityID), nil, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
