package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

type Remote struct {
	OnRead    func(n int, message []byte, err error) error
	write     chan []byte
	cleanup   chan struct{}
}

func (client *Remote) Write(p []byte) (int, error) {
	if client.write == nil {
		return 0, errors.New("cannot call Write() before client is initialized with client.Start()")
	}
	client.write <- p
	return len(p), nil

}

func (client *Remote) Close() error {
	client.cleanup <- struct{}{}
	return nil
}

func (client *Remote) Start(addr string, id string) error {

	log.Infof("Starting client to connect to %v", addr)

	done := make(chan struct{})

	client.cleanup = make(chan struct{})
	client.write = make(chan []byte)

	dialer := websocket.DefaultDialer

	u := url.URL{Scheme: "ws", Host: addr, Path: fmt.Sprintf("/%v", id)}

	log.Infof("connecting to %s", u.String())

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return errors.New(fmt.Sprintf("dial: %v", err))
	}

	log.Info("connected")

	defer conn.Close()

	go func() {
		defer close(done)
		for {
			select {
			default:
				n, msg, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, 1000) {
						log.WithField("action", "readmessage").Error(err)
					}
					return
				}
				if client.OnRead != nil {
					err := client.OnRead(n, msg, err)
					if err != nil {
						client.Write([]byte(err.Error()))
						log.WithField("func", "onread").Error(err)
						return
					}
				}
			}
		}
	}()

	for {
		select {
		case message := <-client.write:
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return err
			}

		case <-done:
			return nil

		case <-client.cleanup:
			log.Info("cleaning up")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			select {
			case <-done:
			case <-time.After(time.Second):
			}

			return err
		}
	}

	log.Info("Done.")

	return nil
}

