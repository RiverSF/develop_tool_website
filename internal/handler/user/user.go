package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/base"
	"develop_tools/internal/handler/render"
	"develop_tools/internal/model"
	"develop_tools/pkg/common"
	"develop_tools/pkg/logger"
)

func User(c *gin.Context) {
	render.Page(c, "user.html")
}

func UserGet(c *gin.Context) {
	type request struct {
		Uuid string `json:"uuid"`
	}
	var req request
	if err := c.ShouldBind(&req); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	userKey := model.NewUserKeyModel().GetUserIdByKey(req.Uuid)

	userName := ""
	if userKey.UserId > 0 {
		user := model.NewUserModel().GetUser(userKey.UserId)
		userName = user.Name
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"message":  "ok",
		"userName": userName,
	})
}

func UserSave(c *gin.Context) {
	type request struct {
		Name string `json:"name"`
		Uuid string `json:"uuid"`
	}
	var req request
	if err := c.ShouldBind(&req); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	user, err := model.NewUserModel().FindOrCreateByName(req.Name)
	if err != nil {
		logger.Error("fail to save user, err=%s", err.Error())
		base.ResponseError(c, err.Error())
		return
	}

	userKey := &model.UserKey{
		UserId:     user.Id,
		BrowserKey: req.Uuid,
		UserAgent:  c.Request.Header.Get("user-agent"),
	}
	if err := model.NewUserKeyModel().UpsertByBrowserKey(userKey); err != nil {
		logger.Error("fail to save user key, err=%s", err.Error())
		base.ResponseError(c, err.Error())
		return
	}

	userKeyList := model.NewUserKeyModel().GetUserKeyList(user.Id)

	userKeys := []string{}
	userAgents := []string{}
	for _, item := range userKeyList {
		userKeys = append(userKeys, item.BrowserKey)
		if !common.IsContainStr(userAgents, item.UserAgent) {
			userAgents = append(userAgents, item.UserAgent)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     200,
		"message":    "ok",
		"userKeys":   userKeys,
		"userAgents": userAgents,
		"userId":     user.Id,
	})
}
