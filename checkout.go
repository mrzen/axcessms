package axcessms

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// PaymentType is an enum representing the different types of supported payments
type PaymentType string

// TransactionCategory is an enum representing the different types of transaction
type TransactionCategory string

// IntegrationMode is an enum representing the integration modes
type IntegrationMode string

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

	// IntegrationModeInternal represents internal gateway simulations
	IntegrationModeInternal IntegrationMode = "INTERNAL"

	// IntegrationModeExternal represents external production gateways
	IntegrationModeExternal IntegrationMode = "EXTERNAL"
)

// CustomParameters is a key/value map of strings
type CustomParameters map[string]string

// CreateCheckoutRequest represents the parameters for creating a checkout
type CreateCheckoutRequest struct {
	EntityID            string
	MerchantID          string
	Amount              int
	Currency            string
	PaymentType         PaymentType
	TransactionCategory TransactionCategory
	IntegrationMode     *IntegrationMode
	Customer            *Customer
	BillingAddress      *BillingAddress
	CustomParameters    *CustomParameters
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

	Custom CustomParameters `json:"customParameters"`
}

// Customer represents data known about the customer
type Customer struct {
	ID           string     `json:"id"`
	FirstName    string     `json:"first_name"`
	MiddleNames  string     `json:"middle_names"`
	LastName     string     `json:"last_name"`
	EmailAddress string     `json:"email"`
	IPAddress    *string    `json:"ip_address"`
	Browser      *UserAgent `json:"browser"`
}

// UserAgent represents confuration and identification information about the customer's user agent
type UserAgent struct {
	Accept        string `json:"accept"`
	Language      string `json:"language"`
	ScreenHeight  uint   `json:"screen_height"`
	ScreenWidth   uint   `json:"screen_width"`
	UTCOffset     int    `json:"utc_offset"`
	UserAgent     string `json:"user_agent"`
	IsJavaEnabled bool   `json:"java_enabled"`
}

// BillingAddress represents the content of the customer's billing address. Used for AVS.
type BillingAddress struct {
	Street1      string  `json:"line_1"`
	Street2      *string `json:"line_2"`
	HouseNumber1 *string `json:"house_number_1"`
	HouseNumber2 *string `json:"house_number_2"`
	City         string  `json:"city"`
	State        *string `json:"state"`
	PostalCode   string  `json:"postal_code"`
	CountryCode  string  `json:"country"`
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

// URLEncode encodes the request as a url.Values
func (r CreateCheckoutRequest) URLEncode() url.Values {

	values := make(url.Values)

	values.Set("entityId", r.EntityID)
	values.Set("merchantTransactionId", r.MerchantID)
	values.Set("amount", strconv.FormatInt(int64(r.Amount), 10))
	values.Set("currency", r.Currency)
	values.Set("paymentType", string(r.PaymentType))
	values.Set("transactionCategory", string(r.TransactionCategory))

	if r.IntegrationMode != nil {
		values.Set("testMode", string(*r.IntegrationMode))
	}

	if r.Customer != nil {
		values.Set("customer.merchantCustomerId", r.Customer.ID)
		values.Set("customer.givenName", r.Customer.FirstName)
		values.Set("customer.middleName", r.Customer.MiddleNames)
		values.Set("customer.surname", r.Customer.LastName)
		values.Set("customer.email", r.Customer.EmailAddress)

		if r.Customer.IPAddress != nil {
			values.Set("customer.ip", *r.Customer.IPAddress)
		}

		if ua := r.Customer.Browser; ua != nil {
			values.Set("customer.browser.acceptHeader", ua.Accept)
			values.Set("customer.browser.userAgent", ua.UserAgent)
			values.Set("customer.browser.language", ua.Language)
			values.Set("customer.browser.timezone", strconv.FormatInt(int64(ua.UTCOffset), 10))
			values.Set("customer.browser.javaEnabled", strconv.FormatBool(ua.IsJavaEnabled))
			values.Set("customer.browser.screenWidth", strconv.FormatUint(uint64(ua.ScreenWidth), 10))
			values.Set("customer.browser.screenHeight", strconv.FormatUint(uint64(ua.ScreenHeight), 10))
		}
	}

	if b := r.BillingAddress; b != nil {
		values.Set("billing.street1", b.Street1)
		values.Set("billing.city", b.City)
		values.Set("billing.postcode", b.PostalCode)
		values.Set("billing.country", b.CountryCode)

		if b.Street2 != nil {
			values.Set("billing.street2", *b.Street2)
		}

		if b.State != nil {
			values.Set("billing.state", *b.State)
		}

		if b.HouseNumber1 != nil {
			values.Set("billing.houseNumber1", *b.HouseNumber1)
		}

		if b.HouseNumber2 != nil {
			values.Set("billing.houseNumber2", *b.HouseNumber2)
		}

	}

	if r.CustomParameters != nil {
		for key, value := range *r.CustomParameters {
			values.Set(fmt.Sprintf("customParameters[%s]", key), value)
		}
	}

	return values
}
