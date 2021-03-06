// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	log = logrus.WithField("cmd", "hub")
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.

type Hub struct {
	sync.Mutex

	//queue of waiting tests
	waiting map[string]chan *Client

	// Registered clients.
	clients map[*Client]bool

	// Send message to clients.
	Send chan *Message

	OnRegister func(*Client)

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type Message struct {
	Data []byte
	To   []*Client
}

func (hub *Hub) Close() error {
	for client := range hub.clients {
		client.Close()
	}
	return nil
}

func NewHub() (*Hub, error) {

	hub := &Hub{
		Send:       make(chan *Message),
		waiting:    make(map[string]chan *Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	go hub.run()

	return hub, nil
}

type Events struct {
	OnClick    func(btn, x, y int)
	OnKeypress func(key int)
}

func (hub *Hub) Handler(ev *Events) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		unescaped, err := url.PathUnescape(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		split := strings.Split(unescaped, "/")

		if split[0] == "" {
			split = split[1:]
		}

		if split[0] == "" {
			http.NotFound(w, r)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUpgradeRequired)
		}

		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), Events: ev}

		client.Id = split[0]

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()

		hub.register <- client

	}
}

func (hub *Hub) run() {
	for {
		select {

		//Register a client
		case client := <-hub.register:

			log.WithField("client", client.Id).Info("registering")

			hub.Lock()
			hub.clients[client] = true

			for id, ch := range hub.waiting {
				if client.Id == id {
					go func(ch chan *Client, client *Client) {
						ch <- client
					}(ch, client)
					delete(hub.waiting, id)
				}
			}

			hub.Unlock()

			log.WithFields(logrus.Fields{"client": client.Id, "addr": client.conn.RemoteAddr()}).Info("registered")

			if hub.OnRegister != nil {
				go hub.OnRegister(client)
			}

		//Unregister a client
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {

				log.WithFields(logrus.Fields{"client": client.Id, "addr": client.conn.RemoteAddr()}).Info("unregistered")

				hub.Lock()
				delete(hub.clients, client)
				hub.Unlock()

			}

		case message := <-hub.Send:

			for _, to := range message.To {
				isin := false
				for registered := range hub.clients {
					isin = isin || to == registered
				}
				if !isin {
					log.WithField("addr", to.conn.RemoteAddr()).Error("does not exist in registry")
					continue
				}
				select {
				case to.send <- message.Data:
					log.WithFields(logrus.Fields{"addr": to.conn.RemoteAddr(), "msg": string(message.Data)}).Info("sending")
				default:
					to.Close()
					delete(hub.clients, to)
				}
			}

		}
	}
}

func (hub *Hub) Subscribe(id string) chan *Client {
	ch := make(chan *Client)

	hub.Lock()
	defer hub.Unlock()

	for cl := range hub.clients {
		if cl.Id == id {
			go func(ch chan *Client, cl *Client) {
				ch <- cl
			}(ch, cl)
			return ch
		}
	}

	hub.waiting[id] = ch

	return ch
}
