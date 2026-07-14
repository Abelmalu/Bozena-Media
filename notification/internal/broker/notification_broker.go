package broker

import (
	"sync"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)

type NotificationBroker struct {
	mu          sync.RWMutex
	userStreams map[int]chan string
	logger      *platform.Logger
}

func NewNotificationBroker(logger *platform.Logger) *NotificationBroker {

	broker := &NotificationBroker{
		userStreams: make(map[int]chan string),
		logger: logger,
	}

	return broker
}



func (b *NotificationBroker) Register(userID int) chan string {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan string, 10)
	b.userStreams[userID] = ch
	return ch
}

func (b *NotificationBroker) Unregister(userID int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if ch, exists := b.userStreams[userID]; exists {
		close(ch)
		delete(b.userStreams, userID)
	}
}

func (b *NotificationBroker) NotifyUser(userID int, message string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if ch, exists := b.userStreams[userID]; exists {
		select {

		case ch <- message:

		default:

			b.logger.Warn("Notification Dropped",zap.Int16("userID",int16(userID)))


		}
	}
}
