package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/magomedcoder/simple-microservice/gateway-service/api/pb"
	"github.com/magomedcoder/simple-microservice/gateway-service/internal/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type RPCPayload struct {
	Name string
	Data string
}

func (c *Config) Gateway(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Шлюз доступен",
	}
	_ = c.writeJSON(w, http.StatusOK, payload)
}

func (c *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	if err := c.readJSON(w, r, &requestPayload); err != nil {
		c.errorJSON(w, err)
		return
	}
	switch requestPayload.Action {
	case "auth":
		c.auth(w, requestPayload.Auth)
	case "log":
		c.logEventViaRabbitMQ(w, requestPayload.Log)
	case "mail":
		c.sendMail(w, requestPayload.Mail)
	default:
		c.errorJSON(w, errors.New("Неизвестное действие"))
	}
}

func (c *Config) auth(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	url := os.Getenv("HTTP_AUTH_SERVICE")
	request, err := http.NewRequest("POST", url+"/auth", bytes.NewBuffer(jsonData))
	if err != nil {
		c.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		c.errorJSON(w, err)
		return
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusUnauthorized {
		c.errorJSON(w, errors.New("Неверные учетные данные"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		c.errorJSON(w, errors.New("Ошибка вызова службы аутентификации"))
		return
	}

	var jsonFromService jsonResponse
	if err = json.NewDecoder(res.Body).Decode(&jsonFromService); err != nil {
		c.errorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		c.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Аутентифицирован"
	payload.Data = jsonFromService.Data
	c.writeJSON(w, http.StatusAccepted, payload)
}

func (c *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	url := os.Getenv("HTTP_LOGGER_SERVICE")
	request, err := http.NewRequest("POST", url+"/log", bytes.NewBuffer(jsonData))
	if err != nil {
		c.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		c.errorJSON(w, err)
		return
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		c.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Записано"
	c.writeJSON(w, http.StatusAccepted, payload)
}

func (c *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")
	url := os.Getenv("HTTP_MAILER_SERVICE")
	request, err := http.NewRequest("POST", url+"/send", bytes.NewBuffer(jsonData))
	if err != nil {
		c.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		c.errorJSON(w, err)
		return
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		c.errorJSON(w, errors.New("Ошибка вызова службы отправки почты"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Сообщение отправлено на " + msg.To
	c.writeJSON(w, http.StatusAccepted, payload)
}

func (c *Config) logEventViaRabbitMQ(w http.ResponseWriter, l LogPayload) {
	if err := c.pushToQueue(l.Name, l.Data); err != nil {
		c.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Записано через RabbitMQ"
	c.writeJSON(w, http.StatusAccepted, payload)
}

func (c *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(c.Rabbit)
	if err != nil {
		return err
	}
	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	if err = emitter.Push(string(j), "log.INFO"); err != nil {
		return err
	}

	return nil
}

func (c *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	tcp := os.Getenv("TCP_LOGGER_SERVICE")
	client, err := rpc.Dial("tcp", tcp)
	if err != nil {
		c.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}
	var result string
	if err = client.Call("RPCServer.LogInfo", rpcPayload, &result); err != nil {
		c.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = result
	c.writeJSON(w, http.StatusAccepted, payload)
}

func (c *Config) LogItemViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	if err := c.readJSON(w, r, &requestPayload); err != nil {
		c.errorJSON(w, err)
		return
	}

	tcp := os.Getenv("TCP_LOGGER_SERVICE")
	conn, err := grpc.Dial(tcp, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		c.errorJSON(w, err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	l := pb.NewLogServiceClient(conn)
	_, err = l.WriteLog(ctx, &pb.LogRequest{
		LogEntry: &pb.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})

	if err != nil {
		c.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Записано через gRPC"
	c.writeJSON(w, http.StatusAccepted, payload)
}
