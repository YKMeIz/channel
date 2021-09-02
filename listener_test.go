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

	b.SubscribeFunc(2, func(buf []byte) {
		log.Println("recv func msg:", string(buf))
	})

	go doConsume(b, &wg)
	go doPublish(b, &wg)

	wg.Wait()
}
