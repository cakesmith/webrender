package websocket_test

import (
	"github.com/cakesmith/webrender/websocket"
	"log"
	"sync"
	"testing"
)


func TestWebsocket(t *testing.T) {

	message := "hello, websocket!"
	clientId := "client 1"

	hub, err := websocket.NewHub("localhost:0")
	if err != nil {
		t.Error(err)
	}
	if !hub.Started() {
		t.Error("hub.Started() should be true")
	}

	defer func() {
		err := hub.Close()
		if err != nil {
			t.Error("error closing hub: " + err.Error())
		}
		if hub.Started() {
			t.Error("hub.Started() should be false")
		}
	}()

	if !hub.Started() {
		t.Error("hub.Started() should be true")
	}
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
		err := remote.Start(hub.Addr(), clientId)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	remote.Close()

	hub.Close()

}
