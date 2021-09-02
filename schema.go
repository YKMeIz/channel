package channel

import "sync"

type Broadcast struct {
	sub sync.Map
	pub chan []byte
}
