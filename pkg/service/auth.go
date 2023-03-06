package service

import (
	"fmt"
	"server/pkg/wxbizmsgcrypt"
)

func (s *service) Auth(msgSignature, timestamp, nonce, echostr string) (string, error) {
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt("p0HrmcDC1goL", "yMA8ePfLF5ChUOSDL0ZynZzA99uOoxJxLb5fhwxffiS", "ww7f756c0495ad8cc4", wxbizmsgcrypt.XmlType)
	res, cryptErr := wxcpt.VerifyURL(msgSignature, timestamp, nonce, echostr)
	if cryptErr != nil {
		fmt.Println(cryptErr.ErrMsg)
		return "", fmt.Errorf("Error:%s", cryptErr.ErrMsg)
	}
	fmt.Println(string(res))
	return string(res), nil
}
