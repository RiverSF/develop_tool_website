package adx

import (
	"develop_tools/internal/config"
	"develop_tools/internal/handler/base"
	"develop_tools/internal/handler/render"
	"develop_tools/internal/model"
	"develop_tools/pkg/common"
	"develop_tools/pkg/logger"
	"develop_tools/pkg/openrtb"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AdxIndex(c *gin.Context) {
	render.Page(c, "tp_adx.html")
}

func AdxCn(c *gin.Context) {
	render.Page(c, "tp_adx_cn.html")
}

func AdxDSP(c *gin.Context) {
	render.Page(c, "tp_adx_dsp.html")
}

func AdxAdm(c *gin.Context) {
	render.Page(c, "tp_adx_adm.html")
}

func AdxDspSave(c *gin.Context) {
	dsp := &model.Dsp{}
	if err := c.ShouldBind(dsp); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	if dsp.IsCn == 1 {
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
	}

	pw := model.NewDspModel()
	if err := pw.Save(dsp); err != nil {
		logger.Error("fail to save dsp, err=%s", err.Error())
		base.ResponseError(c, err.Error())
		return
	}

	base.ResponseOk(c, nil)
}

func AdxGetDspList(c *gin.Context) {
	type request struct {
		IsCn int `json:"is_cn"`
	}
	var req request
	if err := c.ShouldBind(&req); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	dspList := model.NewDspModel().GetDspList(req.IsCn)

	for _, dsp := range dspList {
		dsp.Adm = ""
	}

	c.JSON(http.StatusOK, dspList)
}

func AdxGetDspResponse(c *gin.Context) {

	id, _ := strconv.Atoi(c.Query("id"))

	dsp := model.NewDspModel().GetDspAdResponseById(id)

	base.ResponseOk(c, dsp)
}

func AdxGetDspAdm(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	adm := model.NewDspModel().GetDspAdmById(id)
	c.JSON(http.StatusOK, adm)
}

func AdxSaveDspNotice(c *gin.Context) {
	uniqueKey := c.Params.ByName("uniqueKey")
	noticeType := c.Params.ByName("noticeType")

	dsp := model.NewDspModel().GetDspByUniqueKey(uniqueKey)

	noticeTypeInt, err := model.NewDspNoticeModel().GetNoticeTypeValue(noticeType)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 401, "message": err.Error()})
		return
	}

	netIp := c.RemoteIP()
	ua := c.Request.Header.Get("User-Agent")

	dspNotice := &model.DspNotice{
		DspId:      dsp.Id,
		NoticeType: noticeTypeInt,
		Ip:         netIp,
		Ua:         ua,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = model.NewDspNoticeModel().Save(dspNotice)

	if err != nil {
		logger.Error("fail to save dspNotice, err=%s", err.Error())
	}

	response := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		200,
		"ok",
	}

	c.JSON(http.StatusOK, response)
}

func AdxGetDspNotice(c *gin.Context) {

	dspId, _ := strconv.Atoi(c.Query("id"))
	noticeId, _ := strconv.Atoi(c.Query("noticeId"))
	nowDate := c.Query("nowDate")

	dspNoticeList := model.NewDspNoticeModel().GetDspNoticeByDspId(dspId, noticeId, nowDate)

	type dataIterm struct {
		*model.DspNotice

		NoticeTypeName string `json:"notice_type_name"`
	}

	var data = []*dataIterm{}
	for _, dspNotice := range dspNoticeList {

		//dspNotice.CreateTime =

		var dtIterm = &dataIterm{
			DspNotice: dspNotice,
		}

		noticeTypeName, err := model.NewDspNoticeModel().GetNoticeTypeName(dspNotice.NoticeType)
		if err != nil {
			logger.Error("fail to get notice type name, err=%s", err.Error())
			continue
		}

		dtIterm.NoticeTypeName = noticeTypeName

		data = append(data, dtIterm)
	}

	c.JSON(http.StatusOK, data)
}

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
	json.Unmarshal([]byte(adResponse), &adResponseJson)
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
		redirectUrl = strings.Replace(redirectUrl, AdxUserIdMacro, myDspUid, -1)
		redirectUrl = strings.Replace(redirectUrl, GdprMacro, common.IntToString(myGdpr), -1)
		redirectUrl = strings.Replace(redirectUrl, GdprConsentMacro, myGdprConsent, -1)

		if len(myPassthrough) > 0 {
			redirectUrl += "&" + myPassthrough
		}

		allowedHosts := []string{"localhost", "127.0.0.1"}
		if host := common.HostnameFromURL(config.AppConfig.Host); host != "" {
			allowedHosts = append(allowedHosts, host)
		}
		safeURL, err := common.ValidateRedirectURL(redirectUrl, allowedHosts...)
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
