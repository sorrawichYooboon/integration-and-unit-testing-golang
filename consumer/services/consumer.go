package services

import "github.com/Shopify/sarama"

type consumerService struct {
	eventService IEventService
}

func NewConsumerService(eventService IEventService) sarama.ConsumerGroupHandler {
	return consumerService{eventService}
}

func (obj consumerService) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (obj consumerService) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (obj consumerService) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		obj.eventService.Handle(msg.Topic, msg.Value)
		session.MarkMessage(msg, "")
	}

	return nil
}
