package nats

import (
	"github.com/nats-io/nats.go"
)

type Client struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewClient(url string) (*Client, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		js:   js,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
