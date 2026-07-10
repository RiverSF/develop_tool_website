package router

import (
	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/adx"
	"develop_tools/internal/handler/bundle"
	"develop_tools/internal/handler/conversion"
	"develop_tools/internal/handler/data"
	"develop_tools/internal/handler/dev"
	"develop_tools/internal/handler/home"
	"develop_tools/internal/handler/price"
	"develop_tools/internal/handler/share"
	"develop_tools/internal/handler/sid"
	"develop_tools/internal/handler/trans"
	"develop_tools/internal/handler/user"
	"develop_tools/internal/middleware"
)

func Init(r *gin.Engine) {
	initRouter(r)
}

func initRouter(r *gin.Engine) {
	r.Use(middleware.LoggerMiddleware())

	root := r.Group("/")
	{
		root.GET("", home.Home)
		root.GET("/timestamp", conversion.ConversionTimestamp)
		root.GET("/json", conversion.ConversionJson)
		root.GET("/html", conversion.ConversionHtml)
		root.GET("/md5", conversion.ConversionMd5)
		root.GET("/url", conversion.ConversionUrl)
		root.GET("/base64", conversion.ConversionBase64)
		root.GET("/utf8", conversion.ConversionUtf8)
		root.GET("/unicode", conversion.ConversionUnicode)
		root.GET("/markdown", conversion.MarkDown)

		root.GET("/data-diff", data.DataDiff)
		root.GET("/data-editor", data.DataEditor)
		root.GET("/data-filter", data.DataFilter)
		root.GET("/data-combine", data.DataCombine)
		root.GET("/data-calculator", data.DataCalculcator)
		root.POST("/data-calc", data.DataCalc)
		root.GET("/data-unique", data.DataUnique)
		root.POST("/data-unique-exec", data.DataUniqueExec)
		root.GET("/data-interval", data.DataInterval)
		root.GET("/data-line", data.DataLine)

		root.GET("/aes", conversion.ConversionAes)
		root.GET("/trans", trans.ConversionTrans)
		root.GET("/trans/youdao", trans.ConversionTransYoudao)

		root.GET("/price", price.ConversionPrice)
		root.POST("/price/price-encrypt", price.PriceEncrypt)

		root.GET("/sid", sid.Sid)
		root.POST("/sid/get", sid.GetSid)

		root.GET("/bundle", bundle.BundleIndex)

		root.GET("/gepip", conversion.GeoIp)
		root.GET("/chart-sankey", conversion.ChartSankey)
		root.GET("/timeline", conversion.Timeline)
		root.GET("/map", conversion.Map)
		root.GET("/realtime-translation", conversion.RealtimeTranslation)

		root.GET("/dev/json2gostruct", dev.DevJson2GoStruct)
		root.POST("/dev/json2gostruct", dev.DevJson2GoStructApi)
	}

	adxGroup := r.Group("/adx")
	{
		adxGroup.GET("", adx.AdxIndex)
		adxGroup.GET("/cn", adx.AdxCn)
		adxGroup.GET("/dsp", adx.AdxDSP)
		adxGroup.GET("/adm", adx.AdxAdm)

		adxGroup.POST("/adxGetDspList", adx.AdxGetDspList)
		adxGroup.GET("/adxGetDspAdm", adx.AdxGetDspAdm)
		adxGroup.GET("/adxGetDspResponse", adx.AdxGetDspResponse)

		adxGroup.POST("/adxDspSave", adx.AdxDspSave)
		adxGroup.GET("/adxGetDspNotice", adx.AdxGetDspNotice)

		adxGroup.POST("/bundle/extract", bundle.Extract)

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
