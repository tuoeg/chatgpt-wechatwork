package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func (con *controller) Auth(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json,charset=utf-8")
	msgSignature, _ := url.QueryUnescape(c.Query("msg_signature"))
	timestamp, _ := url.QueryUnescape(c.Query("timestamp"))
	nonce, _ := url.QueryUnescape(c.Query("nonce"))
	echostr, _ := url.QueryUnescape(c.Query("echostr"))
	fmt.Println(msgSignature)
	fmt.Println(timestamp)
	fmt.Println(nonce)
	// 对字符串进行空格转换
	echostr_new := strings.ReplaceAll(echostr, " ", "+")
	fmt.Println(echostr_new)
	res, err := con.Service.Auth(msgSignature, timestamp, nonce, echostr_new)
	if err != nil {
		fmt.Println("err:", err)
		c.JSON(http.StatusBadRequest, fmt.Errorf("Error is err:%s", err.Error()))
		return
	}
	c.Writer.WriteString(res)
}
