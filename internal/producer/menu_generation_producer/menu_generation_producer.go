package menu_generation_producer

type MenuGenerationProducer struct {
	kafkaBroker []string
	topicName   string
}

func NewMenuGenerationProducer(kafkaBroker []string, topicName string) *MenuGenerationProducer {
	return &MenuGenerationProducer{
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
	}
}
