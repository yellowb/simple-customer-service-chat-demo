package constants

const (
	// CustomerServiceMsgNotifyChan 用于在有新的用户发问时广播消息
	CustomerServiceMsgNotifyChan = "customer_service_msg_notify_chan"
	// CustomerMsgPrefix redis中存储单个客户的对话消息的key。customer_msg:<userId>
	CustomerMsgPrefix = "customer_msg:"
)
