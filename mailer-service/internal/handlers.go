package internal

import "net/http"

func (c *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage
	if err := c.readJSON(w, r, &requestPayload); err != nil {
		c.errorJSON(w, err)
		return
	}
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	if err := c.Mailer.SendSMTPMessage(msg); err != nil {
		c.errorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "Отправлено на " + requestPayload.To,
	}
	c.writeJSON(w, http.StatusAccepted, payload)
}
