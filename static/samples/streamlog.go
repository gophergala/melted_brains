package streamlog

import (
	"sync"

	"github.com/youtube/vitess/go/stats"
)

var (
	sendCount         = stats.NewCounters("StreamlogSend")
	deliveredCount    = stats.NewMultiCounters("StreamlogDelivered", []string{"Log", "Subscriber"})
	deliveryDropCount = stats.NewMultiCounters("StreamlogDeliveryDroppedMessages", []string{"Log", "Subscriber"})
)

type subscriber struct {
	name string
}
type StreamLogger struct {
	name       string
	dataQueue  chan interface{}
	mu         sync.Mutex
	subscribed map[chan interface{}]subscriber
}

func New(name string, size int) *StreamLogger {
	logger := &StreamLogger{
		name:       name,
		dataQueue:  make(chan interface{}, size),
		subscribed: make(map[chan interface{}]subscriber),
	}
	go logger.stream()
	return logger
}

func (logger *StreamLogger) Send(message interface{}) {
	logger.dataQueue <- message
	sendCount.Add(logger.name, 1)
}

func (logger *StreamLogger) Subscribe(name string) chan interface{} {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	ch := make(chan interface{}, 1)
	logger.subscribed[ch] = subscriber{name: name}
	return ch
}
