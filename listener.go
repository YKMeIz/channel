package channel

import "reflect"

func (bc *Broadcast) listenAndServe() {
	for {
		select {
		case msg := <-bc.pub:
			go bc.broadcastMsg(msg)
		default:
		}
	}
}

func (bc *Broadcast) broadcastMsg(b []byte) {
	bc.sub.Range(func(key, value interface{}) bool {
		t := reflect.TypeOf(value)
		if t.Kind() == reflect.Func {
			value.(func(b []byte))(b)
		}
		if t.Kind() == reflect.Chan {
			value.(chan []byte) <- b
		}

		return true
	})
}

func NewBroadcast() *Broadcast {
	bc := &Broadcast{
		pub: make(chan []byte),
	}
	go bc.listenAndServe()
	return bc
}

func (bc *Broadcast) Subscribe(id interface{}, c chan []byte) {
	bc.sub.Store(id, c)
}

func (bc *Broadcast) Unsubscribe(id interface{}) {
	bc.sub.Delete(id)
}

func (bc *Broadcast) Publish(b []byte) {
	bc.pub <- b
}

func (bc *Broadcast) SubscribeFunc(id interface{}, f func(b []byte)) {
	bc.sub.Store(id, f)
}
