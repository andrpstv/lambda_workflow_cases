package workflows

import (
	"lambda_workflow_cases/scenarios/oneworkflow"
	"lambda_workflow_cases/scenarios/oneworkflow/activities/datastore"
	"lambda_workflow_cases/scenarios/oneworkflow/activities/models"
	"lambda_workflow_cases/scenarios/oneworkflow/activities/strategies"
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
		StartToCloseTimeout: time.Second * 60,
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
		logger.Error("Activity failed", "Error", err)
		return &result, err
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

	futureCvUcl := workflow.ExecuteActivity(ctx, strategies.CvUclActivity, input)

	err = futureCvUcl.Get(ctx, &result.Strategies.CvUcl)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return &result, err
	}
	logger.Info(Name + "finished")

	return &result, nil
}
