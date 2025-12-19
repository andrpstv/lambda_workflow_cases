package strategies

import (
	"context"
	"encoding/json"
	"errors"
	"lambda_workflow_cases/scenarios/oneworkflow"

	"go.temporal.io/sdk/activity"
)
var ErrorParams = errors.New("params not found")

func CvUclActivity(ctx context.Context, input oneworkflow.Input) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CvUcl Activity started")

	data, ok := input.RemoteExecuteParams["CvUcl"]
	if !ok {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", ErrorParams
	}
	byteReq, err := json.Marshal(data)
	if err != nil {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", err
	}

	byteResp, err := oneworkflow.CallService(oneworkflow.CvUclUrl, byteReq)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
	var resp oneworkflow.RemoteResponse
	err = json.Unmarshal(byteResp, &resp)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
	logger.Debug("CvUcl activity completed.")

	return resp.Content, nil
}
