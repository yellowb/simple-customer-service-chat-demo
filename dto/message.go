package dto

// CustomerMessages 跟某个用户的所有交谈消息集合，包含客户->客服的、客服->客户的
type CustomerMessages struct {
	Messages []*Message `json:"messages"`
}

// Message 一条交谈消息
type Message struct {
	User         string `json:"user"`          // 客户ID
	Helper       string `json:"helper"`        // 客服ID（客户发问时为空）
	Text         string `json:"text"`          // 消息文本
	FromCustomer bool   `json:"from_customer"` // true = 客户发的消息；false = 客服发的消息
	Ts           int64  `json:"ts"`            // 时间戳
}

// GetLatestMessage 获取最新一条消息
func (c *CustomerMessages) GetLatestMessage() *Message {
	return c.Messages[len(c.Messages)-1]
}
