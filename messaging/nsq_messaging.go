package messaging

import "github.com/nsqio/go-nsq"

type MessageHandler struct{}

// HandleMessage implements the Handler interface.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		return nil
	}

	err := processMessage(m.Body)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return err
}
