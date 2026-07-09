package share

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/base"
	"develop_tools/internal/model"
)

func SaveData(c *gin.Context) {
	type request struct {
		Id     int         `json:"id"`
		OpType int         `json:"opType"`
		Uuid   string      `json:"uuid"`
		Path   string      `json:"path"`
		Name   string      `json:"name"`
		Token  string      `json:"token"`
		Data   interface{} `json:"data"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	userID, err := resolveUserID(req.Uuid)
	if err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	if req.Id > 0 {
		if _, err := loadOwnedShare(req.Id, userID); err != nil {
			base.ResponseError(c, err.Error())
			return
		}
	}

	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	token := req.Token
	if req.Id == 0 {
		token, err = newShareToken()
		if err != nil {
			base.ResponseError(c, err.Error())
			return
		}
	}

	share := &model.Share{
		Id:     req.Id,
		OpType: req.OpType,
		UserId: userID,
		Uuid:   req.Uuid,
		Path:   req.Path,
		Name:   req.Name,
		Token:  token,
		Data:   string(dataBytes),
		Status: 0,
	}

	sm := model.NewShareModel()
	if req.Id > 0 {
		err = sm.Update(share)
	} else {
		err = sm.Save(share)
	}
	if err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	base.ResponseOk(c, gin.H{"id": share.Id, "token": share.Token})
}

func GetData(c *gin.Context) {
	type request struct {
		Path  string `json:"path"`
		Token string `json:"token"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	data := model.NewShareModel().GetShareData(req.Path, req.Token)
	if data == "" {
		base.ResponseError(c, "not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "ok",
		"data":    data,
	})
}

func List(c *gin.Context) {
	type request struct {
		OpType int    `json:"opType"`
		Uuid   string `json:"uuid"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	userID, err := resolveUserID(req.Uuid)
	if err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	shareList := model.NewShareModel().GetShareListByUserId(req.OpType, userID)
	grouped := make(map[string][]*model.Share)
	for _, item := range shareList {
		item.UpdatedDate = item.UpdatedAt.Format("2006-01-02 15:04:05")
		grouped[item.Path] = append(grouped[item.Path], item)
	}

	base.ResponseOk(c, grouped)
}

func Delete(c *gin.Context) {
	type request struct {
		Id   int    `json:"id"`
		Uuid string `json:"uuid"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	userID, err := resolveUserID(req.Uuid)
	if err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	share, err := loadOwnedShare(req.Id, userID)
	if err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	share.Status = 1
	if err := model.NewShareModel().Update(share); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	base.ResponseOk(c, nil)
}

func newShareToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
