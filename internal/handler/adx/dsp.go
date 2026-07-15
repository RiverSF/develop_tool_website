package adx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/base"
	"develop_tools/internal/model"
	"develop_tools/pkg/logger"
	"develop_tools/pkg/openrtb"
)

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

	dspID := model.NewDspModel().GetDspIDByUniqueKey(uniqueKey)

	noticeTypeInt, err := model.NewDspNoticeModel().GetNoticeTypeValue(noticeType)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 401, "message": err.Error()})
		return
	}

	netIp := c.RemoteIP()
	ua := c.Request.Header.Get("User-Agent")

	dspNotice := &model.DspNotice{
		DspId:      dspID,
		NoticeType: noticeTypeInt,
		Ip:         netIp,
		Ua:         ua,
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
