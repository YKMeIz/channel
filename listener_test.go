package channel

import (
	"bytes"
	"log"
	"strconv"
	"sync"
	"testing"
)

func doPublish(b *Broadcast, wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		b.Publish([]byte(strconv.Itoa(i)))
		b.SendToID(3, 0, []byte("send to id 2..."))
	}
	wg.Done()
}

func doConsume(b *Broadcast, wg *sync.WaitGroup) {
	c := make(chan []byte)
	b.Subscribe(1, c)
	for {
		select {
		case msg := <-c:
			{
				if bytes.Compare(msg, []byte(strconv.Itoa(4))) == 0 {
					wg.Done()
					break
				}
				log.Println("recv msg:", string(msg))
			}
		default:

		}
	}
}

func TestBroadcast(t *testing.T) {
	b := NewBroadcast()

	var wg sync.WaitGroup
	wg.Add(2)

	b.RegisterHandleFunc(0, func(buf []byte) []byte {
		log.Println("0: recv func msg:", string(buf))
		return nil
	})

	b.RegisterHandleFunc(3, func(buf []byte) []byte {
		log.Println("3: recv func msg:", string(buf))
		return []byte("ok")
	})

	b.SubscribeFunc(2, func(buf []byte) {
		log.Println("2: recv func msg:", string(buf))
	})

	go doConsume(b, &wg)
	go doPublish(b, &wg)

	wg.Wait()
}
