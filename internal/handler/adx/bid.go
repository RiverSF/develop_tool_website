package adx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/config"
	"develop_tools/internal/model"
	"develop_tools/pkg/common"
	"develop_tools/pkg/logger"
	"develop_tools/pkg/openrtb"
)

func AdxBid(c *gin.Context) {
	var adRequest = &openrtb.AdRequest{}

	c.ShouldBind(&adRequest)

	uniqueKey := c.Params.ByName("uniqueKey")

	dsp := model.NewDspModel().GetDspByUniqueKey(uniqueKey)

	adResponse := openrtb.NewAdResponse()

	adResponse.Id = common.If(len(dsp.RequestId) == 0 || dsp.RequestId == "{REQUEST_ID}", adRequest.Id, dsp.RequestId).(string)

	adResponse.SeatBid[0].Bid[0].Price = common.If(dsp.Price == 0, 99.99, dsp.Price).(float64)
	adResponse.SeatBid[0].Bid[0].Adm = dsp.Adm
	adResponse.SeatBid[0].Bid[0].Bundle = dsp.Bundle
	adResponse.SeatBid[0].Bid[0].CrId = dsp.Crid
	adResponse.SeatBid[0].Bid[0].Ext = &openrtb.AdResponseBidExt{
		Deeplink:            dsp.Deeplink,
		DeeplinkFallBackUrl: dsp.Deeplinkfallbackurl,
		Fallback:            dsp.Fallback,
	}

	var pcatpage = 1
	var AutoStore = 2
	var AutoStoreClick = 3
	adResponse.SeatBid[0].Bid[0].Ext.PCta = &openrtb.BidExtPCta{&pcatpage}
	adResponse.SeatBid[0].Bid[0].Ext.AutoStore = &AutoStore
	adResponse.SeatBid[0].Bid[0].Ext.AutoStoreClick = &AutoStoreClick

	adResponse.SeatBid[0].Bid[0].NUrl = config.AppConfig.LocalHost + "/adx/" + uniqueKey + "/nurl"
	adResponse.SeatBid[0].Bid[0].BUrl = config.AppConfig.LocalHost + "/adx/" + uniqueKey + "/burl"
	adResponse.SeatBid[0].Bid[0].LUrl = config.AppConfig.LocalHost + "/adx/" + uniqueKey + "/lurl"

	c.JSON(http.StatusOK, adResponse)
}

func AdxBidCn(c *gin.Context) {
	var adRequest = &openrtb.DspAccessCnRequest{}
	c.ShouldBind(&adRequest)

	requestId := adRequest.Id

	uniqueKey := c.Params.ByName("uniqueKey")

	dsp := model.NewDspModel().GetDspByUniqueKey(uniqueKey)

	var bidcn = &openrtb.BidCn{}
	if err := json.Unmarshal([]byte(dsp.Adm), &bidcn); err != nil {
		response := struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		}{
			201,
			fmt.Sprintf("解析失败：%s", err.Error()),
		}
		c.JSON(http.StatusOK, response)
		return
	}

	adResponse := &openrtb.DspAccessCnResponse{
		Id: requestId,
		SeatBid: []*openrtb.DspAccessCnSeatBid{
			{
				BidCn: []*openrtb.BidCn{
					bidcn,
				},
			},
		},
		BidId: "my bid id",
	}

	c.JSON(http.StatusOK, adResponse)
}

func AdxBidDsp(c *gin.Context) {
	var adRequest = &openrtb.DspAccessCnRequest{}
	c.ShouldBind(&adRequest)

	requestId := adRequest.Id
	//impId := "" todo：
	//if adRequest.

	uniqueKey := c.Params.ByName("uniqueKey")

	dsp := model.NewDspModel().GetDspByUniqueKey(uniqueKey)

	dsp.Adm = strings.ReplaceAll(dsp.Adm, "{RequestId}", requestId)
	dsp.Adm = strings.ReplaceAll(dsp.Adm, "{ImpId}", requestId)

	adResponse := dsp.Adm
	var adResponseJson map[string]interface{}
	_ = json.Unmarshal([]byte(adResponse), &adResponseJson)
	c.JSON(http.StatusOK, adResponseJson)
}

const (
	AdxUserIdMacro   = "{MY_DSP_UID}"
	GdprMacro        = "{MY_GDPR}"
	GdprConsentMacro = "{MY_GDPR_CONSENT}"
)

func AdxUserSync(c *gin.Context) {

	myGdpr := common.StringToInt(c.Query("my_gdpr"))
	myGdprConsent := c.Query("my_gdpr_consent")
	myPassthrough := c.Query("my_passthrough")

	myDspUid := c.Query("my_dsp_uid")
	if len(myDspUid) == 0 {
		myDspUid, _ = c.Cookie("my_uid")
	}

	if len(myDspUid) == 0 {
		myDspUid = fmt.Sprintf("dsp:%s", common.CreateUuid())
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "my_uid",
			Value:    myDspUid,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
	}

	logger.Info("【Adx UserSync】, requestURI=%s, myDspUid=%s", c.Request.RequestURI, myDspUid)

	redirectUrl := c.Query("my_redirect_url")

	if len(redirectUrl) > 0 {
		redirectUrl = strings.ReplaceAll(redirectUrl, AdxUserIdMacro, myDspUid)
		redirectUrl = strings.ReplaceAll(redirectUrl, GdprMacro, common.IntToString(myGdpr))
		redirectUrl = strings.ReplaceAll(redirectUrl, GdprConsentMacro, myGdprConsent)

		if len(myPassthrough) > 0 {
			redirectUrl += "&" + myPassthrough
		}

		safeURL, err := common.ValidateRedirectURL(redirectUrl, userSyncAllowedHosts()...)
		if err != nil {
			logger.Error("invalid redirect url: %s, err=%s", redirectUrl, err.Error())
			c.JSON(http.StatusOK, gin.H{"status": 400, "message": err.Error()})
			return
		}

		logger.Info("【Adx UserSync】, redirectUrl=%s", safeURL)
		c.Redirect(http.StatusFound, safeURL)
		return
	}

	c.JSON(http.StatusOK, nil)
}

var (
	userSyncHostsOnce sync.Once
	userSyncHosts     []string
)

func userSyncAllowedHosts() []string {
	userSyncHostsOnce.Do(func() {
		userSyncHosts = []string{"localhost", "127.0.0.1"}
		if host := common.HostnameFromURL(config.AppConfig.Host); host != "" {
			userSyncHosts = append(userSyncHosts, host)
		}
	})
	return userSyncHosts
}
