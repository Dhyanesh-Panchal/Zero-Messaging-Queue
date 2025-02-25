package utils

import (
	"chatapp/centralserver/globals"
)

func CreateMessageBuffer() chan string {
	return make(chan string, globals.SharedMessageBufferSize)
}
