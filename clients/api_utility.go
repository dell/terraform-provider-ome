package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PaginationData common
type PaginationData struct {
	OdataContext string                   `json:"@odata.context"`
	OdataCount   int64                    `json:"@odata.count"`
	Value        []map[string]interface{} `json:"value"`
	NextLink     string                   `json:"@odata.nextLink"`
}

// AuthReq holds payload for authentication to create a session
type AuthReq struct {
	Username    string `json:"UserName"`
	Password    string `json:"Password"`
	SessionType string `json:"SessionType"`
}

// JobStatus is the status returned by the jobs API
type JobStatus struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// JobResp is the response returned by the jobs API
type JobResp struct {
	ID             int64     `json:"Id"`
	JobName        string    `json:"JobName"`
	JobDescription string    `json:"JobDescription"`
	NextRun        string    `json:"NextRun"`
	LastRun        string    `json:"LastRun"`
	StartTime      string    `json:"StartTime"`
	EndTime        string    `json:"EndTime"`
	Schedule       string    `json:"Schedule"`
	State          string    `json:"State"`
	CreatedBy      string    `json:"CreatedBy"`
	LastRunStatus  JobStatus `json:"LastRunStatus"`
	JobType        JobStatus `json:"JobType"`
	JobStatus      JobStatus `json:"JobStatus"`
}

// LastExecutionDetail is response returned by LastExecutionDetail job API
type LastExecutionDetail struct {
	Value string `json:"Value"`
}

// AuthResp is the payload returned by the response of authentication
type AuthResp struct {
	ID string `json:"Id"`
}

// CreateSession is used to create session in OME
func (c *Client) CreateSession() (*http.Response, error) {
	ar := AuthReq{
		Username:    c.username,
		Password:    c.password,
		SessionType: SessionType,
	}
	body, _ := c.JSONMarshal(ar)
	resp, err := c.Post(SessionAPI, nil, body)
	if resp != nil {
		respBody, _ := c.GetBodyData(resp.Body)

		authResp := AuthResp{}
		_ = c.JSONUnMarshal(respBody, &authResp)

		c.SetSessionParams(resp.Header.Get(AuthTokenHeader), authResp.ID)
	}
	return resp, err
}

// RemoveSession is used to remove session in OME
func (c *Client) RemoveSession() (*http.Response, error) {

	api := fmt.Sprintf(SessionAPI+"('%s')", c.sessionID)

	resp, err := c.Delete(api, nil, nil)

	c.SetSessionParams("", "")

	return resp, err
}

// TrackJob - is used to track job status
func (c *Client) TrackJob(jobID int64, maxRetries int64, sleepInterval int64) (bool, string) {
	var status bool
	var message string
	jobRetries := int64(0)
	api := fmt.Sprintf(JobAPI+"(%d)", jobID)
	isJobCompleted := false
	for jobRetries < maxRetries {
		jobRetries++
		time.Sleep(time.Second * time.Duration(sleepInterval))
		resp, err := c.Get(api, nil, nil)
		if err != nil {
			message = err.Error()
			isJobCompleted = true
			break
		}
		if resp != nil {
			jr := &JobResp{}
			parseResponse(c, resp, &jr)
			lrs := jr.LastRunStatus.ID
			if lrs == SuccessStatusID {
				status = true
				message = SuccessMsg
				isJobCompleted = true
				break
			} else if findElementInArray(FailureStatusIDs, lrs) != -1 {
				ledAPI := fmt.Sprintf(LastExecDetailAPI, jobID)
				ledResp, err := c.Get(ledAPI, nil, nil)
				isJobCompleted = true
				if err != nil {
					message = err.Error()
				} else {
					led := LastExecutionDetail{}
					parseResponse(c, ledResp, &led)
					message = led.Value
				}
				break
			}
		}
	}
	if !isJobCompleted {
		message = fmt.Sprintf(JobIncompleteMsg, jobID, maxRetries)
	}

	return status, message
}

// GetJob - returns a job detail for job id
func (c *Client) GetJob(jobID int64) (JobResp, error) {
	api := fmt.Sprintf(JobAPI+"(%d)", jobID)
	resp, err := c.Get(api, nil, nil)
	if err != nil {
		return JobResp{}, err
	}
	jr := &JobResp{}
	parseResponse(c, resp, &jr)
	return *jr, nil
}

func parseResponse(c *Client, resp *http.Response, in interface{}) {
	data, _ := c.GetBodyData(resp.Body)
	_ = c.JSONUnMarshal(data, in)
}

func findElementInArray(arr []any, find any) any {
	index := -1
	for i := 0; i < len(arr); i++ {
		if arr[i] == find {
			index = i
		}
	}
	return index
}

// GetURL returns the url framed from the given host and port
func GetURL(host string, port int64) string {
	return fmt.Sprintf("https://%s:%d", host, port)
}

// ClientPreReqHook - performs the set of predefined operations
func ClientPreReqHook(c *Client, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	if c.GetSessionToken() != "" {
		r.Header.Set(AuthTokenHeader, c.GetSessionToken())
	}
}

// GetPaginatedData - returns all the paginated data
func (c *Client) GetPaginatedData(url string, in interface{}) error {

	response, err := c.Get(url, nil, nil)
	if err != nil {
		return err
	}
	var allData []map[string]interface{}
	pd := PaginationData{}
	bodyData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(bodyData, &pd)
	if err != nil {
		return err
	}
	allData = append(allData, pd.Value...)
	for pd.NextLink != "" {
		response, err := c.Get(pd.NextLink, nil, nil)
		if err != nil {
			return err
		}
		pd = PaginationData{}
		bodyData, _ := c.GetBodyData(response.Body)
		err = c.JSONUnMarshal(bodyData, &pd)
		if err != nil {
			return err
		}
		allData = append(allData, pd.Value...)
	}

	jsonString, _ := json.Marshal(allData)
	_ = json.Unmarshal(jsonString, &in)

	return nil
}

// GetPaginatedDataWithQueryParam - returns all the paginated data with query params
func (c *Client) GetPaginatedDataWithQueryParam(url string, queryParams map[string]string, in interface{}) error {

	response, err := c.Get(url, nil, queryParams)
	if err != nil {
		return err
	}
	var allData []map[string]interface{}
	pd := PaginationData{}
	bodyData, _ := c.GetBodyData(response.Body)
	err = c.JSONUnMarshal(bodyData, &pd)
	if err != nil {
		return err
	}
	allData = append(allData, pd.Value...)
	for pd.NextLink != "" {
		response, err := c.Get(pd.NextLink, nil, nil)
		if err != nil {
			return err
		}
		pd = PaginationData{}
		bodyData, _ := c.GetBodyData(response.Body)
		err = c.JSONUnMarshal(bodyData, &pd)
		if err != nil {
			return err
		}
		allData = append(allData, pd.Value...)
	}

	jsonString, _ := json.Marshal(allData)
	_ = json.Unmarshal(jsonString, &in)

	return nil
}
