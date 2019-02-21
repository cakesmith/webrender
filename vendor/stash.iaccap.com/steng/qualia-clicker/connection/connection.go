package connection

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/jackwakefield/gopac"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

var log = logrus.WithField("cmd", "connection")

type Client struct {
	OnRead    func(n int, message []byte, err error) error
	interrupt chan os.Signal
	write     chan []byte
	cleanup   chan struct{}
}

func (client *Client) Write(p []byte) (int, error) {
	if client.write == nil {
		return 0, errors.New("cannot call Write() before client is initialized with client.Start()")
	}
	client.write <- p
	return len(p), nil

}

func (client *Client) Close() error {
	if client.interrupt == nil {
		return errors.New("cannot call Close() before client is initialized with client.Start()")
	}
	client.cleanup <- struct{}{}
	return nil
}

func setupClientProxy(addr string) (*websocket.Dialer, error) {
	parser := gopac.Parser{}

	err := parser.ParseUrl("http://127.0.0.1:19876/pac.js")
	if err != nil {
		return nil, errors.New("error parsing proxy url " + err.Error())
	}

	proxy, err := parser.FindProxy(fmt.Sprintf("http://%v", addr), "localhost")
	if err != nil {
		return nil, errors.New("error finding proxy " + err.Error())
	}

	fmt.Printf("found PAC entry %v\n", proxy)

	uri, err := url.Parse("//" + strings.Fields(proxy)[1])
	if err != nil {
		return nil, errors.New("error parsing url " + err.Error())
	}

	uri.Scheme = "http"

	fmt.Printf("setting Dialer proxy to %v\n", uri.String())

	dialer := websocket.Dialer{}
	dialer.Proxy = http.ProxyURL(uri)
	return &dialer, nil
}

func (client *Client) Start(skipProxy bool, addr string, id string) error {

	log.Infof("Starting client to connect to %v ...\n", addr)

	client.interrupt = make(chan os.Signal, 1)
	signal.Notify(client.interrupt, os.Interrupt)

	done := make(chan struct{})

	client.cleanup = make(chan struct{})
	client.write = make(chan []byte)

	var dialer *websocket.Dialer

	if skipProxy {
		log.Info("using default Dialer")
		dialer = websocket.DefaultDialer
	} else {
		log.Info("using proxied Dialer to %v\n", addr)
		var err error
		dialer, err = setupClientProxy(addr)
		if err != nil {
			return errors.New("error setting up client proxy: " + err.Error())
		}
	}

	u := url.URL{Scheme: "ws", Host: addr, Path: fmt.Sprintf("/%v", id)}

	log.Infof("connecting to %s", u.String())

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return errors.New(fmt.Sprintf("dial: %v", err))
	}

	log.Info("client connected.")

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

		case <-client.interrupt:
			log.Info("interrupt received")
			client.cleanup <- struct{}{}

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
