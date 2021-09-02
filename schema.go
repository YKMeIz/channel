package channel

import (
	"sync"
)

type Broadcast struct {
	sub     sync.Map
	uniPool sync.Map
	pub     chan []byte
	uni     chan []byte
}

type message struct {
	Payload              []byte
	ReceiverID, SenderID interface{}
}
