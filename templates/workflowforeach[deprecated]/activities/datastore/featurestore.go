package datastore

import (
	"context"
	"encoding/json"
	"errors"
	workflowforeachstep "lambda_workflow_cases/scenarios/workflowforeach"

	"go.temporal.io/sdk/activity"
)

var ErrorParams = errors.New("params not found")

func FeatureStoreActivity(ctx context.Context, input workflowforeachstep.Input) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("FeatureStoreActivity started")

	data, ok := input.RemoteExecuteParams["FeatureStore"]
	if !ok {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", ErrorParams
	}
	byteReq, err := json.Marshal(data)
	if err != nil {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", err
	}

	byteResp, err := workflowforeachstep.CallService(workflowforeachstep.FeatureStoreURL, byteReq)
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
	logger.Debug("Feature Store activity completed.")

	return resp.Content, nil
}
