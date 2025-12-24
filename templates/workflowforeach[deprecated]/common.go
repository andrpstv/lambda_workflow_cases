package workflowforeachstep

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Input struct {
	GlobalID            string
	EpkID               string
	RemoteExecuteParams map[string]RemoteRequest
}
type RemoteRequest struct {
	RespSizeKb int  `json:"resp_size_kb"`
	Fail       bool `json:"fail"`
	DelaySec   int  `json:"delay_sec"`
}
type RemoteResponse struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
type Output struct {
	DataStore
	Models
	Strategy
}
type DataStore struct {
	FeatureStore string
	RiskAvatar   string
	RiskParams   string
}
type Models struct {
	Model1 string
	Model2 string
	Model3 string
	Model4 string
}
type Strategy struct {
	Strategy string
}

const (
	TaskQueueName   = "scenario1-case2"
	FeatureStoreURL = "http://localhost:8000/execute"
	RiskAvatarURL   = "http://localhost:8010/execute"
	RiskParamsURL   = "http://localhost:8020/execute"
	Model1URL       = "http://localhost:8030/execute"
	Model2URL       = "http://localhost:8040/execute"
	Model3URL       = "http://localhost:8050/execute"
	Model4URL       = "http://localhost:8060/execute"
	StrategyUrl        = "http://localhost:8070/execute"
)

func CallService(endPoint string, request []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(request))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed with errr: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}
	return body, nil
}
