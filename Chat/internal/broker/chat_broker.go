package broker

import (
	"sync"

	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)

type ChatBroker struct {
	mu          sync.RWMutex
	userStreams map[int]chan string
	logger      *platform.Logger
}

func NewChatBroker(logger *platform.Logger) *ChatBroker {

	broker := &ChatBroker{
		userStreams: make(map[int]chan string),
		logger:      logger,
	}

	return broker
}

func (b *ChatBroker) Register(userID int) chan string {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan string, 10)
	b.userStreams[userID] = ch
	return ch
}

func (b *ChatBroker) Unregister(userID int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if ch, exists := b.userStreams[userID]; exists {
		close(ch)
		delete(b.userStreams, userID)
	}
}

func (b *ChatBroker) Publish(userID int, message string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ch, exists := b.userStreams[userID]
	if exists {
		select {

		case ch <- message:

		default:

			b.logger.Warn("Chat Dropped", zap.Int16("userID", int16(userID)))

		}
	}
	if !exists {

		b.logger.Error("[ChatBroker Publish] UserID doesn't exist")
	}

}
