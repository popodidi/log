package handlers

import "github.com/popodidi/log"

// Codec encodes log entries into bytes.
type Codec interface {
	Encode(*log.Entry) []byte
}
