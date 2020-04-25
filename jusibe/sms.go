package jusibe

type smsDeliveryStatus string

const (
	// StatusSMSRejected is delivery status for rejected SMS
	// This indicates that the SMS wasn't sent
	StatusSMSRejected smsDeliveryStatus = "Rejected"

	// StatusSMSSent is delivery status for send SMS
	// This indicates that the SMS was successfully sent but the receipient is yet to receive the SMS
	StatusSMSSent smsDeliveryStatus = "Sent"

	// StatusSMSDelivered is delivery status for delivered SMS
	// This indicates that the receipient received the SMS
	StatusSMSDelivered smsDeliveryStatus = "Delivered"

	// StatusBulkSMSSubmitted is submitted status for Submitted Bulk SMS
	// This indicates that the bulk sms was successfully submitted to the api
	StatusBulkSMSSubmitted smsDeliveryStatus = "Submitted"
)

// SMSResponse is response returned from Jusibe `send_sms` endpoint
type SMSResponse struct {
	Status         string `json:"status"`
	MessageID      string `json:"message_id"`
	SMSCreditsUsed int    `json:"sms_credits_used"`
}

// SMSDeliveryResponse is response returned from Jusibe `delivery_status` endpoint
type SMSDeliveryResponse struct {
	MessageID     string `json:"message_id"`
	Status        string `json:"status"`
	DateSent      string `json:"date_sent"`
	DateDelivered string `json:"date_delivered"`
}

// SMSCreditsResponse is response returned from Jusibe `get_credits` endpoint
type SMSCreditsResponse struct {
	SMSCredits string `json:"sms_credits"`
}

// BulkSMSResponse is response returned form Jusibe `bulk/send_sms` endpoint
type BulkSMSResponse struct {
	Status    string `json:"status"`
	MessageID string `json:"bulk_message_id"`
}

// BulkSMSStatusResponse is response from Jusibe `bulk/status` endpoint
type BulkSMSStatusResponse struct {
	BulkMessageID string `json:"bulk_message_id"`

	Status    string `json:"status"`
	Created   string `json:"created"`
	Processed string `json:"processed"`

	TotalNumbers string `json:"total_numbers"`

	TotalUniqueNumbers  string `json:"total_unique_numbers"`
	TotalValidNumbers   string `json:"total_valid_numbers"`
	TotalInvalidNumbers string `json:"total_invalid_numbers"`
}
