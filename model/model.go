package model

func CreateMessages() *Messages {
	messages := make([]*MessageData, 0)
	for i := 0; i < 7000; i++ {
		messages = append(messages, &MessageData{
			To:      "dL3ZtkNM",
			From:    "uBYiJ1HD",
			Message: "hQARuTmG",
		})
	}
	return &Messages{
		MessageList: messages,
	}
}
