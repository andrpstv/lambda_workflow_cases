package main

import (
	"lambda_workflow_cases/scenarios/oneworkflow/activities/models"
	"lambda_workflow_cases/scenarios/oneworkflow/activities/strategies"
	workflowforeachstep "lambda_workflow_cases/scenarios/workflowforeach"
	"lambda_workflow_cases/scenarios/workflowforeach/activities/datastore"
	"lambda_workflow_cases/scenarios/workflowforeach/workflows"
	"log"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	StoreImage    = "💾"
	ModelImage    = "🛠️"
	StrategyImage = "⚙️"
)

func main() {
	c, err := client.Dial(client.Options{Namespace: "scenario2"})
	if err != nil {
		log.Fatalln("Unable to connect to temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, workflowforeachstep.TaskQueueName, worker.Options{})

	w.RegisterWorkflow(workflows.Main)
	w.RegisterWorkflow(workflows.DataStores)
	w.RegisterWorkflow(workflows.Models)
	w.RegisterWorkflow(workflows.Strategies)

	w.RegisterActivityWithOptions(datastore.FeatureStoreActivity, activity.RegisterOptions{Name: "FeatureStore " + StoreImage})
	w.RegisterActivityWithOptions(datastore.RiskAvatarActivity, activity.RegisterOptions{Name: "RiskAvatar " + StoreImage})
	w.RegisterActivityWithOptions(datastore.RiskParamsActivity, activity.RegisterOptions{Name: "RiskParams " + StoreImage})

	w.RegisterActivityWithOptions(models.Model1Activity, activity.RegisterOptions{Name: "Model1 " + ModelImage})
	w.RegisterActivityWithOptions(models.Model2Activity, activity.RegisterOptions{Name: "Model2 " + ModelImage})
	w.RegisterActivityWithOptions(models.Model3Activity, activity.RegisterOptions{Name: "Model3 " + ModelImage})
	w.RegisterActivityWithOptions(models.Model4Activity, activity.RegisterOptions{Name: "Model4 " + ModelImage})

	w.RegisterActivityWithOptions(strategies.CvUclActivity, activity.RegisterOptions{Name: "CV_UCL " + StrategyImage})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
