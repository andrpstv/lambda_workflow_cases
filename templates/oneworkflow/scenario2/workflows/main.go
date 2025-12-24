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
	logger.Info(Name + "started")

	result := oneworkflow.Output{}

	futureFS := workflow.ExecuteActivity(ctxNoRetry, datastore.FeatureStoreActivity, input)
	futureRA := workflow.ExecuteActivity(ctx, datastore.RiskAvatarActivity, input)
	futureRP := workflow.ExecuteActivity(ctx, datastore.RiskParamsActivity, input)
	err := futureFS.Get(ctx, &result.DataStore.FeatureStore)
	if err != nil {
		var timeoutErr *temporal.TimeoutError
		if errors.As(err, &timeoutErr) {
			logger.Error("FeatureStore activity timed out after retries",
				"TimeoutType", timeoutErr.TimeoutType(),
				"Error", err)
		} else {
			logger.Error("Activity failed", "Error", err)
			return &result, err
		}
	}

	err = futureRA.Get(ctx, &result.DataStore.RiskAvatar)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return &result, err
	}

	err = futureRP.Get(ctx, &result.DataStore.RiskParams)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return &result, err
	}

	err = result.Validate()
	if err != nil {
		logger.Error("DataStores return invalid data: ", err)
		return nil, err
	}

	futureM1 := workflow.ExecuteActivity(ctx, models.Model1Activity, input)
	futureM2 := workflow.ExecuteActivity(ctx, models.Model2Activity, input)
	futureM3 := workflow.ExecuteActivity(ctx, models.Model3Activity, input)
	futureM4 := workflow.ExecuteActivity(ctx, models.Model4Activity, input)

	err = futureM1.Get(ctx, &result.Models.Model1)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return &result, err
	}
	err = futureM2.Get(ctx, &result.Models.Model2)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return &result, err
	}
	err = futureM3.Get(ctx, &result.Models.Model3)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return &result, err
	}
	err = futureM4.Get(ctx, &result.Models.Model4)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return &result, err
	}

	futureStrategy := workflow.ExecuteActivity(ctx, strategy.StrategyActivity, input)

	err = futureStrategy.Get(ctx, &result.Strategy.Strategy)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return &result, err
	}
	logger.Info(Name + "finished")

	return &result, nil
}
