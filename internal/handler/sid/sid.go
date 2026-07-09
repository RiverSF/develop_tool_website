package sid

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/render"
	"develop_tools/pkg/common"
)

func Sid(c *gin.Context) { render.Page(c, "tp_sid.html") }

type sidStruct struct {
	Content string `json:"content"`
}

var domainRegex = regexp.MustCompile(`^[a-zA-Z0-9]+([\-\.]{1}[a-zA-Z0-9]+)*\.[a-zA-Z]{2,}$`)

func GetSid(c *gin.Context) {
	var rq sidStruct
	if err := c.ShouldBind(&rq); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "result": err.Error()})
		return
	}

	rq.Content = strings.Replace(rq.Content, " ", "", -1)
	rq.Content = strings.Replace(rq.Content, "'", "", -1)
	rq.Content = strings.Replace(rq.Content, "\"", "", -1)
	rq.Content = strings.Replace(rq.Content, ",", "", -1)
	rq.Content = strings.Replace(rq.Content, "\t", "", -1)

	records := strings.Split(rq.Content, "\n")

	content2 := ""
	for _, record := range records {
		sid, err := GetSupplyChainSid(record)
		if err != nil {
			content2 += fmt.Sprintf("%s       [error]%s\n", record, err.Error())
		} else {
			content2 += fmt.Sprintf("%s       river.com,%s,DIRECT\n", record, sid)
		}
	}

	response := struct {
		Status int    `json:"status"`
		Result string `json:"result"`
	}{
		200,
		content2,
	}
	c.JSON(http.StatusOK, response)
}

func GetSupplyChainSid(url string) (string, error) {
	mainDomain, err := GetHostMainDomain(url)
	if err != nil {
		return "", err
	}

	if sid, ok := sellidIdMap[mainDomain]; ok {
		return sid, err
	}

	if !domainRegex.MatchString(mainDomain) {
		return "", errors.New("域名格式验证错误")
	}
	domainSplitArr := strings.Split(mainDomain, ".")
	if len(domainSplitArr) >= 3 {
		domainSplitArrLast3 := domainSplitArr[len(domainSplitArr)-3:]
		if domainSplitArrLast3[0] == "www" {
			domainSplitArr = domainSplitArrLast3[1:]
		}
	}
	mainDomain = strings.Join(domainSplitArr, ".")
	sellerId := common.Md5("tp:" + mainDomain)
	return sellerId[:16], nil
}

// 域名解析主域名
func GetHostMainDomain(urlString string) (string, error) {
	if len(urlString) == 0 {
		return "", nil
	}

	if !strings.HasPrefix(urlString, "http") {
		urlString = "https://" + urlString
	}
	urlData, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	mainDomain := strings.ToLower(urlData.Host)

	if !domainRegex.MatchString(mainDomain) {
		return "", errors.New("域名格式验证错误")
	}

	if len(mainDomain) > 4 && strings.HasPrefix(mainDomain, "www.") {
		mainDomain = mainDomain[4:]
	}

	return mainDomain, nil
}

var sellidIdMap = map[string]string{
	"prinanetn.com":                              "1651-3918",
	"ikeepapps.com":                              "914391-9721",
	"gameapp.fly.dev":                            "18488-5066",
	"ludoroyal-online.com":                       "8374-10323",
	"services.paikangorchid.com":                 "6964-3407",
	"wisesweepmaster.xyz":                        "5312-5934",
	"morankj.cn":                                 "9431-4016",
	"allnovel.ltd":                               "65658-4604",
	"fy17wango.web.app":                          "634854-2182",
	"ohanaisland.com":                            "25863-3260",
	"ieasytech.com":                              "78622277-1804",
	"707interactiveplay.com":                     "537828-3659",
	"speedboosterclean.com":                      "667924-9868",
	"flashfilemanager.web.app":                   "6794-12339",
	"apkpure.com":                                "7977-12647",
	"chillyroom.com":                             "73878494-2105",
	"cynking.cn":                                 "5628-2630",
	"whitedotfun.github.io":                      "56839-9686",
	"bilibili.tv":                                "93882-10001",
	"dsby.yaojihuyu.cn":                          "67169-12325",
	"bibliaconsigo.com":                          "3783-8895",
	"umma.id":                                    "7716-11296",
	"heisky0816.firebaseapp.com":                 "76782-12787",
	"wallpaperkapp.com":                          "865367-13333",
	"yundr.xyz":                                  "743871417-6942",
	"filemanagertop.com":                         "1313-14103",
	"yuchaoplay.com":                             "94468-8314",
	"observation-studio.com":                     "9596-13550",
	"singlestar.voisky.com":                      "4451-13200",
	"perfectbrowser.net":                         "9647-14341",
	"mjyx.com":                                   "7566915-10568",
	"fancyrush.net":                              "7569-11961",
	"funfunquizearnmoney.com":                    "2352-14194",
	"f1rockets.com":                              "797-14880",
	"cleanerboosterx.com":                        "68-6928",
	"junglespeed.net":                            "99-14579",
	"gamefps.com":                                "258365-13039",
	"krumthematze.web.app":                       "6835-7845",
	"phonehdwallpaper.com":                       "5593-7180",
	"mygoldenwood.com":                           "962-15587",
	"vitabppro.com":                              "87276-14978",
	"yourhealthsentry.com":                       "35884-16035",
	"a.petll.eu.org":                             "3426357-16063",
	"coconutislandgames.com":                     "8527599-1692",
	"cp-us-east-1-06.s3.us-east-1.amazonaws.com": "4152-3036",
	"mhlptec.com":                                "8383-12451",
	"gstarcad.net":                               "91125-3757",
	"31gamestudio.com":                           "952-4009",
	"foodie.snow.me":                             "9673-4366",
	"mangatoon.mobi":                             "265-4674",
	"yoyogame.top":                               "952-4009",
	//"babybus.com":                                "7929182-5304",
	"fengzl.com":                   "9724229175-8811",
	"games.pujia8.com":             "18488-5066",
	"yomiko.group":                 "6937412-10743",
	"terabox.com":                  "8257-9455",
	"musicgamestudio.top":          "511913-10890",
	"dq.dianchu.com":               "7566915-10568",
	"edaysoft.cn":                  "511913-10890",
	"loklok.video":                 "92466-6186",
	"hy.haoyueqingfeng.com":        "42889188-11394",
	"dgames.mobi":                  "4249581-12318",
	"idailybible.com":              "914391-9721",
	"qixiang.zrbybs.cn":            "67169-12325",
	"k9616f99.app-ads-txt.com":     "6969-1965",
	"biblehelper.net":              "7569-11961",
	"winningstarters.github.io":    "18488-5066",
	"yoyaworld.com":                "56-13984",
	"buahanblast.com":              "2352-14194",
	"bibleinheart.com":             "914391-9721",
	"yodo1.com":                    "934362-12654",
	"dailykjvbibledevotion.com":    "5312-5934",
	"ad9o.yodo1.app":               "934362-12654",
	"loklok.ltd":                   "92466-6186",
	"watermelonemoji.com":          "5312-5934",
	"voiceeffectsoundmagicapp.com": "5312-5934",
	"coconut.is":                   "8527599-1692",
	"mini1.cn":                     "758874998-3050",
	"ttgstudio.com":                "537828-3659",
	"joyient.com":                  "537828-3659",
	"whitedot.fun":                 "56839-9686",
	"paperbride.link":              "8385-13641",
	"islgames.com":                 "258365-13039",
	"yearads.com":                  "84374-6207",
	"lulustudio.top":               "952-4009",
	"miniworldgame.com":            "758874998-3050",
	"play.google.com":              "934362-12654",
}
