package internal

import (
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload
	_ = c.readJSON(w, r, &requestPayload)
	event := LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	if err := c.Models.LogEntry.Insert(event); err != nil {
		c.errorJSON(w, err)
		return
	}
	res := jsonResponse{
		Error:   false,
		Message: "Записано",
	}
	c.writeJSON(w, http.StatusAccepted, res)
}
