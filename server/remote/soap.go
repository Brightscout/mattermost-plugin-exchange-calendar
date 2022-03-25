package remote

import "encoding/xml"

type SOAPEnvelope struct {
	XMLName xml.Name
	Body    SOAPBody `xml:",omitempty"`
}

type SOAPBody struct {
	XMLName xml.Name
	Content interface{} `xml:",omitempty"`
}

type SendNotificationResult struct {
	XMLName            xml.Name
	SubscriptionStatus string `xml:"SubscriptionStatus"`
}
