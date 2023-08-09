package models

// Job for getting job info.
type Job struct {
	ID             int           `json:"Id"`
	JobName        string        `json:"JobName"`
	JobDescription string        `json:"JobDescription"`
	NextRun        any           `json:"NextRun"`
	LastRun        any           `json:"LastRun"`
	StartTime      any           `json:"StartTime"`
	EndTime        any           `json:"EndTime"`
	Schedule       string        `json:"Schedule"`
	State          string        `json:"State"`
	CreatedBy      string        `json:"CreatedBy"`
	UpdatedBy      any           `json:"UpdatedBy"`
	LastRunStatus  LastRunStatus `json:"LastRunStatus"`
	JobType        JobType       `json:"JobType"`
	JobStatus      JobStatus     `json:"JobStatus"`
	Targets        []any         `json:"Targets"`
	Params         []Params      `json:"Params"`
	Visible        bool          `json:"Visible"`
	Editable       bool          `json:"Editable"`
	Builtin        bool          `json:"Builtin"`
	UserGenerated  bool          `json:"UserGenerated"`
	IDUserOwner    int           `json:"IdUserOwner"`
	IDOwner        any           `json:"IdOwner"`
}

// LastRunStatus for job last run status.
type LastRunStatus struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// JobType for job type.
type JobType struct {
	ID                 int    `json:"Id"`
	Name               string `json:"Name"`
	Internal           bool   `json:"Internal"`
	IsShareUsageActive bool   `json:"IsShareUsageActive"`
}

// JobStatus for final job status.
type JobStatus struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// Params for getting job params.
type Params struct {
	JobID int    `json:"JobId"`
	Key   string `json:"Key"`
	Value string `json:"Value"`
}
