# simple-customer-service-chat-demo
一个简单的基于redis、浏览器SSE的客服服务台demo

# 前提需求
1. 需要一个redis跑在本地，6379端口，用于多个Server之间广播消息

# 怎么运行？
1. 启动2个（或多个）cmd/main.go的实例，启动时加 -p 参数指定Server端口。这里假设启动2个，一个是8085端口（Server A），一个是8086端口（Server B）。
2. 打开4个浏览器页面：
    1. http://localhost:8085/ask （Server A的客户提问页面）
    2. http://localhost:8085/answer （Server A的客服回答页面）
    3. http://localhost:8086/ask （Server B的客户提问页面）
    4. http://localhost:8086/answer （Server B的客服回答页面）
3. 在（Server A的客户提问页面）提交一个问题，可以见到（Server A的客服回答页面）、（Server B的客服回答页面）都会收到通知并刷新数据
4. 在（Server A的客服回答页面）提交一个回答，可以见到（Server A的客服回答页面）、（Server B的客服回答页面）都会收到通知并刷新数据

# 注意
1. 只有/answer的页面才有接入SSE自动刷新数据，/ask的提问页面没有SSE需要手工点击刷新 。