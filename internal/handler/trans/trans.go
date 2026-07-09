package trans

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/base"
	"develop_tools/internal/handler/render"
	"develop_tools/pkg/logger"
	"develop_tools/pkg/net"
)

func ConversionTrans(c *gin.Context) {
	render.Page(c, "conversion_trans.html")
}

type transResponse struct {
	Type            string
	ErrorCode       int
	ElapsedTime     int64
	TranslateResult [][]*translateResultItem
}

type translateResultItem struct {
	Src string
	Tgt string
}

func ConversionTransYoudao(c *gin.Context) {
	text := c.Query("t")
	if text == "" {
		base.ResponseError(c, "empty text")
		return
	}

	apiURL := "https://fanyi.youdao.com/translate?&doctype=json&type=AUTO&i=" + url.QueryEscape(text)
	httpResponse, body, err := net.HttpGetRequest(apiURL, net.HttpClient5000)
	if err != nil {
		logger.Error("fail to HttpGetRequest, err=%s", err.Error())
		base.ResponseError(c, "translate request failed")
		return
	}
	if httpResponse.StatusCode != http.StatusOK {
		base.ResponseError(c, fmt.Sprintf("translate upstream status %d", httpResponse.StatusCode))
		return
	}

	var payload transResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		base.ResponseError(c, fmt.Sprintf("fail to trans: %v", err))
		return
	}

	transType, transTgt := extractTranslation(&payload)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"result": gin.H{
			"transType": transType,
			"transTgt":  transTgt,
		},
	})
}

func extractTranslation(payload *transResponse) (string, string) {
	if payload.TranslateResult == nil || len(payload.TranslateResult) == 0 {
		return "", ""
	}
	if payload.TranslateResult[0] == nil || len(payload.TranslateResult[0]) == 0 {
		return "", ""
	}

	transType := payload.Type
	items := payload.TranslateResult[0]
	if len(items) == 1 {
		return transType, items[0].Tgt
	}

	var b strings.Builder
	for idx, item := range items {
		if idx > 0 {
			b.WriteString("\r\n")
		}
		fmt.Fprintf(&b, "%s\r\n%s", item.Src, item.Tgt)
	}
	return transType, b.String()
}
