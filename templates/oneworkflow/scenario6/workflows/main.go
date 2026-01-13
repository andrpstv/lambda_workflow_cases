package workflows

import (
	"lambda_workflow_cases/templates/oneworkflow"
	"lambda_workflow_cases/templates/oneworkflow/scenario6/activities/datastore"
	"lambda_workflow_cases/templates/oneworkflow/scenario6/activities/models"
	"lambda_workflow_cases/templates/oneworkflow/scenario6/activities/strategy"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const Name = "UCL scenario1 workflow"

var ValidationStatusKey = temporal.NewSearchAttributeKeyKeyword("ValidationStatus")

var (
	m1 oneworkflow.Models
	m2 oneworkflow.Models
	m3 oneworkflow.Models
	m4 oneworkflow.Models
)

func Main(ctx workflow.Context, input *oneworkflow.Input) (*oneworkflow.Output, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    3,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// aoShort := workflow.ActivityOptions{
	// 	StartToCloseTimeout: time.Second * 5,
	// 	RetryPolicy: &temporal.RetryPolicy{
	// 		InitialInterval:    1 * time.Second,
	// 		BackoffCoefficient: 2.0,
	// 		MaximumAttempts:    1,
	// 	},
	// }
	// ctxShort := workflow.WithActivityOptions(ctx, aoShort)
	logger := workflow.GetLogger(ctx)
	logger.Info(Name + " started")

	result := oneworkflow.Output{}

	futureRA := workflow.ExecuteActivity(ctx, datastore.RiskAvatarActivity, input)
	futureRP := workflow.ExecuteActivity(ctx, datastore.RiskParamsActivity, input)

	if err := futureRA.Get(ctx, &result.DataStore.RiskAvatar); err != nil {
		logger.Error("RiskAvatar activity failed", "Error", err)
		return &result, err
	}
	if err := futureRP.Get(ctx, &result.DataStore.RiskParams); err != nil {
		logger.Error("RiskParams activity failed", "Error", err)
		return &result, err
	}
	f1, s1 := workflow.NewFuture(ctx)
	f2, s2 := workflow.NewFuture(ctx)
	f3, s3 := workflow.NewFuture(ctx)
	f4, s4 := workflow.NewFuture(ctx)

	workflow.Go(ctx, func(ctx workflow.Context) {
		var fs oneworkflow.DataStore
		if err := workflow.ExecuteActivity(ctx, datastore.FeatureStoreActivity, input).
			Get(ctx, &fs.FeatureStore); err != nil {
			logger.Error("FeatureStoreActivity for Model1 failed", "Error", err)
			s1.Set(nil, err)
			return
		}

		var r oneworkflow.Models
		if err := workflow.ExecuteActivity(ctx, models.Model1Activity, fs).
			Get(ctx, &r.Model1); err != nil {
			logger.Error("Model1Activity failed", "Error", err)
			s1.Set(nil, err)
			return
		}

		s1.Set(r.Model1, nil)
	})

	workflow.Go(ctx, func(ctx workflow.Context) {
		var fs oneworkflow.DataStore
		if err := workflow.ExecuteActivity(ctx, datastore.FeatureStoreActivity, input).
			Get(ctx, &fs.FeatureStore); err != nil {
			logger.Error("FeatureStoreActivity for Model2 failed", "Error", err)
			s2.Set(nil, err)
			return
		}

		var r oneworkflow.Models
		if err := workflow.ExecuteActivity(ctx, models.Model2Activity, fs).
			Get(ctx, &r.Model2); err != nil {
			logger.Error("Model2Activity failed", "Error", err)
			s2.Set(nil, err)
			return
		}

		s2.Set(r.Model2, nil)
	})

	workflow.Go(ctx, func(ctx workflow.Context) {
		var fs oneworkflow.DataStore
		if err := workflow.ExecuteActivity(ctx, datastore.FeatureStoreActivity, input).
			Get(ctx, &fs.FeatureStore); err != nil {
			logger.Error("FeatureStoreActivity for Model3 failed", "Error", err)
			s3.Set(nil, err)
			return
		}

		var r oneworkflow.Models
		if err := workflow.ExecuteActivity(ctx, models.Model3Activity, fs).
			Get(ctx, &r.Model3); err != nil {
			logger.Error("Model3Activity failed", "Error", err)
			s3.Set(nil, err)
			return
		}

		s3.Set(r.Model3, nil)
	})

	workflow.Go(ctx, func(ctx workflow.Context) {
		var fs oneworkflow.DataStore
		if err := workflow.ExecuteActivity(ctx, datastore.FeatureStoreActivity, input).
			Get(ctx, &fs.FeatureStore); err != nil {
			logger.Error("FeatureStoreActivity for Model4 failed", "Error", err)
			s4.Set(nil, err)
			return
		}

		var r oneworkflow.Models
		if err := workflow.ExecuteActivity(ctx, models.Model4Activity, fs).
			Get(ctx, &r.Model4); err != nil {
			logger.Error("Model4Activity failed", "Error", err)
			s4.Set(nil, err)
			return
		}

		s4.Set(r.Model4, nil)
	})

	if err := f1.Get(ctx, &m1.Model1); err != nil {
		return &result, err
	}
	if err := f2.Get(ctx, &m2.Model2); err != nil {
		return &result, err
	}
	if err := f3.Get(ctx, &m3.Model3); err != nil {
		return &result, err
	}
	if err := f4.Get(ctx, &m4.Model4); err != nil {
		return &result, err
	}
	result.Models.Model1 = m1.Model1
	result.Models.Model2 = m2.Model2
	result.Models.Model3 = m3.Model3
	result.Models.Model4 = m4.Model4

	futureStrategy := workflow.ExecuteActivity(ctx, strategy.StrategyActivity, input)
	if err := futureStrategy.Get(ctx, &result.Strategy.Strategy); err != nil {
		logger.Error("StrategyActivity failed", "Error", err)
		return &result, err
	}

	logger.Info(Name + " finished")
	return &result, nil
}
