package main

import (
	"context"
	workflowforeachstep "lambda_workflow_cases/scenarios/workflowforeach"
	"lambda_workflow_cases/scenarios/workflowforeach/workflows"
	"log"
	"math/rand"
	"strconv"
	"time"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{Namespace: "scenario2"})
	if err != nil {
		log.Fatalf("Unable to connect to temporal client: %v", err)
	}
	defer c.Close()

	input := workflowforeachstep.Input{
		GlobalID:            strconv.Itoa(GenerateGlobalID()),
		EpkID:               "1234567890",
		RemoteExecuteParams: map[string]workflowforeachstep.RemoteRequest{},
	}

	input.RemoteExecuteParams["FeatureStore"] = workflowforeachstep.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["RiskAvatar"] = workflowforeachstep.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["RiskParams"] = workflowforeachstep.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Model1"] = workflowforeachstep.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Model2"] = workflowforeachstep.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Model3"] = workflowforeachstep.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Model4"] = workflowforeachstep.RemoteRequest{DelaySec: 1, RespSizeKb: 5}
	input.RemoteExecuteParams["Strategy"] = workflowforeachstep.RemoteRequest{DelaySec: 1, RespSizeKb: 5}

	options := client.StartWorkflowOptions{
		ID:        "ucl_workflow_" + input.GlobalID,
		TaskQueue: workflowforeachstep.TaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.Main, &input)
	if err != nil {
		log.Fatalf("Failed to execute workflow: %v", err)
	}
	log.Println("Started Workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result workflowforeachstep.Output
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalf("Failed to get result: %v", err)
	}
}
func GenerateGlobalID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	minimal := 1_000_000_000
	maximum := 9_999_999_999

	return r.Intn(maximum-minimal+1) + minimal
}
