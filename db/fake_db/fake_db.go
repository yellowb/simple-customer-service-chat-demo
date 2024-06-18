package fake_db

import (
	"sync"
	"yellowb.com/chat-demo/dto"
)

var (
	once sync.Once
	db   *FakeDb
)

// FakeDb 假的DB
type FakeDb struct {
	Data map[string]*dto.CustomerMessages // 跟所有客户的交谈记录。key = 客户user id
}

func GetFakeDb() *FakeDb {
	once.Do(func() {
		db = &FakeDb{
			Data: make(map[string]*dto.CustomerMessages),
		}
	})
	return db
}

// GetCustomerMessages 获取跟某个客户的所有交谈消息集合
func (f *FakeDb) GetCustomerMessages(user string) *dto.CustomerMessages {
	return f.Data[user]
}

// SaveCustomerMessage 保存跟某个客户的一条交谈消息
func (f *FakeDb) SaveCustomerMessage(user string, message *dto.Message) {
	var customerMessages *dto.CustomerMessages
	// 如果某个客户从来没交谈过，则新建
	if customerMessages = f.Data[user]; customerMessages == nil {
		customerMessages = &dto.CustomerMessages{
			Messages: make([]*dto.Message, 0),
		}
		f.Data[user] = customerMessages
	}
	customerMessages.Messages = append(customerMessages.Messages, message)
}
