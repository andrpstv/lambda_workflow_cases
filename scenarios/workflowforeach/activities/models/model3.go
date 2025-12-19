package models

import (
	"context"
	"encoding/json"
	workflowforeachstep "lambda_workflow_cases/scenarios/workflowforeach"

	"go.temporal.io/sdk/activity"
)

func Model3Activity(ctx context.Context, input workflowforeachstep.Input) (string, error) {
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

	byteResp, err := workflowforeachstep.CallService(workflowforeachstep.Model3URL, byteReq)
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
	logger.Debug("Model3 activity completed.")

	return resp.Content, nil
}
