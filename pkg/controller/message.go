package controller

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"server/pkg/model"

	"github.com/gin-gonic/gin"
)

func (con *controller) Message(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json,charset=utf-8")
	var msgXml model.Xml
	msgSignature, _ := url.QueryUnescape(c.Query("msg_signature"))
	timestamp, _ := url.QueryUnescape(c.Query("timestamp"))
	nonce, _ := url.QueryUnescape(c.Query("nonce"))
	fmt.Println(msgSignature, timestamp, nonce)
	b, _ := c.GetRawData()
	//fmt.Println(b)
	if err := xml.Unmarshal(b, &msgXml); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	res, _ := con.Service.SendMsg(msgSignature, timestamp, nonce, b)
	fmt.Println(res)
}
