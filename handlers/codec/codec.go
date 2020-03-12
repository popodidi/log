package codec

import "github.com/popodidi/log/handlers"

// Default returns the default codec.
func Default(withColor bool) handlers.Codec {
	c := &defaultCodec{
		WithColor:  withColor,
		TimeFormat: timeFormat,
	}
	return c
}

// Simple returns the simple codec.
func Simple() handlers.Codec {
	return &simpleCodec{}
}
