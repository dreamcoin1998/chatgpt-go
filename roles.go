// chatgpt可以预定义交互行为，本文件主要用来预定义Chatgpt的行为

package chat

var (
	MongkeyKing = ChatMode{
		ModeName: "孙悟空",
		PreMessages: []Message{{ // 美猴王孙悟空的角色
			Role:    SystemRole,
			Content: "请你扮演孙悟空的角色和我对话",
		}},
	}
)

var PreRoles = []ChatMode{MongkeyKing}
