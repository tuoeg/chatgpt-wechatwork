package model

type TokenResponse struct {
	ErrCode     int    `json:"errcode,omitempty"`
	ErrMsg      string `json:"errmsg,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

type Xml struct {
	ToUserName string `xml:"ToUserName"`
	AgentId    string `xml:"ToAgentID"`
	Encrypt    string `xml:"Encrypt"`
}

type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Msgid        string `xml:"MsgId"`
	Agentid      uint32 `xml:"AgentId"`
}

type MsgResponse struct {
	ToUser                 string `json:"touser,omitempty"`
	ToParty                string `json:"toparty,omitempty"`
	ToTag                  string `json:"totag,omitempty"`
	MsgType                string `json:"msgtype,omitempty"`
	AgentId                string `json:"agentid,omitempty"`
	Text                   Text   `json:"text,omitempty"`
	Safe                   int    `json:"safe,omitempty"`
	EnableIdTrans          int    `json:"enable_id_trans,omitempty"`
	EnableDuplicateCheck   int    `json:"enable_duplicate_check,omitempty"`
	DuplicateCheckInterval int    `json:"duplicate_check_interval,omitempty"`
}

type Text struct {
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string    `json:"model,omitempty"`
	Messages []Message `json:"messages"`
	// Prompt      string `json:"prompt,omitempty"`
	// MaxTokens   int    `json:"max_tokens,omitempty"`
	// Temperature int    `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type OpenAIResponse struct {
	Id      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	//Model   string   `json:"model,omitempty"`
	Choices []Choice `json:"choices,omitempty"`
}

type Choice struct {
	Index        int     `json:"index,omitempty"`
	Message      Message `json:"message,omitempty"`
	FinishReason string  `json:"finish_reason,omitempty"`
}
