package messages

type (
	Text struct {
		Body string `json:"body"`
	}

	MessageRequest struct {
		RecipientType string `json:"recipient_type"`
		To            string `json:"to"`
		Type          string `json:"type"`
		Text          Text   `json:"text"`
	}
)
