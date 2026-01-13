package models

import (
	"context"
	"encoding/json"
	"lambda_workflow_cases/templates/oneworkflow"

	"go.temporal.io/sdk/activity"
)

func Model2Activity(ctx context.Context, ds oneworkflow.DataStore) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Model2Activity started")

	// data, ok := input.RemoteExecuteParams["Model1"]
	// if !ok {
	// 	logger.Error("Activity failed.", "Error", ErrorParams)
	// 	return "", ErrorParams
	// }
	byteReq, err := json.Marshal(ds)
	if err != nil {
		logger.Error("Activity failed.", "Error", ErrorParams)
		return "", err
	}

	byteResp, err := oneworkflow.CallService(oneworkflow.Model2URL, byteReq)
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
	logger.Debug("Feature Store activity completed.")

	return resp.Content, nil
}
