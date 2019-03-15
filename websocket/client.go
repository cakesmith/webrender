// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
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
	//space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	mutex sync.Mutex

	pending map[uint64]*call

	counter uint64

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// call represents an active request
type call struct {
	res   []byte
	done  chan bool
	Error error
}

func newCall() *call {
	done := make(chan bool)
	return &call{
		done: done,
	}
}

func New() *Client {
	return &Client{
		pending: make(map[uint64]*call, 1),
	}
}

type Command struct {
	Name   string
	Params []string
}

func (c Command) MakePacket() []byte {
	return []byte(strings.Join(append([]string{c.Name}, c.Params...), " "))
}

func (c *Client) Request(payload []string) ([]byte, error) {

	c.mutex.Lock()
	id := c.counter
	c.counter++
	call := newCall()
	c.pending[id] = call

	params := append([]string{fmt.Sprint(id)}, payload...)

	_, err := c.Write(Command{
		Name:   "req",
		Params: params,
	}.MakePacket())

	if err != nil {
		delete(c.pending, id)
		c.mutex.Unlock()
		return nil, err
	}

	c.mutex.Unlock()

	select {
	case <-call.done:
	case <-time.After(reqTimeout):
		call.Error = errors.New("request timeout")
	}

	if call.Error != nil {
		return nil, call.Error
	}

	return call.res, nil
}

func (c *Client) Handler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUpgradeRequired)
	}

	c.send = make(chan []byte, 256)
	c.conn = conn

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go c.writePump()
	go c.readPump()

}

func (c *Client) Write(p []byte) (int, error) {
	c.send <- p
	return len(p), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {

	defer c.Close()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(pongWait)) })

	log.WithFields(logrus.Fields{
		"readLimit": maxMessageSize,
	}).Info("starting read pump")

	for {
		select {
		default:

			_, reader, err := c.conn.NextReader()
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
				id, err := strconv.ParseUint(split[1], 10, 64)
				if err != nil {
					log.Error(err)
				}

				res := split[2:]

				c.mutex.Lock()

				call := c.pending[id]
				delete(c.pending, id)

				c.mutex.Unlock()

				if call == nil {
					err = errors.New("no pending request found")
					continue
				}

				call.res = []byte(strings.Join(res, " "))
				call.done <- true

			}

		}
	}
}

// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)

	log.WithField("pingPeriod", pingPeriod).Info("starting write pump")

	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)

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
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				log.Error(errors.Wrap(err, "error closing websocket writer"))
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					log.Error(errors.Wrap(err, "error sending ping"))
				}
				return
			}
		}
	}
}
