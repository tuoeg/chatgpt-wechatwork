package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"server/config"
	"server/pkg/model"
	"server/pkg/util"
	"server/pkg/wxbizmsgcrypt"
	"strings"
)

func (s *service) GetToken() (string, error) {
	var tokenResponse model.TokenResponse
	corpid := config.NewConfig().WxConfig.CorpId
	corpsecret := config.NewConfig().WxConfig.CorpSecret
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
	if err != nil {
		return "", err
	}
	// 校验
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(config.NewConfig().WxConfig.Token, config.NewConfig().WxConfig.EncodingAeskey, config.NewConfig().WxConfig.CorpId, wxbizmsgcrypt.XmlType)
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

	// 异步发送消息到企业微信
	go s.sendMsgToWx(msgContent.FromUsername, msgContent.Content, accessToken)
	return "", nil
}

func (s *service) sendMsgToWx(user, msg, token string) {
	resMsg := model.MsgResponse{
		ToUser:  user,
		MsgType: "text",
		AgentId: "1000002",
		Text: model.Text{
			Content: "机器人正在思考...",
		},
		Safe:                   0,
		EnableIdTrans:          0,
		EnableDuplicateCheck:   0,
		DuplicateCheckInterval: 1800,
	}
	// 回传消息到企业微信
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
	req1, _, err := util.SendHTTPRequest(s.httpClient, url, "POST", resMsg, nil)
	if err != nil {
		return
	}
	log.Println(string(req1))
	var mrr model.MsgResponseResponse
	err = json.Unmarshal(req1, &mrr)
	if err != nil {
		return
	}
	go func(msgId string) {
		responseMsg, err := s.sendMsgToOpenAI(msg, user)
		if err != nil {
			log.Printf("Resquest openai error:%s", err.Error())
			// 撤回提示消息
			resMsg = model.MsgResponse{
				MsgId: msgId,
			}
			url = fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/recall?access_token=%s", token)
			req2, _, err := util.SendHTTPRequest(s.httpClient, url, "POST", resMsg, nil)
			if err != nil {
				return
			}
			log.Println(string(req2))

			resMsg := model.MsgResponse{
				ToUser:  user,
				MsgType: "text",
				AgentId: "1000002",
				Text: model.Text{
					Content: "遇到一点小故障，能否再输入一下问题？",
				},
				Safe:                   0,
				EnableIdTrans:          0,
				EnableDuplicateCheck:   0,
				DuplicateCheckInterval: 1800,
			}
			// 回传消息到企业微信
			url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
			req1, _, err := util.SendHTTPRequest(s.httpClient, url, "POST", resMsg, nil)
			if err != nil {
				return
			}
			log.Println(string(req1))
			return
		}
		// 撤回提示消息
		resMsg = model.MsgResponse{
			MsgId: msgId,
		}
		url = fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/recall?access_token=%s", token)
		req2, _, err := util.SendHTTPRequest(s.httpClient, url, "POST", resMsg, nil)
		if err != nil {
			return
		}
		log.Println(string(req2))

		resMsg = model.MsgResponse{
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
		url = fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
		req3, _, err := util.SendHTTPRequest(s.httpClient, url, "POST", resMsg, nil)
		if err != nil {
			return
		}
		log.Println(string(req3))
	}(mrr.MsgId)
}

func (s *service) sendMsgToOpenAI(msg, user string) (string, error) {
	var openAIResponse model.OpenAIResponse
	openAIRequest := model.OpenAIRequest{
		Model:    config.NewConfig().OpenAIConfig.Model,
		Messages: []model.Message{{Role: "user", Content: msg}},
	}
	req, _, err := util.SendHTTPRequest(s.httpClient, "https://api.openai.com/v1/chat/completions", "POST", openAIRequest, &util.HTTPOptions{Token: fmt.Sprintf("Bearer %s", config.NewConfig().OpenAIConfig.Token), ContentType: "application/json;charset=utf-8"})
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(req, &openAIResponse); err != nil {
		return "", err
	}
	content := strings.TrimLeft(openAIResponse.Choices[0].Message.Content, "\n")
	fmt.Printf("%s", content)
	return content, nil
}
