package main

import (
	"fmt"
	"time"

	"go.uber.org/cadence/workflow"
)

/**
 * This sample workflow executes multiple branches in parallel. The number of branches is controlled by passed in parameter.
 */

const (
	// ApplicationName is the task list for this sample
	ApplicationName = "branchGroup"

	totalBranches = 3
)

// sampleBranchWorkflow workflow decider
func sampleBranchWorkflow(ctx workflow.Context) error {
	var futures []workflow.Future
	// starts activities in parallel
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for i := 1; i <= totalBranches; i++ {
		activityInput := fmt.Sprintf("branch %d of %d.", i, totalBranches)
		future := workflow.ExecuteActivity(ctx, sampleActivity, activityInput)
		futures = append(futures, future)
	}

	// wait until all futures are done
	for _, future := range futures {
		future.Get(ctx, nil)
	}

	workflow.GetLogger(ctx).Info("Workflow completed.")

	return nil
}

func sampleActivity(input string) (string, error) {
	name := "sampleActivity"
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil
}
