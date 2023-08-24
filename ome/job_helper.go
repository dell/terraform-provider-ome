package ome

import (
	"context"
	"errors"
	"fmt"
	"terraform-provider-ome/clients"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type JobRunner struct {
	client         *clients.Client
	jobID          int64
	maxRetries     int64 // derive from the timeout in minutes
	sleepInterval  int64 // keep sleep interval of 10 seconds.
	partialFailure bool  // to check if partial failure need to be checked
}

func (jr *JobRunner) IsCompleted(ctx context.Context) (clients.JobResp, bool) {
	jobResponse, err := jr.client.GetJob(jr.jobID)
	if err != nil {
		return clients.JobResp{}, false
	}
	return jobResponse, jobResponse.LastRunStatus.ID == clients.SuccessStatusID
}

func (jr *JobRunner) Monitor(ctx context.Context) error {
	for jr.maxRetries > 0 {
		jobResponse, err := jr.client.GetJob(jr.jobID)
		if err != nil {
			return err
		}
		if jobResponse.LastRunStatus.ID != clients.SuccessStatusID {
			time.Sleep(time.Second * time.Duration(jr.sleepInterval)) // job polling delta
			jr.maxRetries--
		} else {
			break
		}
	}
	tflog.Info(ctx, "job has been successfully completed")
	if !jr.partialFailure && jr.maxRetries <= 0 {
		return errors.New("timeout error") // to-do: more appropriate error message
	}
	return nil
}

func (jr *JobRunner) GetLastJobExecution(ctx context.Context) (clients.LastExecutionDetail, error) {
	ledAPI := fmt.Sprintf(clients.LastExecDetailAPI, jr.jobID)
	ledResp, err := jr.client.Get(ledAPI, nil, nil)
	if err != nil {
		return clients.LastExecutionDetail{}, errors.New("get job last execution error: " + err.Error()) // to-do: more appropriate error message
	}
	led := clients.LastExecutionDetail{}
	data, _ := jr.client.GetBodyData(ledResp.Body)

	tflog.Info(ctx, "last execution details data "+string(data))
	err = jr.client.JSONUnMarshal(data, led)
	if err != nil {
		return clients.LastExecutionDetail{}, errors.New("failed to unmarshal last execution detail")
	}
	tflog.Info(ctx, "last execution details has been captured.")
	return led, nil
}

func (jr *JobRunner) GetExecutionDetails(ctx context.Context, executionHistoryID int64) (clients.ExecutionHistories, error) {
	execDetailAPI := fmt.Sprintf("/api/JobService/Jobs(%d)/ExecutionHistories(%d)/ExecutionHistoryDetails", jr.jobID, executionHistoryID)
	execDetails, err := jr.client.Get(execDetailAPI, nil, nil)
	if err != nil {
		return clients.ExecutionHistories{}, errors.New("get job execution history error: " + err.Error()) // to-do: more appropriate error message
	}
	ed := clients.ExecutionHistories{}
	data, _ := jr.client.GetBodyData(execDetails.Body)
	err = jr.client.JSONUnMarshal(data, ed)
	if err != nil {
		return clients.ExecutionHistories{}, errors.New("failed to unmarshal execution histories")
	}
	return ed, nil
}

func DiscoverJobRunner(ctx context.Context, omeClient *clients.Client, jobId, timeout int64) (string, error) {
	var results string
	sleepInterval := int64(10)                 // seconds
	maxRetries := timeout * 60 / sleepInterval // timeout in minutes * 60 seconds / sleep interval in seconds = max number of retries
	jobRunner := JobRunner{
		client:         omeClient,
		jobID:          jobId,
		maxRetries:     maxRetries,
		sleepInterval:  sleepInterval,
		partialFailure: false,
	}
	err := jobRunner.Monitor(ctx)
	if err != nil {
		return "error while monitoring the job", err
	}
	led, err := jobRunner.GetLastJobExecution(ctx)
	if err != nil {
		return "error while get the last job execution details", err
	}
	if ehd, err := jobRunner.GetExecutionDetails(ctx, int64(led.ExecutionHistoryId)); err != nil {
		for _, executionDetail := range ehd.ExecutionDetails {
			results += executionDetail.Value
		}
	}
	return results, nil
}

// func JobSchema(cronPathExpression, timeoutPathExpression path.Expression) map[string]schema.Attribute {
// 	return map[string]schema.Attribute{
// 		"timeout": schema.Int64Attribute{
// 			MarkdownDescription: "Provide job monitor timeout in minutes",
// 			Description:         "Provide job monitor timeout in minutes",
// 			Optional:            true,
// 			Validators: []validator.Int64{
// 				int64validator.AtLeast(1),
// 				int64validator.ExactlyOneOf(cronPathExpression),
// 			},
// 		},
// 		"cron": schema.StringAttribute{
// 			MarkdownDescription: "Provide cron to schedule a job.",
// 			Description:         "Provide cron to schedule a job.",
// 			Optional:            true,
// 			Validators: []validator.String{
// 				stringvalidator.LengthAtLeast(1),
// 				stringvalidator.ExactlyOneOf(timeoutPathExpression),
// 			},
// 		},
// 		"job_id": schema.Int64Attribute{
// 			MarkdownDescription: "Job ID.",
// 			Description:         "Job ID.",
// 			Computed:            true,
// 		},
// 	}
// }
