package srv

// req
type Request struct {
	RespSizeKb int  `json:"resp_size_kb"`
	Fail       bool `json:"fail"`
	DelaySec   int  `json:"delay_sec"`
	IsEmptyResp bool `json:"is_empty_resp"`
}

// resp
type Response struct {
	Name    string `json:"name"`
	Request `json:"request"`
	Content string `json:"content"`
}
