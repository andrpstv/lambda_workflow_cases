package srv

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

const (
	SrvName = "Feature Store"
	ex      = "XXXXXXXXXX XXXXXXXXXX XXXXXXXXXX"
)

// healthcheck
func Health(w http.ResponseWriter, r *http.Request) {
	makeJSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// execute
func Execute(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sleep
	if req.DelaySec > 0 {
		time.Sleep(time.Second * time.Duration(req.DelaySec))
	}

	// 400
	if req.Fail {
		makeJSONResponse(w, http.StatusInternalServerError, map[string]string{":error": "fail"})
		return
	}

	// byte size
	targetSize := req.RespSizeKb * 1024

	// make string
	content := strings.Repeat("x", targetSize/65)

	resp := &Response{
		Name: SrvName,
		Request: Request{
			RespSizeKb: req.RespSizeKb,
			Fail:       req.Fail,
			DelaySec:   req.DelaySec,
		},
		Content: content,
	}

	makeJSONResponse(w, http.StatusOK, resp)
}

// http resp
func makeJSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
