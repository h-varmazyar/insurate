package amqp

import (
	"github.com/h-varmazyar/insurate/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"sync"
)

var (
	ConnectionClosedBefore = errors.New("connection closed before")
	ChannelClosedBefore    = errors.New("channel closed before")
	ClientCreationFailed   = errors.NewWithHttp("client creation failed", 1, http.StatusServiceUnavailable)
)

type Client struct {
	lock    *sync.RWMutex
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewClient(configs *Configs) (*Client, error) {
	conn, err := amqp.Dial(configs.DSN)
	if err != nil {
		return nil, ClientCreationFailed.AddOriginalError(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, ClientCreationFailed.AddOriginalError(err)
	}

	c := &Client{
		lock:    new(sync.RWMutex),
		conn:    conn,
		channel: channel,
	}
	return c, nil
}

func (c *Client) Close() error {
	if !c.conn.IsClosed() {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	} else {
		return ConnectionClosedBefore
	}

	if !c.channel.IsClosed() {
		err := c.channel.Close()
		if err != nil {
			return err
		}
	} else {
		return ChannelClosedBefore
	}
	return nil
}

func (c *Client) Channel() *amqp.Channel {
	return c.channel
}
