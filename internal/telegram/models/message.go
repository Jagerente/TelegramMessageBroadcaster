package models

type Message struct {
	Id    uint64
	Text  map[uint64]string
	Photo string
}

type MessageType int

const (
	TextMessage MessageType = iota
	PhotoMessage
)

func (m *Message) GetType() MessageType {
	if m.Photo != "" {
		return PhotoMessage
	}

	return TextMessage
}

func CreateMessage() *Message {
	return &Message{
		Text:  make(map[uint64]string),
		Photo: "",
	}
}
