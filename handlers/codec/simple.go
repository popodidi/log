package codec

import (
	"fmt"

	"github.com/popodidi/log"
)

type simpleCodec struct{}

func (c *simpleCodec) Encode(entry *log.Entry) []byte {
	content := entry.Log
	if entry.Level <= log.Error {
		content = fmt.Sprintf("%s: %s", entry.DebugInfo(), entry.Log)
	}
	return []byte(content)
}
