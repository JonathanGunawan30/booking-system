package kafka

type Registry struct {
	brokers []string
}

type KafkaRegistryInterface interface {
	GetKafkaProducer() KafkaInterface
}

func NewKafkaRegistry(brokers []string) KafkaRegistryInterface {
	return &Registry{brokers: brokers}
}

func (r *Registry) GetKafkaProducer() KafkaInterface {
	return NewKafkaProducer(r.brokers)
}
