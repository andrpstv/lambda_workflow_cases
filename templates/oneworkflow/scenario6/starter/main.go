package main

import (
	"context"
	"lambda_workflow_cases/templates/oneworkflow"
	"lambda_workflow_cases/templates/oneworkflow/scenario2/workflows"
	"log"
	"math/rand"
	"strconv"
	"time"

	"go.temporal.io/api/common/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{Namespace: "scenario2"})
	if err != nil {
		log.Fatalf("Unable to connect to temporal client: %v", err)
	}
	defer c.Close()

	input := oneworkflow.Input{
		GlobalID:            strconv.Itoa(GenerateGlobalID()),
		EpkID:               "1234567890",
		RemoteExecuteParams: map[string]oneworkflow.RemoteRequest{},
	}

	// ВАЖНО: перед определением задержки проверить на соответствие таймауту на целевом сервере
	input.RemoteExecuteParams["FeatureStore"] = oneworkflow.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["RiskAvatar"] = oneworkflow.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["RiskParams"] = oneworkflow.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Model1"] = oneworkflow.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Model2"] = oneworkflow.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Model3"] = oneworkflow.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Model4"] = oneworkflow.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Strategy"] = oneworkflow.RemoteRequest{DelaySec: 1, RespSizeKb: 5, Fail: true}

	options := client.StartWorkflowOptions{
		ID:        "strategy_workflow_" + input.GlobalID,
		TaskQueue: oneworkflow.TaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.Main, &input)
	if err != nil {
		log.Fatalf("Failed to execute workflow: %v", err)
	}
	log.Println("Started Workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result oneworkflow.Output
	err = we.Get(context.Background(), &result)
	if err != nil {
		const workflowTaskEventID int64 = 34

		req := &workflowservice.ResetWorkflowExecutionRequest{
			Namespace: "scenario2",
			WorkflowExecution: &common.WorkflowExecution{
				WorkflowId: we.GetID(),
				RunId:      we.GetRunID(),
			},
			WorkflowTaskFinishEventId: workflowTaskEventID,
			Reason:                    "reset from Go client after failure",
		}

		_, resetErr := c.ResetWorkflowExecution(context.Background(), req)
		if resetErr != nil {
			log.Fatalf("Failed to reset workflow: %v", resetErr)
		}

		log.Fatalf("Failed to get result: %v", err)
	}

	log.Printf("Workflow result: %+v", result)
}
func GenerateGlobalID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	minimal := 1_000_000_000
	maximum := 9_999_999_999

	return r.Intn(maximum-minimal+1) + minimal
}
