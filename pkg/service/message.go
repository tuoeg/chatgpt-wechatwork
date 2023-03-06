package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"server/pkg/model"
	"server/pkg/util"
	"server/pkg/wxbizmsgcrypt"
)

func (s *service) GetToken() (string, error) {
	var tokenResponse model.TokenResponse
	corpid := "ww7f756c0495ad8cc4"
	corpsecret := "YFY0gcQVQQOHSKex7PIPpR58fvtjiXjrCYTlEyKgHw8"
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpid, corpsecret)
	req, _, err := util.SendHTTPRequest(s.httpClient, url, "GET", nil, nil)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(req, &tokenResponse); err != nil {
		return "", err
	}
	return tokenResponse.AccessToken, nil
}

func (s *service) SendMsg(msgSignature, timestamp, nonce string, msg []byte) (string, error) {
	accessToken, err := s.GetToken()
	fmt.Println(accessToken)
	if err != nil {
		return "", err
	}
	// 校验
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt("p0HrmcDC1goL", "yMA8ePfLF5ChUOSDL0ZynZzA99uOoxJxLb5fhwxffiS", "ww7f756c0495ad8cc4", wxbizmsgcrypt.XmlType)
	msg, cryptErr := wxcpt.DecryptMsg(msgSignature, timestamp, nonce, msg)
	if cryptErr != nil {
		return "", fmt.Errorf("Decrypt msg error:%s", cryptErr.ErrMsg)
	}
	// 解析出明文
	var msgContent model.MsgContent
	err = xml.Unmarshal(msg, &msgContent)
	if nil != err {
		return "", err
	}
	// openai解析消息

	// 异步发送消息到企业微信
	go s.sendMsgToWx(msgContent.FromUsername, msgContent.Content, accessToken)
	return "", nil
}

func (s *service) sendMsgToWx(user, msg, token string) {
	responseMsg, err := s.sendMsgToOpenAI(msg, user)
	if err != nil {
		log.Printf("Resquest openai error:%s", err.Error())
		return
	}
	resMsg := model.MsgResponse{
		ToUser:  user,
		MsgType: "text",
		AgentId: "1000002",
		Text: model.Text{
			Content: responseMsg,
		},
		Safe:                   0,
		EnableIdTrans:          0,
		EnableDuplicateCheck:   0,
		DuplicateCheckInterval: 1800,
	}
	// 回传消息到企业微信
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
	req, _, err := util.SendHTTPRequest(s.httpClient, url, "POST", resMsg, nil)
	if err != nil {
		return
	}
	log.Println(req)
}

func (s *service) sendMsgToOpenAI(msg, user string) (string, error) {
	var openAIResponse model.OpenAIResponse
	openAIRequest := model.OpenAIRequest{
		Model:    "text-davinci-003",
		Messages: []model.Message{{Role: user, Content: msg}},
	}
	req, _, err := util.SendHTTPRequest(s.httpClient, "https://api.openai.com/v1/chat/completions", "POST", openAIRequest, &util.HTTPOptions{Token: "Bearer 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", ContentType: "application/json;charset=utf-8"})
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(req, &openAIResponse); err != nil {
		return "", err
	}
	return openAIResponse.Choices[0].Message.Content, nil
}
