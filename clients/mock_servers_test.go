package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"terraform-provider-ome/models"
	"testing"
	"time"
)

func createNewTLSServer(t *testing.T) *httptest.Server {
	// create a listener with the desired port.
	attemtp := 0
	jobRetries := 0
	l, err := net.Listen("tcp", "127.0.0.1:8234")
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		shouldReturn1 := mockJobsAPI(r, w, &jobRetries) || mockSessionAPIs(r, w) || mockTimeoutAPIs(r, &attemtp, w) || mockDeviceAPIs(r, w) || mockUpdateNetworkConfigAPI(r, w) || mockUnassignProfileAPI(r, w) || mockDeleteProfileAPI(r, w)
		if shouldReturn1 {
			return
		}

		shouldReturn2 := mockTemplateAPIs(r, w) || mockViewTypeAPIs(r, w) || mockDeviceTypeAPIs(r, w) || mockUpdateTemplateAPIs(r, w) || mockGetTemplateByNameAPIs(r, w) || mockGetNetworkAttributesAPI(r, w)
		if shouldReturn2 {
			return
		}

		shouldReturn3 := mockGetTemplateByIDAPIs(r, w) || mockTrackTemplateCreationAPIs(r, w, &jobRetries) || mockGetTemplateAttributesAPIs(r, w) || mockGetIdentityPoolAPIs(r, w) || mockCloneTemplateAPI(r, w)
		if shouldReturn3 {
			return
		}

		shouldReturn4 := mockGroupServiceAPIs(r, w) || mockDeployAPIs(r, w) || mockGetServerProfileInfoByTemplateNameAPIs(r, w) || mockNetworkVlanAPI(r, w)
		if shouldReturn4 {
			return
		}

		shouldReturn5 := mockGetServerProfileInfoByTemplateNameAPIs(r, w)
		if shouldReturn5 {
			return
		}

		shouldReturn6 := mockBaselineAPIs(r, w) || mockGetBaselineByIDAPI(r, w) ||
			mockGetBaselineDevComplianceReportByIDAPI(r, w) || mockGetBaselineDevAttrComplianceReportByIDAPI(r, w) || mockGetBaselineByNameAPI(r, w)
		if shouldReturn6 {
			return
		}

		mockGeneralAPIs(r, w)
	}))

	// NewUnstartedServer creates a listener. Close that listener and replace
	// with the one we created.
	ts.Listener.Close()
	ts.Listener = l

	// Start the server.
	ts.StartTLS()
	return ts
}

func createNewTLSServerWithPort(t *testing.T, port int64, handler http.HandlerFunc) *httptest.Server {
	// create a listener with the desired port.

	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewUnstartedServer(http.HandlerFunc(handler))

	// NewUnstartedServer creates a listener. Close that listener and replace
	// with the one we created.
	ts.Listener.Close()
	ts.Listener = l

	// Start the server.
	ts.StartTLS()
	return ts
}

func mockPortUnAuth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CGEN1001",
						"RelatedProperties": [],
						"Message": "Unable to complete the operation because the user is not authenticated.",
						"MessageArgs": [],
						"Severity": "Critical",
						"Resolution": "Log in with valid credentials and retry the operation."
					}
				]
			}
		}`))
}

func mockPortInValidJSON(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{[,]}`))
}

func mockGeneralAPIs(r *http.Request, w http.ResponseWriter) bool {
	if r.URL.Path != "/emptyBody" {
		if r.Method == "GET" || r.Method == "PATCH" || r.Method == "PUT" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`Hello from TLS server`))
		}
		if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`Hello from TLS server post body`))
		}
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
		}
		return true
	}
	return false
}

func mockTimeoutAPIs(r *http.Request, attemtp *int, w http.ResponseWriter) bool {
	if r.URL.Path == "/timeout" {
		time.Sleep(2 * time.Second)
		return true
	}
	if r.URL.Path == "/timeout-success" {
		*attemtp++
		if *attemtp == 3 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`Hello from TLS server`))
		} else {
			time.Sleep(4 * time.Second)
		}
		return true
	}
	return false
}

func mockSessionAPIs(r *http.Request, w http.ResponseWriter) bool {
	if r.URL.Path == "/api/SessionService/Sessions" && r.Method == "POST" {

		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "admin") {
			w.Header().Set("x-auth-token", "13bc3f63-9376-44dc-a09f-3a94591a7c5d")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
					"Id": "e1817fe6-97e5-4ea0-88a9-d865c73021529",
					"Description": "Create Session",
					"Name": "administrator",
					"UserName": "admin",
					"Password": null,
					"Roles": [
					"ADMINISTRATOR"
					],
					"IpAddress": "xx.xx.xx.xx",
					"StartTimeStamp": "2017-04-17 20:07:12.357",
					"LastAccessedTimeStamp": "2017-04-17 20:07:12.357"
				   }`))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
		return true
	}
	if r.URL.Path == "/api/SessionService/Sessions('e1817fe6-97e5-4ea0-88a9-d865c73021529')" && r.Method == "DELETE" {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func mockJobsAPI(r *http.Request, w http.ResponseWriter, jobRetries *int) bool {
	if r.URL.Path == "/api/JobService/Jobs(1)" && r.Method == "GET" {
		w.Write([]byte(buildJobResponse(2060, "success")))
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(2)" && r.Method == "GET" {
		w.Write([]byte(buildJobResponse(2070, "success")))
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(3)" && r.Method == "GET" {
		w.Write([]byte(buildJobResponse(2090, "success")))
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(4)" && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(buildJobResponse(2060, "success")))
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(12345)" && r.Method == "GET" {
		w.Write([]byte(buildJobResponse(2060, "success")))
		return true
	}
	if (r.URL.Path == "/api/JobService/Jobs(23456)" || r.URL.Path == "/api/JobService/Jobs(14567)") && r.Method == "GET" {
		w.Write([]byte(buildJobResponse(2070, "Failure")))
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(34567)" && r.Method == "GET" {
		w.Write([]byte(buildJobResponse(2090, "Warning")))
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(45678)" && r.Method == "GET" {
		*jobRetries++
		if *jobRetries == 2 {
			*jobRetries = 0
			w.Write([]byte(buildJobResponse(2060, "Success")))
		} else {
			w.Write([]byte(buildJobResponse(2050, "Running")))
		}
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(56789)" && r.Method == "GET" {
		w.Write([]byte(buildJobResponse(2050, "Running")))
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(13456)" && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CJOB4041",
							"RelatedProperties": [],
							"Message": "Unable to retrieve or modify the job information because no record exists for the job ID entered.",
							"MessageArgs": [],
							"Severity": "Warning",
							"Resolution": "Enter a valid job ID and retry the operation."
						}
					]
				}
			}`))
		return true
	}
	if (r.URL.Path == "/api/JobService/Jobs(23456)/LastExecutionDetail") && r.Method == "GET" {
		w.Write([]byte(`{"Value": "LastExecutionDetail Failure"}`))
		return true
	}
	if (r.URL.Path == "/api/JobService/Jobs(34567)/LastExecutionDetail") && r.Method == "GET" {
		w.Write([]byte(`{"Value": "LastExecutionDetail Warning"}`))
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(10860)" && r.Method == "GET" {
		w.Write([]byte(buildJobResponse(2060, "success")))
		return true
	}
	if r.URL.Path == "/api/JobService/Jobs(10861)" && r.Method == "GET" {
		w.Write([]byte(buildJobResponse(2070, "Failure")))
		return true
	}
	if (r.URL.Path == "/api/JobService/Jobs(14567)/LastExecutionDetail" || r.URL.Path == "/api/JobService/Jobs(10861)/LastExecutionDetail") && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CJOB5004",
						"RelatedProperties": [],
						"Message": "No recent execution details were found for the provided job id.",
						"MessageArgs": [],
						"Severity": "Warning",
						"Resolution": "Make sure that the provided Job has been executed before and is a valid job."
					}
				]
			}
		}`))
		return true
	}
	return false
}

func buildJobResponse(id int, msg string) string {
	s := `{
	"LastRunStatus": {
		"@odata.type": "#JobService.JobStatus",
		"Id":` + strconv.Itoa(id) + `,
		"Name": "` + msg + `"` + `
		}
	}`
	return s
}

func mockDeviceAPIs(r *http.Request, w http.ResponseWriter) bool {
	if (strings.Contains(r.URL.RawQuery, "SVT123") || strings.Contains(r.URL.RawQuery, "SVT223")) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
				{
					"Id": 123456
				}
			]
		}`))
		return true
	}

	if strings.Contains(r.URL.RawQuery, "SV6789") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": []
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "INVJSON") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
				{
					"Id": 
				},
			]
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "NOAUTH") && r.Method == "GET" {
		w.WriteHeader(http.StatusUnauthorized)
		return true
	}

	if (strings.Contains(r.URL.RawQuery, "123456") || strings.Contains(r.URL.RawQuery, "223456")) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
				{
					"Id": 123456
				}
			]
		}`))
		return true
	}

	if strings.Contains(r.URL.RawQuery, "123457") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": []
		}`))
		return true
	}
	return false
}

func mockViewTypeAPIs(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == TemplateViewTypeAPI) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(TemplateService.TemplateViewTypeModel)",
			"@odata.count": 5,
			"value": [
				{
					"@odata.type": "#TemplateService.TemplateViewTypeModel",
					"Id": 0,
					"Description": "None"
				},
				{
					"@odata.type": "#TemplateService.TemplateViewTypeModel",
					"Id": 1,
					"Description": "Compliance"
				},
				{
					"@odata.type": "#TemplateService.TemplateViewTypeModel",
					"Id": 2,
					"Description": "Deployment"
				},
				{
					"@odata.type": "#TemplateService.TemplateViewTypeModel",
					"Id": 3,
					"Description": "Inventory"
				},
				{
					"@odata.type": "#TemplateService.TemplateViewTypeModel",
					"Id": 4,
					"Description": "Sample"
				}
			]
		}`))
		return true
	}
	return false
}

func mockDeviceTypeAPIs(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == TemplateDeviceTypeAPI) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(TemplateService.TemplateTypeModel)",
			"@odata.count": 4,
			"value": [
				{
					"@odata.type": "#TemplateService.TemplateTypeModel",
					"Id": 0,
					"Name": "None"
				},
				{
					"@odata.type": "#TemplateService.TemplateTypeModel",
					"Id": 2,
					"Name": "Server"
				},
				{
					"@odata.type": "#TemplateService.TemplateTypeModel",
					"Id": 4,
					"Name": "Chassis"
				},
				{
					"@odata.type": "#TemplateService.TemplateTypeModel",
					"Id": 3,
					"Name": "IO Module"
				}
			]
		}`))
		return true
	}
	return false
}

func mockTemplateAPIs(r *http.Request, w http.ResponseWriter) bool {

	if (r.URL.Path == TemplateAPI) && r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "TestTemplate") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`1`))
		} else if strings.Contains(string(body), "ExtTemplate") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CTEM1027",
							"RelatedProperties": [],
							"Message": "Unable to create the template because the template name ExtTemplate already exists.",
							"MessageArgs": [
								"Temp12"
							],
							"Severity": "Warning",
							"Resolution": "Enter a unique template name and retry the operation."
						}
					]
				}
			}`))
		} else if strings.Contains(string(body), "TemplateInvalidDevice") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CTEM1025",
							"RelatedProperties": [],
							"Message": "Unable to create or deploy the template because the device ID 23456 is invalid.",
							"MessageArgs": [
								"1009"
							],
							"Severity": "Warning",
							"Resolution": "Enter a valid device ID and retry the operation."
						}
					]
				}
			}`))
		}

		return true
	}
	return false
}

func mockGetTemplateByNameAPIs(r *http.Request, w http.ResponseWriter) bool {
	if strings.Contains(r.URL.RawQuery, "ValidEmptyTemplate") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": []
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "ValidSingleTemplate") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
				{
					"Id": 26,
					"Name": "ValidSingleTemplate",
					"Description": "This is a test template update1"
				}
			]
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "ValidMultipleTemplate") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
				{
					"Id": 26,
					"Name": "ValidMultipleTemplate1",
					"Description": "This is a test template update1"
				},
				{
					"Id": 27,
					"Name": "ValidMultipleTemplate2",
					"Description": "This is a test template update2"
				}
			]
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "UnauthorisedTemplate") && r.Method == "GET" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CGEN1001",
						"RelatedProperties": [],
						"Message": "Unable to complete the operation because the user is not authenticated.",
						"MessageArgs": [],
						"Severity": "Critical",
						"Resolution": "Log in with valid credentials and retry the operation."
					}
				]
			}
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "UnmarshalErrTemplate") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
		}`))
		return true
	}
	return false
}

func mockGetTemplateByIDAPIs(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == TemplateAPI+"(23)") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"Id": 23,
			"Name": "Test Template",
			"Description": "This is a test template"
		}`))
		return true
	}
	if (r.URL.Path == TemplateAPI+"(24)") && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CTEM1026",
						"RelatedProperties": [],
						"Message": "Unable to process the request because the template ID 24 provided is invalid.",
						"MessageArgs": [
							"36"
						],
						"Severity": "Informational",
						"Resolution": "Enter a valid template ID and retry the operation. For information about valid template IDs. Refer to the API Guide or Product Guide available on the support site."
					}
				]
			}
		}`))
		return true
	}
	if (r.URL.Path == TemplateAPI+"(25)") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"Id": 23
			"Name": "Test Template",
			"Description": "This is a test template"
		}`))
		return true
	}
	return false
}

func mockGetIdentityPoolAPIs(r *http.Request, w http.ResponseWriter) bool {
	if r.URL.Path == IdentityPoolAPI && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
				{
					"Id": 1,
					"Name": "IdPool1"
				}
			]
			
		}`))
		return true
	}
	if r.URL.Path == fmt.Sprintf(IdentityPoolAPI+"(%d)", 123) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"Id": 123,
			"Name": "IdPool1"		
		}`))
		return true
	}
	if r.URL.Path == fmt.Sprintf(IdentityPoolAPI+"(%d)", 124) && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CTEM9046",
						"RelatedProperties": [],
						"Message": "Unable to process the request because the Identity Pool ID 124 provided is invalid.",
						"MessageArgs": [
							"5"
						],
						"Severity": "Warning",
						"Resolution": "Enter a valid Identity Pool ID and retry the operation."
					}
				]
			}
		}`))
		return true
	}

	return false
}

func mockGetNetworkAttributesAPI(r *http.Request, w http.ResponseWriter) bool {
	if r.URL.Path == TemplateAPI+"(50)/Views(4)/AttributeViewDetails" && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"Id": 4,
			"AttributeGroups": [
			  {
				"GroupNameId": 1005,
				"DisplayName": "NicBondingTechnology",
				"SubAttributeGroups": [
				  
				],
				"Attributes": [
				  {
					"CustomId": 0,
					"DisplayName": "Nic Bonding Technology",
					"Value": "NoTeaming"
				  }
				]
			  },
			  {
				"GroupNameId": 1001,
				"DisplayName": "NICModel",
				"SubAttributeGroups": [
				  {
					"GroupNameId": 1,
					"DisplayName": "Integrated NIC 1",
					"SubAttributeGroups": [
					  {
						"GroupNameId": 1,
						"DisplayName": "Port ",
						"SubAttributeGroups": [
						  {
							"GroupNameId": 1,
							"DisplayName": "Partition ",
							"SubAttributeGroups": [
							  
							],
							"Attributes": [
							  {
								"CustomId": 1049,
								"DisplayName": "NIC Bonding Enabled",
								"Value": "false"
							  },
							  {
								"CustomId": 1049,
								"DisplayName": "Vlan Tagged",
								"Value": "10133"
							  },
							  {
								"CustomId": 1049,
								"DisplayName": "Vlan UnTagged",
								"Value": "0"
							  }
							]
						  }
						],
						"Attributes": [
						  
						]
					  }
					],
					"Attributes": [
					  
					]
				  }
				],
				"Attributes": [
				  
				]
			  },
			  {
				"GroupNameId": 1001,
				"DisplayName": "TestAttributeGroup",
				"SubAttributeGroups": [],
				"Attributes": []
			  }
			]
		  }`))

		return true
	} else if r.URL.Path == TemplateAPI+"(51)/Views(4)/AttributeViewDetails" && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`
		{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CGEN1004",
						"RelatedProperties": [],
						"Message": "Unable to complete the operation because the value provided for TemplateId is invalid.",
						"MessageArgs": [
							"TemplateId"
						],
						"Severity": "Critical",
						"Resolution": "Enter a valid value and retry the operation."
					}
				]
			}
		}
		`))
	}
	return false
}

func mockUpdateTemplateAPIs(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == TemplateAPI+"(123)") && r.Method == "PUT" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`123`))
		return true
	} else if (r.URL.Path == TemplateAPI+"(124)") && r.Method == "PUT" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CGEN1002",
						"RelatedProperties": [],
						"Message": "Unable to complete the operation because the requested URI is invalid.",
						"MessageArgs": [],
						"Severity": "Critical",
						"Resolution": "Enter a valid URI and retry the operation."
					}
				]
			}
		}`))
		return true
	}
	return false
}

func mockGetTemplateAttributesAPIs(r *http.Request, w http.ResponseWriter) bool {
	if r.URL.Path == TemplateAPI+"(31)/AttributeDetails" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"AttributeGroups": [
			  {
				"GroupNameId": 1,
				"DisplayName": "BIOS",
				"SubAttributeGroups": [
				  {
					"GroupNameId": 123,
					"DisplayName": "BIOS Boot Settings",
					"SubAttributeGroups": [
					  
					],
					"Attributes": [
					  {
						"AttributeId": 110,
						"DisplayName": "Boot Sequence",
						"Value": "HardDisk.List.1-1",
						"isIgnored": false
					  }
					]
				  }
				],
				"Attributes": [
				  
				]
			  },
			  {
				"GroupNameId": 2,
				"DisplayName": "NIC",
				"SubAttributeGroups": [
				  {
					"GroupNameId": 124,
					"DisplayName": "NIC.Integrated.1-1-1",
					"SubAttributeGroups": [
					  {
						"GroupNameId": 456,
						"DisplayName": "iSCSI General Parameters",
						"SubAttributeGroups": [
						  
						],
						"Attributes": [
						  {
							"AttributeId": 120,
							"DisplayName": "Boot to Target",
							"Value": "Enabled",
							"IsIgnored": false
						  }
						]
					  }
					],
					"Attributes": [
					  
					]
				  }
				],
				"Attributes": [
				  
				]
			  }
			]
		  }`))
		return true
	}
	if (r.URL.Path == TemplateAPI+"(32)/AttributeDetails" ||
		r.URL.Path == TemplateAPI+"(33)/AttributeDetails" ||
		r.URL.Path == TemplateAPI+"(34)/AttributeDetails") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"AttributeGroups" :[
				{
					"GroupNameId": 1,
					"DisplayName": "BIOS",
					"SubAttributeGroups": [
						{
							"GroupNameId": 32700,
							"DisplayName": "BIOS Boot Settings",
							"SubAttributeGroups": [],
							"Attributes": [
								{
								"AttributeId": 110,
								"DisplayName": "Boot Sequence",
								"Value": "HardDisk.List.1-1",
								"IsIgnored": false

							}
							]
						}
					],
					"Attributes": []
				}
			]
		}`))
		return true
	}
	if (r.URL.Path == TemplateAPI+"(35)/AttributeDetails") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"AttributeGroups": [
			  {
				"GroupNameId": 1,
				"DisplayName": "BIOS",
				"SubAttributeGroups": [
				  {
					"GroupNameId": 123,
					"DisplayName": "BIOS Boot Settings",
					"SubAttributeGroups": [
					  
					],
					"Attributes": [
					  {
						"AttributeId": 110,
						"DisplayName": "Boot Sequence",
						"Value": "HardDisk.List.1-1",
						"isIgnored": false
					  }
					]
				  }
				],
				"Attributes": [
				  
				]
			  },
			  {
				"GroupNameId": 2,
				"DisplayName": "NIC",
				"SubAttributeGroups": [
				  {
					"GroupNameId": 124,
					"DisplayName": "NIC.Integrated.1-1-1",
					"SubAttributeGroups": [
					  {
						"GroupNameId": 456,
						"DisplayName": "iSCSI General Parameters",
						"SubAttributeGroups": [
						  
						],
						"Attributes": [
						  {
							"AttributeId": 120,
							"DisplayName": "Boot to Target",
							"Value": "Enabled",
							"IsIgnored": false
						  }
						]
					  }
					],
					"Attributes": [
					  
					]
				  }
				],
				"Attributes": [
				  
				]
			  }
			]
		  }`))
		return true
	}
	if (r.URL.Path == TemplateAPI+"(36)/AttributeDetails") && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CGEN1004",
						"RelatedProperties": [],
						"Message": "Unable to complete the operation because the value provided for TemplateId is invalid.",
						"MessageArgs": [
							"36"
						],
						"Severity": "Informational",
						"Resolution": "Enter a valid value and retry the operation."
					}
				]
			}
		}`))
	}
	if (r.URL.Path == TemplateAPI+"(37)/AttributeDetails") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{`))
	}
	return false
}

func mockTrackTemplateCreationAPIs(r *http.Request, w http.ResponseWriter, jobRetries *int) bool {
	if (r.URL.Path == TemplateAPI+"(25)") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"Id": 25,
			"Name": "TestTemplate",
			"Description": "This is test template",
			"Status": 2060
		}`))
		return true
	}
	if (r.URL.Path == TemplateAPI+"(26)") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"Id": 26,
			"Name": "Temp12",
			"Description": "This is temp12",
			"Status": 2070
		}`))
		return true
	}

	if (r.URL.Path == TemplateAPI+"(27)") && r.Method == "GET" {
		*jobRetries++
		w.WriteHeader(http.StatusOK)
		if *jobRetries == 3 {
			*jobRetries = 0
			w.Write([]byte(`{
				"Id": 27,
				"Name": "Temp12",
				"Description": "This is temp12",
				"Status": 2060
			}`))
		} else {
			w.Write([]byte(`{
				"Id": 27,
				"Name": "Temp12",
				"Description": "This is temp12",
				"Status": 2050
			}`))
		}
		return true
	}

	if (r.URL.Path == TemplateAPI+"(28)") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"Id": 28,
			"Name": "TestTemplate",
			"Description": "This is TestTemplate",
			"Status": 2050
		}`))
		return true
	}
	return false
}

func mockUpdateNetworkConfigAPI(r *http.Request, w http.ResponseWriter) bool {

	if (r.URL.Path == UpdateNetworkConfigAPI) && r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "NoTeaming") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`1`))
		} else if strings.Contains(string(body), "None") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CGEN6002",
							"RelatedProperties": [],
							"Message": "Unable to complete the request because the input value for VlanAttributes is missing or an invalid value is entered.",
							"MessageArgs": [
								"VlanAttributes"
							],
							"Severity": "Critical",
							"Resolution": "Enter a valid value and retry the operation."
						}
					]
				}
			}`))
		}

		return true
	}
	return false
}

func mockGroupServiceAPIs(r *http.Request, w http.ResponseWriter) bool {
	if r.URL.Path == GroupAPI {
		if r.URL.RawQuery == "Name=valid_group1" && r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"@odata.context": "/api/$metadata#Collection(GroupService.Group)",
				"@odata.count": 1,
				"value": [
					{
						"Id": 1011,
						"Name": "Linux Servers"
					}
				]
			}`))
			return true
		} else if r.URL.RawQuery == "Name=valid_group2" && r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"@odata.context": "/api/$metadata#Collection(GroupService.Group)",
				"@odata.count": 1,
				"value": [
					{
						"Id": 1012,
						"Name": "IDRAC Servers"
					}
				]
			}`))
			return true
		} else if r.URL.RawQuery == "Name=invalid_group1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"@odata.context": "/api/$metadata#Collection(GroupService.Group)",
				"@odata.count": 0,
				"value": []
			}`))
			return true
		}
	}

	if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, 1011)) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(DeviceService.Device)",
			"@odata.count": 1,
			"value": [
				{
					"@odata.type": "#DeviceService.Device",
					"@odata.id": "/api/DeviceService/Devices(10337)",
					"Id": 10337,
					"DeviceServiceTag": "SvcTag-1"
				}
			]
		}`))
		return true
	} else if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, 1012)) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(DeviceService.Device)",
			"@odata.count": 1,
			"value": [
				{
					"@odata.type": "#DeviceService.Device",
					"@odata.id": "/api/DeviceService/Devices(10338)",
					"Id": 10338,
					"DeviceServiceTag": "SvcTag-2"
				}
			]
		}`))
		return true
	} else if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, 1013)) && r.Method == "GET" && !strings.Contains(r.URL.RawQuery, "skip=1&top=1") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
				"@odata.context": "/api/$metadata#Collection(DeviceService.Device)",
				"@odata.count": 2,
				"value": [
					{
						"@odata.type": "#DeviceService.Device",
						"@odata.id": "/api/DeviceService/Devices(10337)",
						"Id": 10337,
						"DeviceServiceTag": "SvcTag-1"
					}
				],
				"@odata.nextLink": "/api/GroupService/Groups(1013)/Devices?skip=1&top=1"
			}`))
		return true
	} else if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, 1013)) && r.Method == "GET" && strings.Contains(r.URL.RawQuery, "skip=1&top=1") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(DeviceService.Device)",
			"@odata.count": 1,
			"value": [
				{
					"@odata.type": "#DeviceService.Device",
					"@odata.id": "/api/DeviceService/Devices(10338)",
					"Id": 10338,
					"DeviceServiceTag": "SvcTag-2"
				}
			]
		}`))
		return true
	} else if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, 1014)) && r.Method == "GET" && !strings.Contains(r.URL.RawQuery, "skip=1&top=1") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
				"@odata.context": "/api/$metadata#Collection(DeviceService.Device)",
				"@odata.count": 2,
				"value": [
					{
						"@odata.type": "#DeviceService.Device",
						"@odata.id": "/api/DeviceService/Devices(10337)",
						"Id": 10337,
						"DeviceServiceTag": "SvcTag-1"
					}
				],
				"@odata.nextLink": "/api/GroupService/Groups(1014)/Devices?skip=1&top=1"
			}`))
		return true
	} else if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, 1014)) && r.Method == "GET" && strings.Contains(r.URL.RawQuery, "skip=1&top=1") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CGRP9008",
						"RelatedProperties": [],
						"Message": "Unable to complete the operation because data entered in the field Group Id is invalid.",
						"MessageArgs": [
							"Group Id"
						],
						"Severity": "Warning",
						"Resolution": "Enter valid data in the field and retry the operation. For information about valid data permitted in a field, see the OpenManage Enterprise-Modular User's Guide available on the support site."
					}
				]
			}
		}`))
		return true
	} else if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, 1015)) && r.Method == "GET" && !strings.Contains(r.URL.RawQuery, "skip=1&top=1") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
				"@odata.context": "/api/$metadata#Collection(DeviceService.Device)",
				"@odata.count": 2,
				"value": [
					{
						"@odata.type": "#DeviceService.Device",
						"@odata.id": "/api/DeviceService/Devices(10337)",
						"Id": 10337,
						"DeviceServiceTag": "SvcTag-1",
					},
				],
				"@odata.nextLink": "/api/GroupService/Groups(1015)/Devices?skip=1&top=1"
			}`))
		return true
	} else if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, 1016)) && r.Method == "GET" && !strings.Contains(r.URL.RawQuery, "skip=1&top=1") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
				"@odata.context": "/api/$metadata#Collection(DeviceService.Device)",
				"@odata.count": 2,
				"value": [
					{
						"@odata.type": "#DeviceService.Device",
						"@odata.id": "/api/DeviceService/Devices(10337)",
						"Id": 10337,
						"DeviceServiceTag": "SvcTag-1"
					}
				],
				"@odata.nextLink": "/api/GroupService/Groups(1016)/Devices?skip=1&top=1"
			}`))
		return true
	} else if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, 1016)) && r.Method == "GET" && strings.Contains(r.URL.RawQuery, "skip=1&top=1") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(DeviceService.Device)",
			"@odata.count": 1,
			"value": [
				{
					"@odata.type": "#DeviceService.Device",
					"@odata.id": "/api/DeviceService/Devices(10338)",
					"Id": 10338,
					"DeviceServiceTag": "SvcTag-2",
				},
			]
		}`))
		return true
	} else if r.URL.Path == (fmt.Sprintf(GroupServiceDevicesAPI, -1)) && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CGRP9008",
						"RelatedProperties": [],
						"Message": "Unable to complete the operation because data entered in the field Group Id is invalid.",
						"MessageArgs": [
							"Group Id"
						],
						"Severity": "Warning",
						"Resolution": "Enter valid data in the field and retry the operation. For information about valid data permitted in a field, see the OpenManage Enterprise-Modular User's Guide available on the support site."
					}
				]
			}
		}`))
		return true
	}

	return false
}

func mockGetServerProfileInfoByTemplateNameAPIs(r *http.Request, w http.ResponseWriter) bool {
	if strings.Contains(r.URL.RawQuery, "ValidEmptyProfileTemplateName") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": []
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "ValidSingleProfileTemplateName") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
				{
					"Id": 10848,
            		"ProfileName": "Profile from template 'test_deployment' 00001",
            		"TemplateId": 585
				}
			]
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "ValidMultipleProfileTemplateName") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
				{
					"Id": 10848,
            		"ProfileName": "Profile from template 'test_deployment' 00001",
            		"TemplateId": 585
				},
				{
					"Id": 10849,
            		"ProfileName": "Profile from template 'test_deployment' 00002",
            		"TemplateId": 585
				}
			]
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "UnauthorisedProfileTemplateName") && r.Method == "GET" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CGEN1001",
						"RelatedProperties": [],
						"Message": "Unable to complete the operation because the user is not authenticated.",
						"MessageArgs": [],
						"Severity": "Critical",
						"Resolution": "Log in with valid credentials and retry the operation."
					}
				]
			}
		}`))
		return true
	}
	if strings.Contains(r.URL.RawQuery, "UnmarshalErrProfileTemplateName") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"value": [
		}`))
		return true
	}
	return false
}

func mockDeployAPIs(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == DeployAPI) && r.Method == "POST" {
		requestStruct := models.OMETemplateDeployRequest{}
		err := json.NewDecoder(r.Body).Decode(&requestStruct)
		if err != nil {
			return false
		}
		if requestStruct.ID == 1 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`1234`))
		} else if requestStruct.ID == 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CTEM9021",
							"RelatedProperties": [],
							"Message": "Unable to deploy the template test_deployment because 100.96.24.28 has a profile assigned.",
							"MessageArgs": [
								"test_deployment",
								"100.96.24.28 has a profile assigned"
							],
							"Severity": "Warning",
							"Resolution": "Review the reason and initiate necessary resolution."
						}
					]
				}
			}`))
		}
		return true
	}
	return false
}

func mockUnassignProfileAPI(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == UnAssignProfileAPI) && r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "10850") || strings.Contains(string(body), "10852") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`10860`))
		} else if strings.Contains(string(body), "10853") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`10861`))
		} else if strings.Contains(string(body), "10851") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CGEN6002",
							"RelatedProperties": [],
							"Message": "Unable to complete the request because the input value for ProfileIds is missing or an invalid value is entered.",
							"MessageArgs": [
								"ProfileIds"
							],
							"Severity": "Critical",
							"Resolution": "Enter a valid value and retry the operation."
						}
					]
				}
			}`))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CGEN6002",
							"RelatedProperties": [],
							"Message": "Unable to complete the request because the input value for ProfileIds is missing or an invalid value is entered.",
							"MessageArgs": [
								"ProfileIds"
							],
							"Severity": "Critical",
							"Resolution": "Enter a valid value and retry the operation."
						}
					]
				}
			}`))
		}

		return true
	}
	return false
}

func mockDeleteProfileAPI(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == DeleteProfileAPI) && r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "10850") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`0`))
		} else if strings.Contains(string(body), "10852") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CGEN6002",
							"RelatedProperties": [],
							"Message": "Unable to complete the request because the input value for ProfileIds is invalid.",
							"MessageArgs": [
								"ProfileIds"
							],
							"Severity": "Critical",
							"Resolution": "Enter a valid value and retry the operation."
						}
					]
				}
			}`))
		}

		return true
	}
	return false
}

func mockBaselineAPIs(r *http.Request, w http.ResponseWriter) bool {

	if (r.URL.Path == BaselineAPI) && r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "TestAccCreateBaselineFailure") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CTEM1026",
							"RelatedProperties": [],
							"Message": "Unable to process the request because the template ID -1 provided is invalid.",
							"MessageArgs": [
								"3226"
							],
							"Severity": "Informational",
							"Resolution": "Enter a valid template ID and retry the operation. For information about valid template IDs. Refer to the API Guide or Product Guide available on the support site."
						}
					]
				}
			}`))
		} else if strings.Contains(string(body), "TestAccCreateBaseline") {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
			"Id": 21,
			"Name": "TestAccCreateBaseline",
			"Description": "Test Acc Description for create baseline",
			"LastRun": null,
			"TemplateId": 326,
			"TemplateName": null,
			"TemplateType": 0,
			"TaskId": 0,
			"PercentageComplete": null,
			"TaskStatus": 0,
			"ConfigComplianceSummary": null,
			"BaselineTargets": [
				{
					"Id": 10093,
					"Type": {
						"Id": 1,
						"Name": "DEVICE"
					}
				},
				{
					"Id": 10104,
					"Type": {
						"Id": 1,
						"Name": "DEVICE"
					}
				}
			],
			"NotificationSettings": {
				"NotificationType": "NOTIFY_ON_SCHEDULE",
				"Schedule": {
					"RunNow": false,
					"RunLater": false,
					"Cron": "0 00 00 * * ? *",
					"StartTime": null,
					"EndTime": null
				},
				"EmailAddresses": [
					"test@testdell.com"
				],
				"OutputFormat": "HTML"
			}
		}`))
		}
		return true
	}
	if r.URL.Path == fmt.Sprintf(BaselineAPI+"(%d)", 101) && r.Method == "PUT" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CTEM1026",
							"RelatedProperties": [],
							"Message": "Unable to process the request because the template ID -1 provided is invalid.",
							"MessageArgs": [
								"3226"
							],
							"Severity": "Informational",
							"Resolution": "Enter a valid template ID and retry the operation. For information about valid template IDs. Refer to the API Guide or Product Guide available on the support site."
						}
					]
				}
			}`))
	}
	if r.URL.Path == fmt.Sprintf(BaselineAPI+"(%d)", 100) && r.Method == "PUT" {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
		"Id": 100,
		"Name": "TestAccCreateBaseline",
		"Description": "Test Acc Description for create baseline",
		"LastRun": null,
		"TemplateId": 326,
		"TemplateName": null,
		"TemplateType": 0,
		"TaskId": 0,
		"PercentageComplete": null,
		"TaskStatus": 0,
		"ConfigComplianceSummary": null,
		"BaselineTargets": [
			{
				"Id": 10093,
				"Type": {
					"Id": 1,
					"Name": "DEVICE"
				}
			},
			{
				"Id": 10104,
				"Type": {
					"Id": 1,
					"Name": "DEVICE"
				}
			}
		],
		"NotificationSettings": {
			"NotificationType": "NOTIFY_ON_SCHEDULE",
			"Schedule": {
				"RunNow": false,
				"RunLater": false,
				"Cron": "0 00 00 * * ? *",
				"StartTime": null,
				"EndTime": null
			},
			"EmailAddresses": [
				"test@testdell.com"
			],
			"OutputFormat": "HTML"
		}
	}`))
		return true
	}
	if r.URL.Path == BaseLineRemoveAPI && r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "10001") {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`10001`))
		}
		if strings.Contains(string(body), "10002") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CGEN1004",
							"RelatedProperties": [],
							"Message": "Unable to complete the operation because the value provided for {0} is invalid.",
							"MessageArgs": [],
							"Severity": "Critical",
							"Resolution": "Enter a valid value and retry the operation."
						}
					]
				}
			}`))
		}

		return true
	}
	return false
}

func mockGetBaselineByIDAPI(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == BaselineAPI+"(1)") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"Id": 1,
			"Name": "Baseline Name",
			"Description": "Description 1.",
			"TemplateId": 326,
			"TemplateName": "test-compliance-template",
			"TemplateType": 2,
			"TaskId": 11988,
			"PercentageComplete": null,
			"TaskStatus": 0,
			"LastRun": null,
			"BaselineTargets": [
				{
					"Id": 10093,
					"Type": {
						"Id": 1000,
						"Name": "DEVICE"
					}
				},
				{
					"Id": 10104,
					"Type": {
						"Id": 1000,
						"Name": "DEVICE"
					}
				}
			],
			"ConfigComplianceSummary": {
				"ComplianceStatus": "OK",
				"NumberOfCritical": 0,
				"NumberOfWarning": 0,
				"NumberOfNormal": 0,
				"NumberOfIncomplete": 0
			},
			"NotificationSettings": {
				"NotificationType": "NOTIFY_ON_NON_COMPLIANCE",
				"Schedule": {
					"Cron": "0 00 00 * * ? *"
				},
				"EmailAddresses": [
					"naveen.patil@dell.com"
				],
				"OutputFormat": "html"
			},
			"DeviceConfigComplianceReports@odata.navigationLink": "/api/TemplateService/Baselines(20)/DeviceConfigComplianceReports"
		}`))
		return true
	}
	if (r.URL.Path == BaselineAPI+"(-1)") && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": "Base.1.0.GeneralError",
				"message": "A general error has occurred. See ExtendedInfo for more information.",
				"@Message.ExtendedInfo": [
					{
						"MessageId": "CGEN1008",
						"RelatedProperties": [],
						"Message": "Unable to process the request because an error occurred.",
						"MessageArgs": [],
						"Severity": "Critical",
						"Resolution": "Retry the operation. If the issue persists, contact your system administrator."
					}
				]
			}
		}`))
		return true
	}
	return false
}

func mockGetBaselineDevComplianceReportByIDAPI(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == fmt.Sprintf(BaselineDeviceComplianceReportsAPI, 14)) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(TemplateService.DeviceConfigComplianceReports)",
			"@odata.count": 2,
			"value": [
				{
					"@odata.type": "#TemplateService.DeviceConfigComplianceReports",
					"@odata.id": "/api/TemplateService/Baselines(14)/DeviceConfigComplianceReports(11803)",
					"Id": 11803,
					"DeviceName": "WIN-02GODDHDJTC",
					"IpAddress": "100.68.168.50",
					"IpAddresses": [
						"100.68.168.50"
					],
					"Model": "PowerEdge MX840c",
					"ServiceTag": "MX84002",
					"ComplianceStatus": 2,
					"DeviceType": 1000,
					"InventoryTime": "2022-11-09 00:01:14.619974",
					"DeviceComplianceDetails": {
						"@odata.id": "/api/TemplateService/Baselines(14)/DeviceConfigComplianceReports(11803)/DeviceComplianceDetails"
					}
				},
				{
					"@odata.type": "#TemplateService.DeviceConfigComplianceReports",
					"@odata.id": "/api/TemplateService/Baselines(14)/DeviceConfigComplianceReports(10337)",
					"Id": 10337,
					"DeviceName": "WIN-MX740.wacdev.com",
					"IpAddress": null,
					"IpAddresses": [],
					"Model": "PowerEdge MX740c",
					"ServiceTag": "6H6GNX2",
					"ComplianceStatus": 3,
					"DeviceType": 1000,
					"InventoryTime": null,
					"DeviceComplianceDetails": {
						"@odata.id": "/api/TemplateService/Baselines(14)/DeviceConfigComplianceReports(10337)/DeviceComplianceDetails"
					}
				}
			]
		}`))
		return true
	} else if (r.URL.Path == fmt.Sprintf(BaselineDeviceComplianceReportsAPI, -1)) && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(TemplateService.DeviceConfigComplianceReports)",
			"@odata.count": 0,
			"value": []
		}`))
		return true
	}
	return false
}

func mockGetBaselineDevAttrComplianceReportByIDAPI(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == fmt.Sprintf(BaselineDeviceAttrComplianceReportsAPI, 14, 11803)) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#TemplateService.DeviceComplianceDetail",
			"@odata.type": "#TemplateService.DeviceComplianceDetail",
			"@odata.id": "/api/TemplateService/Baselines(14)/DeviceConfigComplianceReports(11803)/DeviceComplianceDetails",
			"DeviceId": 11803,
			"DeviceName": "WIN-02GODDHDJTC",
			"BaselineId": 14,
			"BaselineName": "Baseline Name 1-updated",
			"TemplateId": 326,
			"TemplateName": "test-compliance-template",
			"ComplianceAttributeGroups": [
				{
					"GroupNameId": 2,
					"DisplayName": "LifecycleController",
					"ComplianceStatus": 2,
					"ComplianceReason": "One or more attributes on the target device(s) does not match the compliance template.",
					"ComplianceSubAttributeGroups": [
						{
							"GroupNameId": 33,
							"DisplayName": "Lifecycle Controller Attributes",
							"ComplianceStatus": 2,
							"ComplianceReason": "One or more attributes on the target device(s) does not match the compliance template.",
							"ComplianceSubAttributeGroups": [],
							"Attributes": [
								{
									"AttributeId": 721728,
									"CustomId": 0,
									"DisplayName": "LCAttributes 1 Automatic Backup Feature",
									"Description": null,
									"Value": null,
									"ExpectedValue": "Disabled",
									"ComplianceStatus": 2,
									"ComplianceReason": "Missing template attribute value."
								},
								{
									"AttributeId": 721727,
									"CustomId": 0,
									"DisplayName": "LCAttributes 1 Automatic Update Feature",
									"Description": null,
									"Value": null,
									"ExpectedValue": "Disabled",
									"ComplianceStatus": 1,
									"ComplianceReason": "All attributes on the target device(s) match the compliance template."
								}
							]
						}
					],
					"Attributes": []
				}
			]
		}`))
		return true
	} else if (r.URL.Path == fmt.Sprintf(BaselineDeviceAttrComplianceReportsAPI, -1, -1)) && r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#TemplateService.DeviceComplianceDetail",
			"@odata.type": "#TemplateService.DeviceComplianceDetail",
			"@odata.id": "/api/TemplateService/Baselines(-1)/DeviceConfigComplianceReports(-1)/DeviceComplianceDetails",
			"DeviceId": -1,
			"DeviceName": null,
			"BaselineId": -1,
			"BaselineName": null,
			"TemplateId": 0,
			"TemplateName": null,
			"ComplianceAttributeGroups": []
		}`))
		return true
	}
	return false
}

func mockCloneTemplateAPI(r *http.Request, w http.ResponseWriter) bool {
	if (r.URL.Path == CloneTemplateAPI) && r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "dep-dep-template") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`` + fmt.Sprint(TestClonedDeploymentTemplateID) + ``))
		} else if strings.Contains(string(body), "dep-comp-template") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`` + fmt.Sprint(TestClonedComplianceTemplateID) + ``))
		} else if strings.Contains(string(body), "comp-comp-template") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`` + fmt.Sprint(TestClonedComplianceTemplateID) + ``))
		} else if strings.Contains(string(body), "test-invalid-viewtype-id") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`` + fmt.Sprint(TestClonedComplianceTemplateID) + ``))
		} else if strings.Contains(string(body), "test-invalid-template-id") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CTEM9022",
							"RelatedProperties": [],
							"Message": "Unable to clone the template clone example because Source template does not exist..",
							"MessageArgs": [
								"clone example",
								"Source template does not exist."
							],
							"Severity": "Warning",
							"Resolution": "Review the reason and initiate necessary resolution."
						}
					]
				}
			}`))
		} else if strings.Contains(string(body), "test-existing-template-name") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": "Base.1.0.GeneralError",
					"message": "A general error has occurred. See ExtendedInfo for more information.",
					"@Message.ExtendedInfo": [
						{
							"MessageId": "CTEM1027",
							"RelatedProperties": [],
							"Message": "Unable to create the template because the template name test-existing-template-name already exists.",
							"MessageArgs": [
								"clone example1"
							],
							"Severity": "Warning",
							"Resolution": "Enter a unique template name and retry the operation."
						}
					]
				}
			}`))
		}
		return true
	}
	return false
}

func mockNetworkVlanAPI(r *http.Request, w http.ResponseWriter) bool {
	if r.URL.Path == VlanNetworksAPI && r.Method == "GET" {
		if r.URL.RawQuery == "" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"@odata.context": "/api/$metadata#Collection(NetworkConfigurationService.Network)",
				"@odata.count": 2,
				"value": [
					{
						"Id": 1234,
						"Name": "VLAN1"
					}
				],
				"@odata.nextLink": "/api/NetworkConfigurationService/Networks?skip=1&top=1"
			}`))
			return true
		} else if strings.Contains(r.URL.RawQuery, "skip=1&top=1") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"@odata.context": "/api/$metadata#Collection(NetworkConfigurationService.Network)",
				"@odata.count": 2,
				"value": [
					{
						"Id": 1235,
						"Name": "VLAN2"
					}
				]
			}`))
			return true
		}
	}
	return false
}

func mockGetBaselineByNameAPI(r *http.Request, w http.ResponseWriter) bool {
	// first page
	if (r.URL.Path == BaselineAPI && r.URL.RawQuery == "") && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(TemplateService.Baseline)",
			"@odata.count": 3,
			"value": [
				{
					"@odata.type": "#TemplateService.Baseline",
					"@odata.id": "/api/TemplateService/Baselines(162)",
					"Id": 162,
					"Name": "TestAccBaseline1",
					"Description": null,
					"TemplateId": 745,
					"TemplateName": "test_acc_compliance_template",
					"TemplateType": 2,
					"TaskId": 12399,
					"PercentageComplete": "100",
					"TaskStatus": 2060,
					"LastRun": "2022-11-18 04:32:16.838",
					"BaselineTargets": [
						{
							"Id": 12152,
							"Type": {
								"Id": 1000,
								"Name": "DEVICE"
							}
						}
					],
					"ConfigComplianceSummary": {
						"ComplianceStatus": "NOT_INVENTORIED",
						"NumberOfCritical": 0,
						"NumberOfWarning": 0,
						"NumberOfNormal": 0,
						"NumberOfIncomplete": 1
					},
					"NotificationSettings": {
						"NotificationType": "NOTIFY_ON_SCHEDULE",
						"Schedule": {
							"Cron": "0 10 5 * * ? *",
							"StartTime": "2022-11-18 05:10:59.780",
							"EndTime": null
						},
						"EmailAddresses": [
							"dell@dell.com"
						],
						"OutputFormat": "html"
					},
					"DeviceConfigComplianceReports@odata.navigationLink": "/api/TemplateService/Baselines(162)/DeviceConfigComplianceReports"
				}
			],
			"@odata.nextLink": "/api/TemplateService/Baselines?skip=1&top=1"
		}`))
		return true
	}
	// second page
	if (r.URL.Path == BaselineAPI && strings.EqualFold(r.URL.RawQuery, "skip=1&top=1")) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(TemplateService.Baseline)",
			"@odata.count": 3,
			"value": [
				{
					"@odata.type": "#TemplateService.Baseline",
					"@odata.id": "/api/TemplateService/Baselines(163)",
					"Id": 163,
					"Name": "TestAccBaseline2",
					"Description": null,
					"TemplateId": 745,
					"TemplateName": "test_acc_compliance_template",
					"TemplateType": 2,
					"TaskId": 12400,
					"PercentageComplete": "100",
					"TaskStatus": 2060,
					"LastRun": "2022-11-18 04:32:16.838",
					"BaselineTargets": [
						{
							"Id": 12152,
							"Type": {
								"Id": 1000,
								"Name": "DEVICE"
							}
						}
					],
					"ConfigComplianceSummary": {
						"ComplianceStatus": "NOT_INVENTORIED",
						"NumberOfCritical": 0,
						"NumberOfWarning": 0,
						"NumberOfNormal": 0,
						"NumberOfIncomplete": 1
					},
					"NotificationSettings": {
						"NotificationType": "NOTIFY_ON_SCHEDULE",
						"Schedule": {
							"Cron": "0 10 5 * * ? *",
							"StartTime": "2022-11-18 05:10:59.780",
							"EndTime": null
						},
						"EmailAddresses": [
							"dell@dell.com"
						],
						"OutputFormat": "html"
					},
					"DeviceConfigComplianceReports@odata.navigationLink": "/api/TemplateService/Baselines(163)/DeviceConfigComplianceReports"
				}
			],
			"@odata.nextLink": "/api/TemplateService/Baselines?skip=2&top=1"
		}`))
		return true
	}
	// third page
	if (r.URL.Path == BaselineAPI && strings.EqualFold(r.URL.RawQuery, "skip=2&top=1")) && r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"@odata.context": "/api/$metadata#Collection(TemplateService.Baseline)",
			"@odata.count": 3,
			"value": [
				{
					"@odata.type": "#TemplateService.Baseline",
					"@odata.id": "/api/TemplateService/Baselines(164)",
					"Id": 164,
					"Name": "TestAccBaseline3",
					"Description": null,
					"TemplateId": 745,
					"TemplateName": "test_acc_compliance_template",
					"TemplateType": 2,
					"TaskId": 12401,
					"PercentageComplete": "100",
					"TaskStatus": 2050,
					"LastRun": "2022-11-18 04:32:16.838",
					"BaselineTargets": [
						{
							"Id": 12152,
							"Type": {
								"Id": 1000,
								"Name": "DEVICE"
							}
						}
					],
					"ConfigComplianceSummary": {
						"ComplianceStatus": "NOT_INVENTORIED",
						"NumberOfCritical": 0,
						"NumberOfWarning": 0,
						"NumberOfNormal": 0,
						"NumberOfIncomplete": 1
					},
					"NotificationSettings": {
						"NotificationType": "NOTIFY_ON_SCHEDULE",
						"Schedule": {
							"Cron": "0 10 5 * * ? *",
							"StartTime": "2022-11-18 05:10:59.780",
							"EndTime": null
						},
						"EmailAddresses": [
							"dell@dell.com"
						],
						"OutputFormat": "html"
					},
					"DeviceConfigComplianceReports@odata.navigationLink": "/api/TemplateService/Baselines(164)/DeviceConfigComplianceReports"
				}
			]
		}`))
		return true
	}
	return false
}
