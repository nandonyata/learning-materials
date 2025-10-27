package product_service

import (
	"learn-graphql-go-gorm/datalayer/models"
	"learn-graphql-go-gorm/graph/model"
	"sync"
	"time"

	"github.com/google/uuid"
)

// EventBroadcaster manages subscription channels
type EventBroadcaster struct {
	mu          sync.RWMutex
	subscribers map[string]chan *model.ProductChangeEvent
}

var broadcaster *EventBroadcaster
var once sync.Once

// GetBroadcaster returns singleton instance
func GetBroadcaster() *EventBroadcaster {
	once.Do(func() {
		// This code runs ONLY ONCE, even if called 1000 times simultaneously
		broadcaster = &EventBroadcaster{
			subscribers: make(map[string]chan *model.ProductChangeEvent),
		}
	})
	return broadcaster
}

func (b *EventBroadcaster) Subscribe() <-chan *model.ProductChangeEvent {
	b.mu.Lock()         // Lock the list
	defer b.mu.Unlock() // Unlock when done

	ch := make(chan *model.ProductChangeEvent, 1) // Create a mailbox
	id := uuid.New().String()                     // Give it a unique ID
	b.subscribers[id] = ch                        // Direct assignment

	return ch // Return the mailbox so client can read from it
}

func (b *EventBroadcaster) BroadcastProductChange(action model.ProductAction, product *models.Product) {
	b.mu.RLock() // Lock for reading
	defer b.mu.RUnlock()

	// Create the message
	event := &model.ProductChangeEvent{
		Action:    action, // CREATED, UPDATED
		Product:   product.ProductToGraphQLProduct(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Send to ALL subscribers
	for _, ch := range b.subscribers {
		select {
		case ch <- event: // Try to send
		default:
			// Skip if mailbox is full (non-blocking)
		}
	}
}
