package finnotech

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type BaseResponse struct {
	Status  string `json:"status"`
	TrackID string `json:"trackId"`
	Error   Error  `json:"error"`
}

type NegativeScore struct {
	BaseResponse
	ResponseCode string `json:"responseCode"`
	Result       struct {
		LicenceNumber string  `json:"LicenceNumber"`
		NegativeScore string  `json:"NegativeScore"`
		OffenceCount  *string `json:"OffenceCount"`
		Rule          string  `json:"Rule"`
	}
}
