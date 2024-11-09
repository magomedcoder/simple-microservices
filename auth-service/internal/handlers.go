package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func (c *Config) Auth(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.readJSON(w, r, &requestPayload); err != nil {
		c.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := c.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		c.errorJSON(w, errors.New("Неверные учетные данные"), http.StatusUnauthorized)
		return
	}
	validPass, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !validPass {
		c.errorJSON(w, errors.New("Неверные учетные данные"), http.StatusUnauthorized)
		return
	}

	if err = c.logRequest("Событие аутентификации", fmt.Sprintf("%s вошел в систему", user.Email)); err != nil {
		c.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("пользователь %s успешно вошел в систему", user.Email),
		Data:    user,
	}
	c.writeJSON(w, http.StatusAccepted, payload)
}

func (c *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	url := os.Getenv("HTTP_LOGGER_SERVICE")
	req, err := http.NewRequest("POST", url+"/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
