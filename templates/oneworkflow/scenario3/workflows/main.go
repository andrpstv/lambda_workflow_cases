package workflows

import (
	"errors"
	"lambda_workflow_cases/templates/oneworkflow"
	"lambda_workflow_cases/templates/oneworkflow/scenario2/activities/datastore"
	"lambda_workflow_cases/templates/oneworkflow/scenario2/activities/models"
	"lambda_workflow_cases/templates/oneworkflow/scenario2/activities/strategy"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const Name = "UCL scenario1 workflow"

type StepStatus string

const (
	StepPending StepStatus = "PENDING"
	StepOK      StepStatus = "OK"
	StepFailed  StepStatus = "FAILED"
)

type WorkflowState struct {
	Aborted bool

	FeatureStore StepStatus
	RiskAvatar   StepStatus
	RiskParams   StepStatus

	Model1 StepStatus
	Model2 StepStatus
	Model3 StepStatus
	Model4 StepStatus

	Strategy StepStatus
}

func Main(ctx workflow.Context, input *oneworkflow.Input) (*oneworkflow.Output, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 60,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	aoNoRetry := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 15,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctxNoRetry := workflow.WithActivityOptions(ctx, aoNoRetry)

	logger := workflow.GetLogger(ctx)
	logger.Info(Name + " started")

	result := oneworkflow.Output{}

	state := &WorkflowState{
		FeatureStore: StepPending,
		RiskAvatar:   StepPending,
		RiskParams:   StepPending,
		Model1:       StepPending,
		Model2:       StepPending,
		Model3:       StepPending,
		Model4:       StepPending,
		Strategy:     StepPending,
	}
	err := workflow.SetQueryHandler(ctx, "getState", func() (*WorkflowState, error) {
		return state, nil
	})
	if err != nil {
		logger.Error("SetQueryHandler failed", "Error", err)
		return &result, err
	}

	futureFS := workflow.ExecuteActivity(ctxNoRetry, datastore.FeatureStoreActivity, input)
	futureRA := workflow.ExecuteActivity(ctx, datastore.RiskAvatarActivity, input)
	futureRP := workflow.ExecuteActivity(ctx, datastore.RiskParamsActivity, input)

	err = futureFS.Get(ctx, &result.DataStore.FeatureStore)
	if err != nil {
		var timeoutErr *temporal.TimeoutError
		if errors.As(err, &timeoutErr) {
			logger.Error("FeatureStore activity timed out after retries",
				"TimeoutType", timeoutErr.TimeoutType(),
				"Error", err)
			state.FeatureStore = StepFailed
		} else {
			logger.Error("FeatureStore activity failed", "Error", err)
			state.FeatureStore = StepFailed
			return &result, err
		}
	} else {
		state.FeatureStore = StepOK
	}

	err = futureRA.Get(ctx, &result.DataStore.RiskAvatar)
	if err != nil {
		logger.Error("RiskAvatar activity failed", "Error", err)
		state.RiskAvatar = StepFailed
		return &result, err
	}
	state.RiskAvatar = StepOK

	err = futureRP.Get(ctx, &result.DataStore.RiskParams)
	if err != nil {
		logger.Error("RiskParams activity failed", "Error", err)
		state.RiskParams = StepFailed
		return &result, err
	}
	state.RiskParams = StepOK

	if state.Aborted {
		logger.Info("Workflow aborted after DataStore block")
		return &result, nil
	}

	err = result.Validate()
	if err != nil {
		logger.Error("DataStores return invalid data", "Error", err)
		return nil, err
	}

	futureM1 := workflow.ExecuteActivity(ctx, models.Model1Activity, input)
	futureM2 := workflow.ExecuteActivity(ctx, models.Model2Activity, input)
	futureM3 := workflow.ExecuteActivity(ctx, models.Model3Activity, input)
	futureM4 := workflow.ExecuteActivity(ctx, models.Model4Activity, input)

	err = futureM1.Get(ctx, &result.Models.Model1)
	if err != nil {
		logger.Error("Model1 activity failed", "Error", err)
		state.Model1 = StepFailed
		return &result, err
	}
	state.Model1 = StepOK

	err = futureM2.Get(ctx, &result.Models.Model2)
	if err != nil {
		logger.Error("Model2 activity failed", "Error", err)
		state.Model2 = StepFailed
		return &result, err
	}
	state.Model2 = StepOK

	err = futureM3.Get(ctx, &result.Models.Model3)
	if err != nil {
		logger.Error("Model3 activity failed", "Error", err)
		state.Model3 = StepFailed
		return &result, err
	}
	state.Model3 = StepOK

	err = futureM4.Get(ctx, &result.Models.Model4)
	if err != nil {
		logger.Error("Model4 activity failed", "Error", err)
		state.Model4 = StepFailed
		return &result, err
	}
	state.Model4 = StepOK

	if state.Aborted {
		logger.Info("Workflow aborted after Models block")
		return &result, nil
	}
	decision := &oneworkflow.DecisionSignal{}
	chDecision := workflow.GetSignalChannel(ctx, "decision")
	logger.Info("Waiting for decision signal (APPROVED / REJECTED)")
	chDecision.Receive(ctx, decision)

	switch decision.Decision {
	case oneworkflow.DecisionApproved:
		logger.Info("Decision APPROVED, continue to Strategy block")

	case oneworkflow.DecisionRejected:
		logger.Info("Decision REJECTED, finishing workflow",
			"Reason", decision.Reason)
		return &result, nil

	default:
		logger.Error("Unknown decision", "Decision", decision.Decision)
		return &result, temporal.NewNonRetryableApplicationError(
			"unknown decision",
			"DecisionError",
			nil,
		) //не ретраимся
	}
	futureStrategy := workflow.ExecuteActivity(ctx, strategy.StrategyActivity, input)
	err = futureStrategy.Get(ctx, &result.Strategy.Strategy)
	if err != nil {
		logger.Error("Strategy activity failed", "Error", err)
		state.Strategy = StepFailed
		return &result, err
	}
	state.Strategy = StepOK

	logger.Info(Name + " finished")
	return &result, nil
}
