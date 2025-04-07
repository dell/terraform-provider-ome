# Managing State

### Overview

This guide is has some best practices for handling different scenarios around managing the terraform state. 

### Scenario 1

*The user has created a terraform resource and it is part of the state. However outside of terraform someone has manually deleted the resource.*

#### **Warning Signs this issue may have occured**

*The user, while doing a terraform apply, gets an error during the terraform read/refresh before the apply*

#### **Steps to remediate and re-create the resource**

1. Analyze the terraform read error and find the resource which was effected. 

Example Error: 
```
example.resource_2: Refreshing state... [id=1440]
╷
│ Error: Unable to read example
│
│   with example.resource_2,
│   on resource.tf line 33, in resource "example" "resource_2":
│   33: resource "example" "resource_2" {
│
```

*Verify by checking the system that the resource `example.resource_2` was actually manually removed from the target system.*

2. Run a `terraform state list`

**Ouput:**
```
example.resource_1
example.resource_2
```

3. Run the command `terraform state rm example.resource_2` to remove the resource from the terraform state. 

**Ouput:**
```
Removed example.resource_2
Successfully removed 1 resource instance(s).
```

3. Re-run `terraform apply` to re-create the deleted resource. Then you are back to a normailzed state.
