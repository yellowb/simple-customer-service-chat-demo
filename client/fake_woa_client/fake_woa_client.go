package fake_woa_client

import (
	"log"
	"yellowb.com/chat-demo/db/fake_db"
)

// FakeWoaClient 虚假WOA Client
type FakeWoaClient struct{}

func NewFakeWoaClient() *FakeWoaClient {
	return &FakeWoaClient{}
}

func (c *FakeWoaClient) PushMsgToWoa(msg *fake_db.Message) {
	log.Printf("[FakeWoaClient] push msg [%s] to woa user [%s]", msg.Text, msg.User)
	// Do nothing
}
