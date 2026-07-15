package tools

import (
	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/render"
)

func ConversionJson(c *gin.Context)      { render.Page(c, "conversion_json.html") }
func ConversionHtml(c *gin.Context)      { render.Page(c, "conversion_html.html") }
func ConversionXml(c *gin.Context)       { render.Page(c, "conversion_xml.html") }
func ConversionTimestamp(c *gin.Context) { render.Page(c, "conversion_timestamp.html") }
func ConversionMd5(c *gin.Context)       { render.Page(c, "conversion_md5.html") }
func ConversionUrl(c *gin.Context)       { render.Page(c, "conversion_url.html") }
func ConversionBase64(c *gin.Context)    { render.Page(c, "conversion_base64.html") }
func ConversionUtf8(c *gin.Context)      { render.Page(c, "conversion_utf8.html") }
func ConversionUnicode(c *gin.Context)   { render.Page(c, "conversion_unicode.html") }
func ConversionAes(c *gin.Context)       { render.Page(c, "tp_aes.html") }
func MarkDown(c *gin.Context)            { render.Page(c, "markdown.html") }
func Timeline(c *gin.Context)            { render.Page(c, "timeline.html") }
func GeoIp(c *gin.Context)               { render.Page(c, "tp_ip.htm") }
func Map(c *gin.Context)                 { render.Page(c, "baidu_map.html") }
func ChartSankey(c *gin.Context)         { render.Page(c, "chart_sankey.html") }
func RealtimeTranslation(c *gin.Context) { render.Page(c, "realtime_translation.html") }
