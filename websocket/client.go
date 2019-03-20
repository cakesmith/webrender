// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// Request timeout
	reqTimeout = time.Second * 2
)

var (
	log     = logrus.New()
	newline = []byte{'\n'}
)

type Receiver interface {
	OnRecv([]byte)
}

type Registrar interface {
	OnRegister()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Command struct {
	Name   string
	Params []string
}

// This format is used for a fire and forget command
func (cmd Command) MakePacket() []byte {
	return []byte(strings.Join(cmd.MakePayload(), " "))
}

// This format is used for a Request command
func (cmd Command) MakePayload() []string {
	return append([]string{cmd.Name}, cmd.Params...)
}

type Client struct {
	mutex sync.Mutex

	pending map[string]*call

	counter uint64

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type Message struct {
	*Client
	Data []byte
}

func NewClient() *Client {
	return &Client{
		pending: make(map[string]*call, 1),
	}
}

// call represents an active request
type call struct {
	res   []byte
	done  chan bool
	Error error
}

type Actions interface {
	OnRecv([]byte)
	OnRegister()
}

func newCall() *call {
	done := make(chan bool)
	return &call{
		done: done,
	}
}

func (client *Client) Request(payload []string) ([]byte, error) {

	guid := xid.New().String()

	client.mutex.Lock()
	call := newCall()
	client.pending[guid] = call

	params := append([]string{fmt.Sprint(guid)}, payload...)

	_, err := client.Write(Command{
		Name:   "req",
		Params: params,
	}.MakePacket())

	if err != nil {
		delete(client.pending, guid)
		client.mutex.Unlock()
		return nil, err
	}

	client.mutex.Unlock()

	select {
	case <-call.done:
	case <-time.After(reqTimeout):
		call.Error = errors.New("request timeout")
	}

	if call.Error != nil {
		client.mutex.Lock()
		delete(client.pending, guid)
		client.mutex.Unlock()
		return nil, call.Error
	}

	return call.res, nil
}

func (client *Client) MakeHandler(actions Actions) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUpgradeRequired)
		}

		client.send = make(chan []byte, 256)
		client.conn = conn

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump(actions.OnRecv)

		actions.OnRegister()

	}
}

func (client *Client) Write(p []byte) (int, error) {
	client.send <- p
	return len(p), nil
}

func (client *Client) Close() error {
	return client.conn.Close()
}

// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (client *Client) readPump(onrecv func([]byte)) {

	defer client.Close()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { return client.conn.SetReadDeadline(time.Now().Add(pongWait)) })

	log.WithFields(logrus.Fields{
		"readLimit": maxMessageSize,
	}).Info("starting read pump")

	for {
		select {
		default:

			_, reader, err := client.conn.NextReader()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Error(errors.Wrap(err, "error getting next reader"))
				}
				return
			}

			message, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Error(errors.Wrap(err, "connection read error"))
				return
			}

			logrus.WithField("msg", string(message)).Info("received")

			split := strings.Split(string(message), " ")

			// see if this is a response
			if split[0] == "res" {

				id := string(split[1])

				res := split[2:]

				client.mutex.Lock()

				call := client.pending[id]
				delete(client.pending, id)

				client.mutex.Unlock()

				if call == nil {
					err = errors.New("no pending request found")
					continue
				}

				call.res = []byte(strings.Join(res, " "))
				call.done <- true

			} else {
				go onrecv(message)
			}

		}
	}
}

// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (client *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)

	log.WithField("pingPeriod", pingPeriod).Info("starting write pump")

	defer func() {
		ticker.Stop()
		client.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:

			client.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Error(errors.Wrap(err, "error getting next writer"))
				}
				return
			}

			log.WithField("msg", string(message)).Trace("writing")

			if _, err := w.Write(message); err != nil {
				log.Error(errors.Wrap(err, "error writing"))
				return
			}

			// Add queued messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				log.Error(errors.Wrap(err, "error closing websocket writer"))
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					log.Error(errors.Wrap(err, "error sending ping"))
				}
				return
			}
		}
	}
}
