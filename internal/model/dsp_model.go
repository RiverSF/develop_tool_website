package model

type Dsp struct {
	Id                  int     `gorm:"primary_key" json:"id"`
	Name                string  `json:"name"`
	UniqueKey           string  `json:"unique_key"`
	IsCn                int     `json:"is_cn"`
	RequestId           string  `json:"request_id"`
	Price               float64 `json:"price"`
	Adm                 string  `json:"adm"`
	Crid                string  `json:"crid"`
	Bundle              string  `json:"bundle"`
	Deeplink            string  `json:"deeplink"`
	Deeplinkfallbackurl string  `json:"deeplinkfallbackurl"`
	Fallback            string  `json:"fallback"`
}

func (m Dsp) TableName() string {
	return "my_dsp"
}

type DspModel struct {
}

func NewDspModel() *DspModel {
	return &DspModel{}
}

func (m *DspModel) GetDspList(isCn int) (dspList []*Dsp) {
	dspList = []*Dsp{}
	db.Where("is_cn = ?", isCn).Omit("adm").Order("updatetime desc").Find(&dspList)
	return
}

func (m *DspModel) GetDspAdmById(id int) (adm string) {
	_ = db.Table("my_dsp").Where("id = ?", id).Select("adm").Scan(&adm).Error
	return
}

func (m *DspModel) GetDspAdResponseById(id int) *Dsp {
	var dsp = &Dsp{}
	db.Where("id = ?", id).Find(&dsp)
	return dsp
}

func (m *DspModel) GetDspByUniqueKey(uniqueKey string) *Dsp {
	var dsp = &Dsp{}
	db.Table("my_dsp").Where("unique_key = ?", uniqueKey).First(&dsp)
	return dsp
}

func (m *DspModel) GetDspIDByUniqueKey(uniqueKey string) int {
	var id int
	_ = db.Table("my_dsp").Where("unique_key = ?", uniqueKey).Select("id").Scan(&id).Error
	return id
}

func (m *DspModel) Save(dsp *Dsp) error {
	return db.Save(dsp).Error
}
