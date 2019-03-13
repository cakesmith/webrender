package websocket_test

import (
	"github.com/cakesmith/webrender/websocket"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestWebsocket(t *testing.T) {

	message := "hello, websocket!"
	clientId := "client 1"

	hub, err := websocket.NewHub()
	if err != nil {
		t.Error(err)
	}

	server := httptest.NewServer(http.HandlerFunc(hub.Handler(&websocket.Events{})))

	defer func() {
		server.Close()
		err := hub.Close()
		if err != nil {
			t.Error("error closing hub: " + err.Error())
		}

	}()

	arrival := hub.Subscribe(clientId)

	var wg sync.WaitGroup
	wg.Add(2)

	remote := websocket.Remote{}

	remote.OnRead = func(n int, messageBytes []byte, err error) error {
		defer wg.Done()
		if string(messageBytes) != message {
			t.Errorf("epected %v received %v", message, string(messageBytes))
		}
		return nil
	}

	go func() {
		defer wg.Done()
		client := <-arrival
		if client.Id != clientId {
			t.Errorf("expected %v, received %v", clientId, client.Id)
		}
		_, err := client.Write([]byte(message))
		if err != nil {
			t.Error(err)
		}

	}()

	go func() {
		err := remote.Start(server.URL[7:], clientId)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	remote.Close()

}
