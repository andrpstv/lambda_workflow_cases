package models

import (
	"context"
	"encoding/json"
	workflowforeachstep "lambda_workflow_cases/scenarios/workflowforeach"

	"go.temporal.io/sdk/activity"
)

func Model4Activity(ctx context.Context, input workflowforeachstep.Input) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Model4Activity started")

	data, ok := input.RemoteExecuteParams["Model4"]
	if !ok {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", ErrorParams
	}
	byteReq, err := json.Marshal(data)
	if err != nil {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", err
	}

	byteResp, err := workflowforeachstep.CallService(workflowforeachstep.Model4URL, byteReq)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
	var resp workflowforeachstep.RemoteResponse
	err = json.Unmarshal(byteResp, &resp)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
	logger.Debug("Model4 activity completed.")

	return resp.Content, nil
}
