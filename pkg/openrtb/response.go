package openrtb

type AdResponse struct {
	Id         string         `json:"id,omitempty"`
	SeatBid    []*SeatBid     `json:"seatbid,omitempty"`
	BidId      string         `json:"bidid,omitempty"`
	Cur        string         `json:"cur,omitempty"`
	CustomData string         `json:"customdata,omitempty"`
	Nbr        int            `json:"nbr,omitempty"`
	Ext        *AdResponseExt `json:"ext,omitempty"`
}

type SeatBid struct {
	Bid   []*Bid      `json:"bid,omitempty"`
	IsCn  int         `json:"iscn,omitempty"`
	BidCn *BidCn      `json:"bidcn,omitempty"`
	Seat  string      `json:"seat,omitempty"`
	Group int         `json:"grup,omitempty"`
	Ext   *SeatBidExt `json:"ext,omitempty"`
}

type SeatBidExt struct {
	DspCompanyName string `json:"dc"`
	DspAccountId   int    `json:"daid"`
	DspAccountName string `json:"da"`
}

type AdResponseExt struct {
	TP               *AdResponseExtTP          `json:"tp,omitempty"`
	EffectiveDisplay *EffectiveDisplay         `json:"effective_display"`          //有效展示配置
	RenderStyle      *AdResponseExtRenderStyle `json:"render_style"`               //APP渲染样式
	AutoRedirect     *AutoRedirect             `json:"auto_redirect"`              //自动跳转
	CnSplashConfig   *CnSplashConfig           `json:"cn_splash_config,omitempty"` //国内开屏配置
	IsExclusive      int                       `json:"is_exclusive"`               //是否独占
}

type AdResponseExtTP struct {
	AppId        int `json:"app_id"`
	AdseatId     int `json:"adseat_id"`
	SegmentId    int `json:"segment_id"`
	BucketId     int `json:"bucket_id"`
	AspId        int `json:"asp_id"`
	DspAccountId int `json:"dsp_account_id"`
	DspAdType    int `json:"dsp_ad_type"`
}

type EffectiveDisplay struct {
	CheckVisible bool `json:"check_visible"`  //View 是否可见
	MinAreaRatio int  `json:"min_area_ratio"` //最小广告展示面积比例
	MinDuration  int  `json:"min_duration"`   //最小广告展示持续时长
}

type AdResponseExtRenderStyle struct {
	//国内native配置
	TemplateRenderType int `json:"template_render_type"` //模板渲染类型 :| 1 左图右文 2 左文右图 3 上图下文 4 上文下图 5 三图 优先级: 3 > 4 > 1 > 2

	//海外video配置
	*RenderStyle
}

type AutoRedirect struct {
	FilterRatio int `json:"filter_ratio"` //过滤比例
}

type RenderStyle struct {
	VideoClickArea      int      `json:"video_click_area"`     //点击区域
	SkipBtnRatio        int      `json:"skip_btn_ratio"`       //跳过/关闭按钮大小
	VideoSkipTime       int      `json:"video_skip_time"`      //视频跳过出现时间
	EndcardCloseTime    int      `json:"endcard_close_time"`   //endcard关闭出现时间 (endcard1)
	IsEndcard2          int      `json:"is_endcard2"`          //是否endcard2类型， SDK取值：0 否 1 是 ； 后台取值：0，1 是  2 否
	Endcard2CloseTime   int      `json:"endcard2_close_time"`  //endcard2 关闭出现时间
	Endcard2IconUrl     string   `json:"endcard2_icon_url"`    //endcard2 icon
	Endcard2Title       string   `json:"endcard2_title"`       //endcard2 标题
	Endcard2ShowApp     int      `json:"endcard2_show_app"`    //endcard2 是否显示APP详情：0 否 1 是
	Endcard2Screenshots []string `json:"endcard2_screenshots"` //endcard2 APP详情截图列表
	CountdownStyle      int      `json:"countdown_style"`      //倒计时样式 1-数字，2-上进度条，3-下进度条
	CountdownColor      string   `json:"countdown_color"`      //倒计时颜色

	//不下发字段
	ClickArea int `json:"click_area,omitempty"` //点击区域
}

// 国内开屏配置
type CnSplashConfig struct {
	ClickType       int `json:"click_type"`
	SwayType        int `json:"sway_type"`
	ClickArea       int `json:"click_area"`
	SlideUpDistance int `json:"slide_up_distance"`
	SkipTime        int `json:"skip_time"`
}

type Bid struct {
	Id             string            `json:"id"`
	ImpId          string            `json:"impid"`
	Price          float64           `json:"price"`
	NUrl           string            `json:"nurl,omitempty"`
	BUrl           string            `json:"burl,omitempty"`
	LUrl           string            `json:"lurl,omitempty"`
	Adm            string            `json:"adm"`
	ADid           string            `json:"adid,omitempty"`
	ADomain        []string          `json:"adomain,omitempty"`
	Bundle         string            `json:"bundle,omitempty"`
	IUrl           string            `json:"iurl,omitempty"`
	CId            string            `json:"cid,omitempty"`
	CrId           string            `json:"crid,omitempty"`
	Cat            []string          `json:"cat,omitempty"`
	Attr           []int             `json:"attr,omitempty"`
	Api            int               `json:"api,omitempty"`
	Protocol       int               `json:"protocol,omitempty"`
	QagMediaRating int               `json:"qagmediarating,omitempty"`
	DealId         string            `json:"dealid,omitempty"`
	W              int               `json:"w,omitempty"`
	H              int               `json:"h,omitempty"`
	WRatio         int               `json:"wratio,omitempty"`
	HRatio         int               `json:"hratio,omitempty"`
	Exp            int64             `json:"exp,omitempty"`
	MType          int               `json:"mtype,omitempty"`
	Ext            *AdResponseBidExt `json:"ext,omitempty"`
}

type AdResponseBidExt struct {
	NUrl   []string `json:"nurl,omitempty"`
	LUrl   []string `json:"lurl,omitempty"`
	BUrl   []string `json:"burl,omitempty"`
	ImpUrl []string `json:"impurl,omitempty"`
	ClkUrl []string `json:"clkurl,omitempty"`

	Skadn *AdxBidExtSkadn `json:"skadn,omitempty"`

	Deeplink            string `json:"deeplink,omitempty"`
	DeeplinkFallBackUrl string `json:"deeplinkfallbackurl,omitempty"`
	Fallback            string `json:"fallback,omitempty"`

	Vxec           *int `json:"vxec,omitempty"`           //Indicates whether or not to add VX EC after the video or after the playable endcard.
	AutoStore      *int `json:"autostore,omitempty"`      //Indicates whether DTX should enable automatically displaying the Store 1 - true; 0 - false
	AutoStoreClick *int `json:"autostoreclick,omitempty"` //Indicates whether DT Exchange should fire click trackers when displaying the Store Kit: 1 - true; 0 - false

	PCta *BidExtPCta `json:"pcta,omitempty"`
}

type BidExtPCta struct {
	PCtaPage1 *int `json:"pctapage1,omitempty"`
}

type AdxBidExtSkadn struct {
	Version          string            `json:"version"`
	Network          string            `json:"network"`
	SourceIdentifier string            `json:"sourceidentifier,omitempty"`
	Campaign         string            `json:"campaign,omitempty"`
	Itunesitem       string            `json:"itunesitem"`
	ProductPageId    string            `json:"productpageid"`
	Sourceapp        string            `json:"sourceapp"`
	Fidelities       []*SkadnFidelitie `json:"fidelities,omitempty"`
	Nonce            string            `json:"nonce,omitempty"`
	Timestamp        string            `json:"timestamp,omitempty"`
	Signature        string            `json:"signature,omitempty"`
}

type SkadnFidelitie struct {
	Fidelity  int    `json:"fidelity"`
	Nonce     string `json:"nonce"`
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
}

// 验证格式是否正确
func (resp *AdResponse) VerifyFormat() bool {
	if len(resp.SeatBid) == 0 ||
		len(resp.SeatBid[0].Bid) == 0 ||
		len(resp.SeatBid[0].Bid[0].Adm) == 0 {
		return false
	}

	return true
}

func (bid *Bid) GetADomain() string {
	if len(bid.ADomain) > 0 {
		return bid.ADomain[0]
	}
	return ""
}
