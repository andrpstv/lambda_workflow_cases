package workflows

import (
	"lambda_workflow_cases/scenarios/workflowforeach"
	"lambda_workflow_cases/scenarios/workflowforeach/activities/datastore"
	"lambda_workflow_cases/scenarios/workflowforeach/activities/models"
	"lambda_workflow_cases/scenarios/workflowforeach/activities/strategies"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	MainFlowName       = "UCL scenario1 workflow"
	DataStoreFlowName  = "DataStores workflow"
	ModelsFlowName     = "Models workflow"
	StrategiesFlowName = "Strategies workflow"
)

func Main(ctx workflow.Context, input *workflowforeachstep.Input) (*workflowforeachstep.Output, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info(MainFlowName, "started")

	result := workflowforeachstep.Output{}

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{WorkflowID: "datastore-" + input.GlobalID})

	err := workflow.ExecuteChildWorkflow(ctx, DataStores, input).Get(ctx, &result.DataStore)
	if err != nil {
		logger.Error("Main executiob received Datastore execution failure", "Error", err)
		return &result, err
	}

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{WorkflowID: "models-" + input.GlobalID})

	err = workflow.ExecuteChildWorkflow(ctx, Models, input).Get(ctx, &result.Models)
	if err != nil {
		logger.Error("Main executiob received Models execution failure", "Error", err)
		return &result, err
	}

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{WorkflowID: "strategies-" + input.GlobalID})

	err = workflow.ExecuteChildWorkflow(ctx, Strategies, input).Get(ctx, &result.Strategies)
	if err != nil {
		logger.Error("Main executiob received Strategies execution failure", "Error", err)
		return &result, err
	}
	logger.Info(MainFlowName, "completed")

	return &result, nil
}

func DataStores(ctx workflow.Context, input workflowforeachstep.Input) (*workflowforeachstep.DataStore, error) {
	aoNoRetry := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 60,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, aoNoRetry)

	logger := workflow.GetLogger(ctx)
	logger.Info(DataStoreFlowName + "started")

	result := workflowforeachstep.DataStore{}

	futureFS := workflow.ExecuteActivity(ctx, datastore.FeatureStoreActivity, input)
	futureRA := workflow.ExecuteActivity(ctx, datastore.RiskAvatarActivity, input)
	futureRP := workflow.ExecuteActivity(ctx, datastore.RiskParamsActivity, input)

	err := futureFS.Get(ctx, &result.FeatureStore)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return nil, err
	}

	err = futureRA.Get(ctx, &result.RiskAvatar)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return nil, err
	}

	err = futureRP.Get(ctx, &result.RiskParams)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return nil, err
	}

	logger.Info(DataStoreFlowName + "finished")

	return &result, nil
}
func Models(ctx workflow.Context, input workflowforeachstep.Input) (*workflowforeachstep.Models, error) {
	aoNoRetry := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 60,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, aoNoRetry)

	logger := workflow.GetLogger(ctx)
	logger.Info(ModelsFlowName + "started")

	result := workflowforeachstep.Models{}

	futureM1 := workflow.ExecuteActivity(ctx, models.Model1Activity, input)
	futureM2 := workflow.ExecuteActivity(ctx, models.Model2Activity, input)
	futureM3 := workflow.ExecuteActivity(ctx, models.Model3Activity, input)
	futureM4 := workflow.ExecuteActivity(ctx, models.Model4Activity, input)

	err := futureM1.Get(ctx, &result.Model1)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return nil, err
	}
	err = futureM2.Get(ctx, &result.Model2)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return nil, err
	}
	err = futureM3.Get(ctx, &result.Model3)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return nil, err
	}
	err = futureM4.Get(ctx, &result.Model4)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return nil, err
	}

	logger.Info(ModelsFlowName, "workflow completed")
	return &result, nil
}
func Strategies(ctx workflow.Context, input *workflowforeachstep.Input) (workflowforeachstep.Strategies, error) {
    ao := workflow.ActivityOptions{StartToCloseTimeout: time.Second * 60}
    ctx = workflow.WithActivityOptions(ctx, ao)

    logger := workflow.GetLogger(ctx)
    logger.Info(StrategiesFlowName + "started")

    result := workflowforeachstep.Strategies{}

    var raw string
    futureCvUcl := workflow.ExecuteActivity(ctx, strategies.CvUclActivity, input)
    err := futureCvUcl.Get(ctx, &raw)
    if err != nil {
        logger.Error("Activity failed", "Error", err)
        return result, err
    }

    result.CvUcl = raw

    logger.Info(StrategiesFlowName + "completed")
    return result, nil
}
