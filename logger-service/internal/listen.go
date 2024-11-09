package internal

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

func (c *Config) RpcListen() error {
	port := os.Getenv("RPC_PORT")
	log.Println("Запуск сервера RPC на порту ", port)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
