package models

type StockTrackingInfo struct {
	BoardID                    string      `json:"boardId"`
	Isin                       string      `json:"isin"`
	AdminStatus                string      `json:"adminStatus"`
	CaStatus                   *string     `json:"caStatus"`
	Ceiling                    float64     `json:"ceiling"`
	CompanyNameEn              string      `json:"companyNameEn"`
	CompanyNameVi              string      `json:"companyNameVi"`
	CorporateEvents            []any       `json:"corporateEvents"`
	CouponRate                 float64     `json:"couponRate"`
	CoveredWarrantType         string      `json:"coveredWarrantType"`
	Exchange                   string      `json:"exchange"`
	ExercisePrice              float64     `json:"exercisePrice"`
	FirstTradingDate           string      `json:"firstTradingDate"`
	Floor                      float64     `json:"floor"`
	IssuerName                 string      `json:"issuerName"`
	LastTradingDate            string      `json:"lastTradingDate"`
	Market                     string      `json:"market"`
	MaturityDate               string      `json:"maturityDate"`
	ParValue                   int         `json:"parValue"`
	PermaHalt                  bool        `json:"permaHalt"`
	RefPrice                   float64     `json:"refPrice"`
	StockSymbol                string      `json:"stockSymbol"`
	StockType                  string      `json:"stockType"`
	TradingCurrencyISOCode     string      `json:"tradingCurrencyISOCode"`
	TradingDate                string      `json:"tradingDate"`
	TradingStatus              string      `json:"tradingStatus"`
	TradingUnit                int         `json:"tradingUnit"`
	ContractMultiplier         int         `json:"contractMultiplier"`
	PriorClosePrice            float64     `json:"priorClosePrice"`
	ProductID                  string      `json:"productId"`
	LastMFSeq                  int         `json:"lastMFSeq"`
	RemainForeignQtty          int         `json:"remainForeignQtty"`
	Best1Bid                   float64     `json:"best1Bid"`
	Best1BidVol                int         `json:"best1BidVol"`
	Best1Offer                 float64     `json:"best1Offer"`
	Best1OfferVol              int         `json:"best1OfferVol"`
	Best2Bid                   float64     `json:"best2Bid"`
	Best2BidVol                int         `json:"best2BidVol"`
	Best2Offer                 float64     `json:"best2Offer"`
	Best2OfferVol              int         `json:"best2OfferVol"`
	Best3Bid                   float64     `json:"best3Bid"`
	Best3BidVol                int         `json:"best3BidVol"`
	Best3Offer                 float64     `json:"best3Offer"`
	Best3OfferVol              int         `json:"best3OfferVol"`
	ExpectedLastUpdate         int64       `json:"expectedLastUpdate"`
	ExpectedMatchedPrice       float64     `json:"expectedMatchedPrice"`
	ExpectedMatchedVolume      int         `json:"expectedMatchedVolume"`
	ExpectedPriceChange        float64     `json:"expectedPriceChange"`
	ExpectedPriceChangePercent float64     `json:"expectedPriceChangePercent"`
	LastMESeq                  int         `json:"lastMESeq"`
	AvgPrice                   float64     `json:"avgPrice"`
	Highest                    float64     `json:"highest"`
	Lowest                     float64     `json:"lowest"`
	MatchedPrice               float64     `json:"matchedPrice"`
	MatchedVolume              int         `json:"matchedVolume"`
	NmTotalTradedQty           int         `json:"nmTotalTradedQty"`
	NmTotalTradedValue         int64       `json:"nmTotalTradedValue"`
	OpenPrice                  float64     `json:"openPrice"`
	PriceChange                float64     `json:"priceChange"`
	PriceChangePercent         float64     `json:"priceChangePercent"`
	StockBUVol                 int         `json:"stockBUVol"`
	StockVol                   int         `json:"stockVol"`
	StockSDVol                 int         `json:"stockSDVol"`
	BuyForeignQtty             int         `json:"buyForeignQtty"`
	BuyForeignValue            int64       `json:"buyForeignValue"`
	LastMTSeq                  int         `json:"lastMTSeq"`
	SellForeignQtty            int         `json:"sellForeignQtty"`
	SellForeignValue           int64       `json:"sellForeignValue"`
	Session                    string      `json:"session"`
	OddSession                 string      `json:"oddSession"`
	SessionPt                  string      `json:"sessionPt"`
	OddSessionPt               string      `json:"oddSessionPt"`
	SessionRt                  string      `json:"sessionRt"`
	OddSessionRt               string      `json:"oddSessionRt"`
	OddSessionRtStart          int64       `json:"oddSessionRtStart"`
	SessionRtStart             int64       `json:"sessionRtStart"`
	SessionStart               int64       `json:"sessionStart"`
	OddSessionStart            int64       `json:"oddSessionStart"`
	ExchangeSession            string      `json:"exchangeSession"`
	IsPreSessionPrice          bool        `json:"isPreSessionPrice"`
}

type StockListResponse struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Data    []StockTrackingInfo `json:"data"`
}

type StockInfoResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Data    StockTrackingInfo `json:"data"`
}
