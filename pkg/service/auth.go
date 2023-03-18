package service

import (
	"fmt"
	"server/config"
	"server/pkg/wxbizmsgcrypt"
)

func (s *service) Auth(msgSignature, timestamp, nonce, echostr string) (string, error) {
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(config.NewConfig().WxConfig.Token, config.NewConfig().WxConfig.EncodingAeskey, config.NewConfig().WxConfig.CorpId, wxbizmsgcrypt.XmlType)
	res, cryptErr := wxcpt.VerifyURL(msgSignature, timestamp, nonce, echostr)
	if cryptErr != nil {
		fmt.Println(cryptErr.ErrMsg)
		return "", fmt.Errorf("Error:%s", cryptErr.ErrMsg)
	}
	return string(res), nil
}
