Axcessms ![Tests](https://github.com/mrzen/axcessms/workflows/Go/badge.svg)
=============

API Client library for [Axcessms](https://axcessms.com/) payment gateway.

Usage
-----

````go
import "fmt"
import "github.com/mrzen/axcessms"

func main() {
    client := accessms.New(axcessms.EnvironmentTokenProvider("AXCESSMS_TOKEN"))

    checkout, err := client.Checkout(&axcessms.CreateCheckoutRequest{
        EntityId: EntityId,
        Currency: "USD",
        Amount: 6969,
        PaymentType: "DB"
    })

    if err != nil {
        fmt.Println("Got an error: %s", err)
    } else {
        fmt.Println("Checkout Token: %s", checkout.ID)
    }
}

````