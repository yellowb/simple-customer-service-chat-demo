package fake_db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
	"yellowb.com/chat-demo/client/redis_client"
	"yellowb.com/chat-demo/constants"
	"yellowb.com/chat-demo/dto"
)

var (
	once sync.Once
	db   *FakeDb
)

// FakeDb 假的DB
type FakeDb struct {
	Data        map[string]*dto.CustomerMessages // 跟所有客户的交谈记录。key = 客户user id
	redisClient *redis_client.Client
}

func GetFakeDb() *FakeDb {
	once.Do(func() {
		db = &FakeDb{
			Data:        make(map[string]*dto.CustomerMessages),
			redisClient: redis_client.GetClient(),
		}
	})
	return db
}

// GetCustomerMessages 获取跟某个客户的所有交谈消息集合
func (f *FakeDb) GetCustomerMessages(user string) (*dto.CustomerMessages, error) {
	val, err := f.redisClient.Client.Get(context.Background(), fmt.Sprintf("%s%s", constants.CustomerMsgPrefix, user)).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Printf("Failed to get customer messages from redis %v", err)
			return nil, err
		}
		return nil, nil
	}
	result := &dto.CustomerMessages{}
	err = json.Unmarshal([]byte(val), result)
	if err != nil {
		log.Printf("Failed to unmarshal customer messages %v", err)
		return nil, err
	}

	return result, nil
}

// SaveCustomerMessage 保存跟某个客户的一条交谈消息
func (f *FakeDb) SaveCustomerMessage(user string, message *dto.Message) error {

	customerMessages, err := f.GetCustomerMessages(user)
	if err != nil {
		return err
	}
	if customerMessages == nil {
		customerMessages = &dto.CustomerMessages{
			Messages: make([]*dto.Message, 0),
		}
	}
	customerMessages.Messages = append(customerMessages.Messages, message)

	marshal, _ := json.Marshal(customerMessages)
	err = f.redisClient.Client.SetEx(context.Background(), fmt.Sprintf("%s%s", constants.CustomerMsgPrefix, user), string(marshal), 5*time.Minute).Err()

	return err

	//var customerMessages *dto.CustomerMessages
	//// 如果某个客户从来没交谈过，则新建
	//if customerMessages = f.Data[user]; customerMessages == nil {
	//	customerMessages = &dto.CustomerMessages{
	//		Messages: make([]*dto.Message, 0),
	//	}
	//	f.Data[user] = customerMessages
	//}
	//customerMessages.Messages = append(customerMessages.Messages, message)
}
