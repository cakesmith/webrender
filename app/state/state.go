package state

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

type State map[string]interface{}

// TextMarshaler is the interface implemented by an object that
// can marshal itself into a textual form.

//MarshalText encodes the receiver into UTF-8-encoded text and returns the result.
func (state *State) MarshalText() (text []byte, err error) {

	b := new(bytes.Buffer)
	b64 := new(bytes.Buffer)

	e := gob.NewEncoder(b)

	// Encoding the map
	err = e.Encode(*state)
	if err != nil {
		return
	}

	e64 := base64.NewEncoder(base64.StdEncoding, b64)
	_, err = e64.Write(b.Bytes())
	if err != nil {
		return
	}

	// Must close the encoder when finished to flush any partial blocks.
	// If you comment out the following line, the last partial block "r"
	// won't be encoded.
	err = e64.Close()
	if err != nil {
		return
	}

	return b64.Bytes(), nil
}

//TextUnmarshaler is the interface implemented by an object
// that can unmarshal a textual representation of itself.
//UnmarshalText must be able to decode the form generated
// by MarshalText. UnmarshalText must copy the text if it wishes
// to retain the text after returning.

func (state *State) UnmarshalText(text []byte) error {

	b64 := new(bytes.Buffer)

	d64 := base64.NewDecoder(base64.StdEncoding, b64)

	_, err := d64.Read(text)
	if err != nil {
		return err
	}

	d := gob.NewDecoder(b64)

	// Decoding the map
	err = d.Decode(state)
	if err != nil {
		return err
	}

	return nil
}

type Stateful interface {
	GetState() State
	SetState(State)
}
