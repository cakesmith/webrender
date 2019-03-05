// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"io/ioutil"
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
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Id string

	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) Write(p []byte) (int, error) {
	c.send <- p
	return len(p), nil
}

func (c *Client) Close() error {
	c.hub.unregister <- c
	return c.conn.Close()
}

// readPump pumps messages from the websocket connection to the hub.
//
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
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Error(errors.Wrap(err, "error getting next reader"))
				}
				return
			}

			message, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Error(errors.Wrap(err, "connection read error"))
				return
			}
			//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			log.WithFields(logrus.Fields{"client": c.Id, "msg": string(message)}).Info("received")
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
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
				// The hub closed the channel.
				log.Info("hub closed channel")
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

			if err := w.Close(); err != nil {
				log.Error(errors.Wrap(err, "error closing websocket writer"))
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error(errors.Wrap(err,"error sending ping"))
				return
			}
		}
	}
}
