package datastore

import (
	"context"
	"encoding/json"
	"lambda_workflow_cases/templates/oneworkflow"

	"go.temporal.io/sdk/activity"
)

func RiskAvatarActivity(ctx context.Context, input oneworkflow.Input) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("RiskAvatarActivity started")

	data, ok := input.RemoteExecuteParams["RiskAvatar"]
	if !ok {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", ErrorParams
	}
	byteReq, err := json.Marshal(data)
	if err != nil {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", err
	}

	byteResp, err := oneworkflow.CallService(oneworkflow.RiskAvatarURL, byteReq)
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
	logger.Debug("Risk Avatar activity completed.")

	return resp.Content, nil
}
