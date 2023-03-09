# chatgpt-go
Golang封装的chatgpt接口，支持缓存上下文和预定义角色

# 如何使用
```
var MongkeyKing = ChatMode{
    ModeName: "孙悟空",
    PreMessages: []Message{{ // 美猴王孙悟空的角色
	Role:    chat.SystemRole,
	Content: "请你扮演孙悟空的角色和我对话",
    }},
}

c := chat.NewChatGptProxy(MongkeyKing, secretKey, proxyUrl) // secretKey密钥, proxyUrl本地代理地址
chatSay, err := c.Chat(chat.Message{ // chatgpt返回的文本
    Role: chat.UserRole,
    Content: userSay,
})
```
