package fake_db

// FakeDb 假的DB
type FakeDb struct {
	Data map[string]*CustomerMessages // 跟所有客户的交谈记录。key = 客户user id
}

// GetCustomerMessages 获取跟某个客户的所有交谈消息集合
func (f *FakeDb) GetCustomerMessages(user string) *CustomerMessages {
	return f.Data[user]
}

// SaveCustomerMessage 保存跟某个客户的一条交谈消息
func (f *FakeDb) SaveCustomerMessage(user string, message *Message) {
	var customerMessages *CustomerMessages
	// 如果某个客户从来没交谈过，则新建
	if customerMessages = f.Data[user]; customerMessages == nil {
		customerMessages = &CustomerMessages{
			Messages: make([]*Message, 0),
		}
		f.Data[user] = customerMessages
	}
	customerMessages.Messages = append(customerMessages.Messages, message)
}

// CustomerMessages 跟某个用户的所有交谈消息集合，包含客户->客服的、客服->客户的
type CustomerMessages struct {
	Messages []*Message
}

// Message 一条交谈消息
type Message struct {
	User         string `json:"user"`          // 谁发的消息
	Text         string `json:"text"`          // 消息文本
	FromCustomer bool   `json:"from_customer"` // true = 客户发的消息；false = 客服发的消息
	Ts           int64  `json:"ts"`            // 时间戳
}

// GetLatestMessage 获取最新一条消息
func (c *CustomerMessages) GetLatestMessage() *Message {
	return c.Messages[len(c.Messages)-1]
}
