package channel

import (
	"bytes"
	"log"
	"strconv"
	"sync"
	"testing"
)

const (
	broadcastCount = 5
	singleMsgCount = 7
)

func doPublish(b *Broadcast, wg *sync.WaitGroup) {
	for i := 1; i <= broadcastCount; i++ {
		log.Println("doPush: broadcast message ...")
		b.Publish([]byte(strconv.Itoa(i)))
	}

	for i := 1; i <= singleMsgCount; i++ {
		log.Println("doPush: send message to ID:3 ...")
		b.SendToID(3, 0, []byte("send to id 2..."))
	}

	wg.Done()
}

func doConsume(b *Broadcast, wg *sync.WaitGroup, t *testing.T) {
	c := make(chan []byte)
	b.Subscribe(1, c)

	count := 0

	for {
		select {
		case msg := <-c:
			{
				if bytes.Compare(msg, []byte(strconv.Itoa(4))) == 0 {
					log.Println("doConsume: stop as received message 4 ...")
					wg.Done()
					break
				}
				log.Println("doConsume: recv msg:", string(msg))
				count++
			}
		default:

		}
	}

	if count != broadcastCount {
		t.Error("channel subscribe has received missed or duplicated message:", count, broadcastCount)
	}
}

func TestBroadcast(t *testing.T) {
	b := NewBroadcast()

	var wg sync.WaitGroup
	wg.Add(2)

	var (
		reg0, reg3, reg2 int
	)

	b.RegisterHandleFunc(0, func(buf []byte) []byte {
		log.Println("ID:0: recv func msg:", string(buf))
		reg0++
		return nil
	})

	b.RegisterHandleFunc(3, func(buf []byte) []byte {
		log.Println("ID:3: recv func msg:", string(buf))
		reg3++
		return []byte("ok")
	})

	b.SubscribeFunc(2, func(buf []byte) {
		log.Println("ID:2: recv func msg:", string(buf))
		reg2++
	})

	go doConsume(b, &wg, t)
	go doPublish(b, &wg)

	wg.Wait()

	if reg0 != singleMsgCount {
		t.Error("handler func 0 has received missed or duplicated message:", reg0, singleMsgCount)
	}

	if reg3 != singleMsgCount {
		t.Error("handler func 3 has received missed or duplicated message:", reg3, singleMsgCount)
	}

	if reg2 != broadcastCount {
		t.Error("subscribe func 2 has received missed or duplicated message:", reg2, broadcastCount)
	}
}

func TestUnicasetWithResponse(t *testing.T) {
	b := NewBroadcast()
	respMsg := []byte("b1: hello from main")

	b.RegisterHandleFunc("b1", func(buf []byte) []byte {
		log.Println("ID:b1: recv func msg:", string(buf))
		return respMsg
	})

	resp := b.SendToIDResponse("b1", nil, []byte("hello from main"))
	if bytes.Compare(resp, respMsg) != 0 {
		t.Error("message sent and received mismatch")
	}
}
