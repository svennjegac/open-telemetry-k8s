package bootstrap

import "user-events-kafka/internal/userevents"

func Consumer() *userevents.Consumer {
	return userevents.NewConsumer()
}
