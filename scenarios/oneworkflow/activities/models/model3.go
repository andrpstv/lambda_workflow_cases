package models

import (
	"context"
	"encoding/json"
	"lambda_workflow_cases/scenarios/oneworkflow"

	"go.temporal.io/sdk/activity"
)

func Model3Activity(ctx context.Context, input oneworkflow.Input) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Model3Activity started")

	data, ok := input.RemoteExecuteParams["Model3"]
	if !ok {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", ErrorParams
	}
	byteReq, err := json.Marshal(data)
	if err != nil {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", err
	}

	byteResp, err := oneworkflow.CallService(oneworkflow.Model3URL, byteReq)
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
	logger.Debug("Model3 activity completed.")

	return resp.Content, nil
}
