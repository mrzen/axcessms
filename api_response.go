package axcessms

import (
	"strconv"
	"time"
)

// APIResponse is an embeddable struct of standard API response fields
type APIResponse struct {
	BuildNumber string    `json:"buildNumber"`
	Timestamp   Timestamp `json:"timestamp"`
	NDC         string    `json:"ndc"`

	Result struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"result"`
}

// Timestamp is a time but with a custom JSON codec because life must be made harder.
type Timestamp time.Time

// TimestampFormat is the format the API uses for timestamps
const TimestampFormat = "2006-01-02 15:04:05-0700"

// UnmarshalJSON unmarshals a Timestamp from a JSON value
func (ts *Timestamp) UnmarshalJSON(data []byte) error {

	str, _ := strconv.Unquote(string(data))
	t, err := time.Parse(TimestampFormat, str)

	*ts = Timestamp(t)

	return err
}

// MarshalJSON encodes a Timestamp to a JSON value
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	str := time.Time(ts).Format(TimestampFormat)
	return []byte(strconv.Quote(str)), nil
}
