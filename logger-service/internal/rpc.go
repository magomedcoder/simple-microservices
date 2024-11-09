package internal

import (
	"context"
	"log"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, res *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Ошибка записи в MongoDB:", err)
		return err
	}
	*res = "Обработано через RPC: " + payload.Name
	return nil
}
