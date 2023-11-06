package connect

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"deskor/chat"
	"fmt"
	"github.com/joho/godotenv"
	"net"
	"os"
)

func Setup() (net.Listener, error) {
	err := godotenv.Load(".env.server")
	if err != nil {
		return nil, err
	}
	port := os.Getenv("PORT")

	cert, err := tls.LoadX509KeyPair("./cert/server.pem", "./cert/server.key")
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile("./cert/ca.crt")
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}
	config.Rand = rand.Reader
	return tls.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port), &config)
}

func Accept(l net.Listener) (chat.Client, error) {
	conn, err := l.Accept()
	if err != nil {
		return chat.Client{}, err
	}

	client := chat.Client{
		Conn:     conn,
		Messages: make(chan chat.Message),
	}

	return client, nil
}
