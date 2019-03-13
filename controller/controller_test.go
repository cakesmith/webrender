package controller

import (
	"fmt"
	"github.com/goburrow/modbus"
	"log"
	"os"
	"testing"
	"time"
)

func XTestController(t *testing.T) {

	// Modbus TCP
	handler := modbus.NewTCPClientHandler("192.168.1.3:502")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0x01
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	defer handler.Close()
	if err != nil {
		t.Error(err)
	}

	client := modbus.NewClient(handler)
	results, err := client.ReadHoldingRegisters(0x0000, 2)
	if err != nil {
		t.Error(err)
	}

	whole := results[0:2]
	fract := results[2:4]

	scaling := GetScaling(whole, fract)
	fmt.Println(scaling)

	results, err = client.ReadHoldingRegisters(0x0018, 1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(results)

	fmt.Println("Voltage")
	fmt.Println(Scale(scaling, results))

	results, err = client.ReadHoldingRegisters(0x100F, 2)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(results)

}

func TestGetScaling(t *testing.T) {
	whole := []byte{0x00, 0x7b}
	fract := []byte{0xE0, 0x41}

	expected := float32(123.87599)

	actual := GetScaling(whole, fract)

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

}

func XTestScale(t *testing.T) {

}
