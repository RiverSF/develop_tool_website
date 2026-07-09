package openrtb

func NewAdResponse() *AdResponse {
	adResponse := &AdResponse{
		Id: "",
		SeatBid: []*SeatBid{
			{
				Bid: []*Bid{
					{
						Id:             "",
						ImpId:          "",
						Price:          0,
						NUrl:           "",
						BUrl:           "",
						LUrl:           "",
						Adm:            "",
						ADid:           "",
						ADomain:        []string{},
						Bundle:         "",
						IUrl:           "",
						CId:            "",
						CrId:           "",
						Cat:            []string{},
						Attr:           []int{},
						Api:            0,
						Protocol:       0,
						QagMediaRating: 0,
						DealId:         "",
						W:              0,
						H:              0,
						WRatio:         0,
						HRatio:         0,
						Exp:            0,
						Ext:            &AdResponseBidExt{},
					},
				},
				Seat:  "",
				Group: 0,
				Ext:   &SeatBidExt{},
			},
		},
		BidId:      "",
		Cur:        "",
		CustomData: "",
		Nbr:        0,
		Ext:        &AdResponseExt{},
	}
	return adResponse
}
