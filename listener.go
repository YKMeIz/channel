package channel

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

func (bc *Broadcast) listenAndServe() {
	for {
		select {
		case msg := <-bc.pub:
			go bc.broadcastMsg(msg)
		case msg := <-bc.uni:
			go bc.unicastMsg(msg)
		default:
		}
	}
}

func (bc *Broadcast) broadcastMsg(b []byte) {
	bc.sub.Range(func(key, value interface{}) bool {
		t := reflect.TypeOf(value)
		if t.Kind() == reflect.Func {
			value.(func([]byte))(b)
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
		uni: make(chan []byte),
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

func (bc *Broadcast) SubscribeFunc(id interface{}, f func([]byte)) {
	bc.sub.Store(id, f)
}

//func (bc *Broadcast) SendToID(recvID, sendID interface{}, payload []byte) {
//	v, ok := bc.sub.Load(recvID)
//	if !ok {
//		return
//	}
//
//	m := message{
//		Payload:  payload,
//		ReceiverID: recvID,
//		SenderID: sendID,
//	}
//
//	if reflect.TypeOf(v).Kind() == reflect.Func {
//		v.(func(b []byte))(m.gobEncode())
//	}
//
//}
//
func (m *message) gobEncode() []byte {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(m); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (m *message) gobDecode(b []byte) {
	var buf bytes.Buffer
	buf.Write(b)
	if err := gob.NewDecoder(&buf).Decode(&m); err != nil {
		panic(err)
	}
}

func (bc *Broadcast) SendToID(recvID, sendID interface{}, payload []byte) {
	m := message{
		Payload:    payload,
		ReceiverID: recvID,
		SenderID:   sendID,
	}

	b := m.gobEncode()
	bc.unicastMsg(b)
}

func (bc *Broadcast) unicastMsg(msg []byte) {
	var m message
	m.gobDecode(msg)
	v, ok := bc.uniPool.Load(m.ReceiverID)
	if !ok {
		return
	}
	if reflect.TypeOf(v).Kind() == reflect.Func {
		if resp := v.(func(payload []byte) []byte)(m.Payload); resp != nil {
			bc.SendToID(m.SenderID, m.ReceiverID, resp)
		}
	}
}

func (bc *Broadcast) SendToIDResponse(recvID, sendID interface{}, payload []byte) []byte {
	m := message{
		Payload:    payload,
		ReceiverID: recvID,
		SenderID:   sendID,
	}

	b := m.gobEncode()
	return bc.unicastMsgResponse(b)
}

func (bc *Broadcast) unicastMsgResponse(msg []byte) []byte {
	var m message
	m.gobDecode(msg)
	v, ok := bc.uniPool.Load(m.ReceiverID)
	if !ok {
		return nil
	}
	if reflect.TypeOf(v).Kind() == reflect.Func {
		return v.(func(payload []byte) []byte)(m.Payload)
	}
	return nil
}

func (bc *Broadcast) RegisterHandleFunc(id interface{}, handler func([]byte) []byte) {
	bc.uniPool.Store(id, handler)
}
