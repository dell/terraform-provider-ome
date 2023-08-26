package ome

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"terraform-provider-ome/clients"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	// CompletedWithSuccess to get job complete with success status
	CompletedWithSuccess = 2060
	// Failed to get job failed status
	Failed = 2070
	// CompletedWithError to get job complete with error status
	CompletedWithError = 2090
	// Aborted to get job aborted status
	Aborted = 2100
	// Stopped to get job stop status
	Stopped = 2102
	// Cancelled to get job cancel status
	Cancelled = 2103
)

// JobRunner to construct the job data and method for job runs.
type JobRunner struct {
	client *clients.Client
	// job id of the task
	jobID int64
	// derive from the timeout in minutes
	maxRetries int64
	// keep sleep interval of 10 seconds.
	sleepInterval int64
	// to check if partial failure need to be checked
	partialFailure bool
}

// ParseResponse to unmarshal the job json models.
func (jr *JobRunner) ParseResponse(ctx context.Context, resp *http.Response, in interface{}) error {
	data, err := jr.client.GetBodyData(resp.Body)
	if err != nil {
		return err
	}
	tflog.Info(ctx, "job runner data "+string(data))
	err = jr.client.JSONUnMarshal(data, in)
	if err != nil {
		return err
	}
	tflog.Info(ctx, "job runner interface has been succesfully unmarshal.")
	return nil
}

// Monitor to monitor the job till timeout.
func (jr *JobRunner) Monitor(ctx context.Context) error {
	for jr.maxRetries > 0 {
		jobResponse, err := jr.client.GetJob(jr.jobID)
		if err != nil {
			return err
		}
		tflog.Info(ctx, "LASTRUNSTATUS ID: "+fmt.Sprint(jobResponse.LastRunStatus.ID))
		// stop monitoring in case of job termination states - failed, aborted, stopped, cancelled, completed with success and completed with error.
		if jobResponse.LastRunStatus.ID == Failed ||
			jobResponse.LastRunStatus.ID == Aborted ||
			jobResponse.LastRunStatus.ID == Stopped ||
			jobResponse.LastRunStatus.ID == Cancelled ||
			jobResponse.LastRunStatus.ID == CompletedWithSuccess ||
			jobResponse.LastRunStatus.ID == CompletedWithError {
			break
		} else if jobResponse.LastRunStatus.ID == CompletedWithError && !jr.partialFailure {
			// if completed with error and ignore partial failure is false, then error out.
			return errors.New("job completed with errors")
		} else {
			// job polling delta
			time.Sleep(time.Second * time.Duration(jr.sleepInterval))
			jr.maxRetries--
		}
	}
	tflog.Info(ctx, "job has been successfully completed")
	if !jr.partialFailure && jr.maxRetries <= 0 {
		return errors.New("")
	}
	return nil
}

// GetLastJobExecution to get the last job execution.
func (jr *JobRunner) GetLastJobExecution(ctx context.Context) (clients.LastExecutionDetail, error) {
	ledAPI := fmt.Sprintf(clients.LastExecDetailAPI, jr.jobID)
	ledResp, err := jr.client.Get(ledAPI, nil, nil)
	if err != nil {
		return clients.LastExecutionDetail{}, errors.New("get job last execution error: " + err.Error())
	}
	led := clients.LastExecutionDetail{}
	err = jr.ParseResponse(ctx, ledResp, &led)
	if err != nil {
		return clients.LastExecutionDetail{}, err
	}
	return led, nil
}

// GetExecutionDetails to get the execution detail of job runs.
func (jr *JobRunner) GetExecutionDetails(ctx context.Context, executionHistoryID int64) (clients.ExecutionHistories, error) {
	execDetailAPI := fmt.Sprintf("/api/JobService/Jobs(%d)/ExecutionHistories(%d)/ExecutionHistoryDetails", jr.jobID, executionHistoryID)
	execDetails, err := jr.client.Get(execDetailAPI, nil, nil)
	if err != nil {
		return clients.ExecutionHistories{}, errors.New("get job execution details error: " + err.Error())
	}
	ed := clients.ExecutionHistories{}
	err = jr.ParseResponse(ctx, execDetails, &ed)
	if err != nil {
		return clients.ExecutionHistories{}, err
	}
	return ed, nil
}

// DiscoverJobRunner to track the discover job.
func DiscoverJobRunner(ctx context.Context, omeClient *clients.Client, jobId, timeout int64, partialFailure bool) ([]string, error) {
	results := make([]string, 0)
	sleepInterval := int64(10)
	maxRetries := timeout * 60 / sleepInterval
	jobRunner := JobRunner{
		client:         omeClient,
		jobID:          jobId,
		maxRetries:     maxRetries,
		sleepInterval:  sleepInterval,
		partialFailure: partialFailure,
	}
	/*
		The job runner needs to wait for an ideal sleep interval before monitoring so that the latest execution status is refreshed on the job.
		If an update operation is performed, the job runner monitor will exit immediately. In such cases, it will fetch the last execution status, which may have already been completed.
		However, this will not point to the case where the job has been updated. Therefore, a sleep interval is necessary to ensure that we fetch the latest execution status and not any historical execution completed status.
	*/
	time.Sleep(time.Second * time.Duration(sleepInterval))

	err := jobRunner.Monitor(ctx)
	if err != nil {
		return results, err
	}
	led, err := jobRunner.GetLastJobExecution(ctx)
	tflog.Info(ctx, "last execution details model "+fmt.Sprint(led))
	if err != nil {
		return results, err
	}
	ehd, err := jobRunner.GetExecutionDetails(ctx, int64(led.ExecutionHistoryId))
	if err != nil {
		return results, err
	}
	tflog.Info(ctx, "execution details model "+fmt.Sprint(ehd))
	for _, ehd := range ehd.ExecutionDetails {
		results = append(results, ehd.Value)
	}
	tflog.Info(ctx, results[0])
	return results, nil
}
