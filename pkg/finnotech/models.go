package finnotech

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type BaseResponse struct {
	Status       string `json:"status"`
	TrackID      string `json:"trackId"`
	Error        *Error `json:"error"`
	ResponseCode string `json:"responseCode"`
}

type NegativeResult struct {
	LicenceNumber string  `json:"LicenceNumber"`
	NegativeScore string  `json:"NegativeScore"`
	OffenceCount  *string `json:"OffenceCount"`
	Rule          string  `json:"Rule"`
}

type NegativeScore struct {
	*BaseResponse
	Result *NegativeResult `json:"result"`
}

type DrivingOffenceBill struct {
	ID             string      `json:"id"`
	Type           string      `json:"type"`
	Description    string      `json:"description"`
	Code           string      `json:"code"`
	Price          int32       `json:"price"`
	City           *string     `json:"city"`
	Location       string      `json:"location"`
	Date           string      `json:"date"`
	PlateCode      string      `json:"serial"`
	DataValue      string      `json:"dataValue"`
	Barcode        interface{} `json:"barcode"`
	Licence        string      `json:"license"`
	BillID         int32       `json:"billId"`
	PaymentID      int32       `json:"paymentId"`
	NormalizedDate string      `json:"normalizedDate"`
	IsPayable      bool        `json:"isPayable"`
	PolicemanCode  interface{} `json:"policemanCode"`
	HasImage       bool        `json:"hasImage"`
}

type DrivingOffenceResult struct {
	Bills       []*DrivingOffenceBill `json:"Bills"`
	TotalAmount int32                 `json:"TotalAmount"`
}

type DrivingOffence struct {
	*BaseResponse
	Result *DrivingOffenceResult `json:"result"`
}

type Plate struct {
	Alphabet    string
	StartNumber int8
	EndNumber   int8
	RegionCode  int8
}
