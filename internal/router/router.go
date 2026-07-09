package router

import (
	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/adx"
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

	g0 := r.Group("/")
	{
		g0.GET("", home.Home)
		g0.GET("/timestamp", conversion.ConversionTimestamp)
		g0.GET("/json", conversion.ConversionJson)
		g0.GET("/html", conversion.ConversionHtml)
		g0.GET("/md5", conversion.ConversionMd5)
		g0.GET("/url", conversion.ConversionUrl)
		g0.GET("/base64", conversion.ConversionBase64)
		g0.GET("/utf8", conversion.ConversionUtf8)
		g0.GET("/unicode", conversion.ConversionUnicode)
		g0.GET("/markdown", conversion.MarkDown)

		g0.GET("/data-diff", data.DataDiff)
		g0.GET("/data-editor", data.DataEditor)
		g0.GET("/data-filter", data.DataFilter)
		g0.GET("/data-combine", data.DataCombine)
		g0.GET("/data-calculator", data.DataCalculcator)
		g0.POST("/data-calc", data.DataCalc)
		g0.GET("/data-unique", data.DataUnique)
		g0.POST("/data-unique-exec", data.DataUniqueExec)
		g0.GET("/data-interval", data.DataInterval)
		g0.GET("/data-line", data.DataLine)

		g0.GET("/aes", conversion.ConversionAes)
		g0.GET("/trans", trans.ConversionTrans)
		g0.GET("/trans/youdao", trans.ConversionTransYoudao)

		g0.GET("/price", price.ConversionPrice)
		g0.POST("/price/price-encrypt", price.PriceEncrypt)

		g0.GET("/sid", sid.Sid)
		g0.POST("/sid/get", sid.GetSid)

		g0.GET("/gepip", conversion.GeoIp)
		g0.GET("/chart-sankey", conversion.ChartSankey)
		g0.GET("/timeline", conversion.Timeline)
		g0.GET("/map", conversion.Map)
		g0.GET("/realtime-translation", conversion.RealtimeTranslation)

		g0.GET("/dev/json2gostruct", dev.DevJson2GoStruct)
		g0.POST("/dev/json2gostruct", dev.DevJson2GoStructApi)
	}

	g5 := r.Group("/adx")
	{
		g5.GET("", adx.AdxIndex)
		g5.GET("/cn", adx.AdxCn)
		g5.GET("/dsp", adx.AdxDSP)
		g5.GET("/adm", adx.AdxAdm)

		g5.POST("/adxGetDspList", adx.AdxGetDspList)
		g5.GET("/adxGetDspAdm", adx.AdxGetDspAdm)
		g5.GET("/adxGetDspResponse", adx.AdxGetDspResponse)

		g5.POST("/adxDspSave", adx.AdxDspSave)
		g5.GET("/adxGetDspNotice", adx.AdxGetDspNotice)

		g5.POST("/cn/:uniqueKey", adx.AdxBidCn)
		g5.GET("/:uniqueKey/:noticeType", adx.AdxSaveDspNotice)
		g5.POST("/:uniqueKey", adx.AdxBid)

		g5.POST("/dsp/:uniqueKey", adx.AdxBidDsp)
		g5.GET("/userSync", adx.AdxUserSync)
	}

	g8 := r.Group("/share")
	{
		g8.POST("/saveData", share.SaveData)
		g8.POST("/getData", share.GetData)
		g8.POST("/list", share.List)
		g8.POST("/delete", share.Delete)
	}

	g11 := r.Group("/user")
	{
		g11.GET("", user.User)
		g11.POST("/get", user.UserGet)
		g11.POST("/save", user.UserSave)
	}
}
