package price

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/render"
	"develop_tools/pkg/encrypt"
)

func ConversionPrice(c *gin.Context) {
	render.Page(c, "tp_price.html")
}

func PriceEncrypt(c *gin.Context) {
	type Request struct {
		Type   int    `json:"type"`
		Price  string `json:"price"`
		Secret string `json:"secret"`
	}

	var req Request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "message": err.Error()})
		return
	}

	if b, _ := encrypt.VerifyEncryptKey(req.Secret); !b {
		c.JSON(http.StatusOK, gin.H{"status": 400, "message": "invalid secret"})
		return
	}

	var (
		priceStr string
		err      error
	)

	if req.Type == 1 {
		priceStr = encrypt.EcpmEncrypt(req.Price, req.Secret)
	} else {
		priceStr, err = encrypt.EcpmDecrypt(req.Price, req.Secret)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": 401, "message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "data": priceStr})
}
