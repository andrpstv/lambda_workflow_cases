package srv

// req
type Request struct {
	FeatureStore    string `json:"feature_store"`
	RespSizeKb int    `json:"resp_size_kb"`
	Fail       bool   `json:"fail"`
	DelaySec   int    `json:"delay_sec"`
}

// resp
type Response struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
