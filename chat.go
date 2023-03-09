package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// 用来预定义chatgpt角色
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (
	AssistantRole = "assistant" // 助手
	SystemRole    = "system"    // 系统消息，用来格式化chatgpt，设置它的预定行为
	UserRole      = "user"      // 用户的消息
)

type Body struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Usage struct {
	PromptToken      int32 `json:"prompt_tokens"`
	CompletionTokens int32 `json:"completion_tokens"`
	TotalTokens      int32 `json:"total_tokens"`
}

type Response struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   Usage  `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason interface{} `json:"finish_reason"`
		Index        int32       `json:"index"`
	} `json:"choices"`
}

// 实例化body
func NewBody(messagesContents []Message) Body {
	var body Body
	body.Model = "gpt-3.5-turbo"
	body.Messages = messagesContents
	return body
}

// 以不同的角色模式开启与chatgpt对话
type ChatMode struct {
	ModeName    string    // 模式名称
	PreMessages []Message // 预先设置好Chatgpt的角色
}

// 新建一个模式
func NewChatMode(modeName string, preMessages []Message) ChatMode {
	// TODO: 预先处理，比如说切换为全英文或者全中文
	return ChatMode{
		ModeName:    modeName,
		PreMessages: preMessages,
	}
}

type ChatGpt struct {
	ChatMode  ChatMode  // 初始化和缓存历史对象
	SecretKey string    // 密钥
	ProxyUrl  *url.URL  // 代理设置
	Messages  []Message // 缓存对话消息
}

// 实例化ChatGpt
func NewChatGptProxy(chatMode ChatMode, secretKey, proxyUrl string) *ChatGpt {
	var chat *ChatGpt
	if proxyUrl != "" {
		proxyUrlParsed, _ := url.Parse(proxyUrl)
		chat = &ChatGpt{
			ChatMode:  chatMode,
			SecretKey: secretKey,
			ProxyUrl:  proxyUrlParsed,
		}
	} else {
		chat = &ChatGpt{
			ChatMode:  chatMode,
			SecretKey: secretKey,
		}
	}
	chat.Messages = chatMode.PreMessages
	return chat
}

// 调取chatgpt接口
func (c *ChatGpt) postChatGpt() (*Response, error) {
	var client *http.Client
	if c.ProxyUrl != nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(c.ProxyUrl),
			},
		}
	} else {
		client = &http.Client{}
	}
	body := NewBody(c.Messages)
	bodyByte, _ := json.Marshal(body)
	var url string = "https://api.openai.com/v1/chat/completions"
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyByte))
	authorization := fmt.Sprintf("Bearer %s", c.SecretKey)
	request.Header.Set("Authorization", authorization)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var res Response
	err = json.Unmarshal(responseBodyBytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// 聊天
func (c *ChatGpt) Chat(message Message) (string, error) {
	// 1.将用户信息添加到message数组中
	c.Messages = append(c.Messages, message)
	// 2.调用chatgpt接口获取响应
	res, err := c.postChatGpt()
	if err != nil {
		return "", err
	}
	// 3.将响应消息放入message数组
	var responseMessage Message
	responseMessage.Role = AssistantRole
	responseMessage.Content = res.Choices[len(res.Choices)-1].Message.Content
	c.Messages = append(c.Messages, responseMessage)
	return responseMessage.Content, nil
}
