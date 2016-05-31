package consumer

// Message is a message that is pulled from a data source and is bound for consumption
type Message struct {
	Bytes []byte
}
