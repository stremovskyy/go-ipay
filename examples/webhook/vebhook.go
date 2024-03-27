package main

import (
	"fmt"
	"time"

	"github.com/stremovskyy/go-ipay/internal/ipay"
	"github.com/stremovskyy/go-ipay/private"
)

func main() {
	xmlData := []byte(private.ValidationXMLSuccess)

	payment, err := ipay.ParsePaymentXML(xmlData)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return
	}

	// For demonstration, print out the payment ID and status
	fmt.Printf("PaymentURL ID: %d, Status: %s, Card Token: %s, Card Is Prepaid: %s\n", payment.ID, payment.Status.String(), payment.CardToken, payment.CardIsPrepaid)

	// To demonstrate timestamp conversion to a readable format
	timestamp := time.Unix(payment.Timestamp, 0)
	fmt.Println("Timestamp:", timestamp)
}
