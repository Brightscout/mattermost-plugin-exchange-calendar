package remote

import "encoding/xml"

type WebhookResponseEnvelope struct {
	XMLName xml.Name
	Body    Body
}

type Body struct {
	XMLName          xml.Name
	SendNotification SendNotification `xml:"SendNotification"`
}

type SendNotification struct {
	XMLName          xml.Name
	ResponseMessages ResponseMessages `xml:"ResponseMessages"`
}

type ResponseMessages struct {
	XMLName                         xml.Name
	SendNotificationResponseMessage SendNotificationResponseMessage `xml:"SendNotificationResponseMessage"`
}

type SendNotificationResponseMessage struct {
	XMLName       xml.Name
	ResponseClass string           `xml:"ResponseClass,attr"`
	ResponseCode  string           `xml:"ResponseCode"`
	Notification  NotificationData `xml:"Notification"`
}

type NotificationData struct {
	XMLName        xml.Name
	SubscriptionID string       `xml:"SubscriptionId"`
	CreatedEvent   CreatedEvent `xml:"CreatedEvent"`
	StatusEvent    *StatusEvent `xml:"StatusEvent"`
}

type CreatedEvent struct {
	XMLName xml.Name
	Item    Item `xml:"ItemId"`
}

type StatusEvent struct {
	XMLName xml.Name
}

type Item struct {
	XMLName xml.Name
	EventID string `xml:"Id,attr"`
}

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
