package router

import (
	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/adx"
	"develop_tools/internal/handler/share"
	"develop_tools/internal/handler/tools"
	"develop_tools/internal/handler/user"
	"develop_tools/internal/middleware"
	"develop_tools/pkg/path"
)

// New creates the HTTP engine with middleware, routes, and static assets.
func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	registerRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Static("/assets", path.Join("assets"))
	return r
}

func registerRoutes(r *gin.Engine) {
	root := r.Group("/")
	{
		root.GET("", tools.Home)
		root.GET("/timestamp", tools.ConversionTimestamp)
		root.GET("/json", tools.ConversionJson)
		root.GET("/html", tools.ConversionHtml)
		root.GET("/md5", tools.ConversionMd5)
		root.GET("/url", tools.ConversionUrl)
		root.GET("/base64", tools.ConversionBase64)
		root.GET("/utf8", tools.ConversionUtf8)
		root.GET("/unicode", tools.ConversionUnicode)
		root.GET("/markdown", tools.MarkDown)

		root.GET("/data-diff", tools.DataDiff)
		root.GET("/data-editor", tools.DataEditor)
		root.GET("/data-filter", tools.DataFilter)
		root.GET("/data-combine", tools.DataCombine)
		root.GET("/data-calculator", tools.DataCalculator)
		root.POST("/data-calc", tools.DataCalc)
		root.GET("/data-unique", tools.DataUnique)
		root.POST("/data-unique-exec", tools.DataUniqueExec)
		root.GET("/data-interval", tools.DataInterval)
		root.GET("/data-line", tools.DataLine)

		root.GET("/aes", tools.ConversionAes)
		root.GET("/trans", tools.ConversionTrans)
		root.GET("/trans/youdao", tools.ConversionTransYoudao)

		root.GET("/price", tools.ConversionPrice)
		root.POST("/price/price-encrypt", tools.PriceEncrypt)

		root.GET("/sid", tools.Sid)
		root.POST("/sid/get", tools.GetSid)

		root.GET("/bundle", adx.BundleIndex)

		root.GET("/gepip", tools.GeoIp)
		root.GET("/chart-sankey", tools.ChartSankey)
		root.GET("/timeline", tools.Timeline)
		root.GET("/map", tools.Map)
		root.GET("/realtime-translation", tools.RealtimeTranslation)

		root.GET("/dev/json2gostruct", tools.DevJson2GoStruct)
		root.POST("/dev/json2gostruct", tools.DevJson2GoStructApi)
	}

	adxGroup := r.Group("/adx")
	{
		// pages
		adxGroup.GET("", adx.AdxIndex)
		adxGroup.GET("/cn", adx.AdxCn)
		adxGroup.GET("/dsp", adx.AdxDSP)
		adxGroup.GET("/adm", adx.AdxAdm)

		// DSP CRUD + notice
		adxGroup.POST("/adxGetDspList", adx.AdxGetDspList)
		adxGroup.GET("/adxGetDspAdm", adx.AdxGetDspAdm)
		adxGroup.GET("/adxGetDspResponse", adx.AdxGetDspResponse)
		adxGroup.POST("/adxDspSave", adx.AdxDspSave)
		adxGroup.GET("/adxGetDspNotice", adx.AdxGetDspNotice)

		adxGroup.POST("/bundle/extract", adx.Extract)

		// bid（notice 回调与 param 路由保持原注册顺序）
		adxGroup.POST("/cn/:uniqueKey", adx.AdxBidCn)
		adxGroup.GET("/:uniqueKey/:noticeType", adx.AdxSaveDspNotice)
		adxGroup.POST("/:uniqueKey", adx.AdxBid)
		adxGroup.POST("/dsp/:uniqueKey", adx.AdxBidDsp)
		adxGroup.GET("/userSync", adx.AdxUserSync)
	}

	shareGroup := r.Group("/share")
	{
		shareGroup.POST("/saveData", share.SaveData)
		shareGroup.POST("/getData", share.GetData)
		shareGroup.POST("/list", share.List)
		shareGroup.POST("/delete", share.Delete)
	}

	userGroup := r.Group("/user")
	{
		userGroup.GET("", user.User)
		userGroup.POST("/get", user.UserGet)
		userGroup.POST("/save", user.UserSave)
	}
}
