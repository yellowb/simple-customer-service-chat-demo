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

// FakeDb 假的DB，跟每一个客户的交谈消息都保存在redis中
type FakeDb struct {
	redisClient *redis_client.Client
}

func GetFakeDb() *FakeDb {
	once.Do(func() {
		db = &FakeDb{
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
	// redis中的value就是整个dto.CustomerMessages的json序列化字符串
	err = f.redisClient.Client.SetEx(context.Background(), fmt.Sprintf("%s%s", constants.CustomerMsgPrefix, user), string(marshal), 5*time.Minute).Err()

	return err
}
