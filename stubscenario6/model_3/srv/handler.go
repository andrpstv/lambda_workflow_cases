package srv

import (
	"encoding/json"
	"net/http"
)

const (
	SrvName = "Model 3"
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

	// // sleep
	// if req.DelaySec > 0 {
	// 	time.Sleep(time.Second * time.Duration(req.DelaySec))
	// }

	// // 400
	// if req.Fail {
	// 	makeJSONResponse(w, http.StatusInternalServerError, map[string]string{":error": "fail"})
	// 	return
	// }

	resp := &Response{
		Name:    SrvName,
		Content: req.FeatureStore + " checked by " + SrvName + " Service.",
	}

	makeJSONResponse(w, http.StatusOK, resp)
}

// http resp
func makeJSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
