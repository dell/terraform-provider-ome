
## Dell OME

|  Endpoints  |  Verb  |  Terraform Actions  |
| :-----------------------------------------------------------------: | :----: | :----------------: |
| **Certificate Datasource**                |
| <sub>/api/ApplicationService/Certificate </sub> | GET | [Read]
| **Baseline Report Datasource**                |
| <sub>/api/TemplateService/Baseline </sub> | GET | [Read]
| <sub>/api/TemplateService/Baselines(%d)/DeviceConfigComplianceReports </sub> | GET | [Read]
| <sub>/api/TemplateService/Baselines(%d)/DeviceConfigComplianceReports(%d)/DeviceComplianceDetails </sub> | GET | [Read]
| **Device Datasource**                |
| <sub>/api/DeviceService/Devices </sub> | GET | [Read]
| <sub>/api/DeviceService/Devices(%d)/InventoryDetails </sub> | GET | [Read]
| <sub>/api/DeviceService/Devices(%d)/InventoryDetails('%s') </sub> | GET | [Read]
| **Group Device Datasource**                |
| <sub>/api/GroupService/Groups </sub> | GET | [Read]
| <sub>/api/GroupService/Groups(%d)/Devices </sub> | GET | [Read]
| **Template Datasource**                |
| <sub>/api/TemplateService/Templates </sub> | GET | [Read]
| <sub>/api/TemplateService/Templates(%d)/AttributeDetails </sub> | GET | [Read]
| <sub>/api/TemplateService/Template(%d)/Views(4)/AttributeViewDetails </sub> | GET | [Read]
| **VLAN Datasource**                |
| <sub>api/NetworkConfigurationService/Networks </sub> | GET | [Read]
| **Certificate Resource**                |
| <sub>/api/ApplicationService/Actions/ApplicationService.UploadCertificate </sub> | POST | [Create]
| **CSR Resource**                |
| <sub>/api/ApplicationService/Actions/ApplicationService.GenerateCSR </sub> | POST | [Create]
| **Devices Resource**                |
| <sub>/api/DeviceService/Devices </sub> | GET | [Create, Read, Import]
| <sub>/api/DeviceService/Actions/DeviceService.RemoveDevices </sub> | POST | [Update]
| **Device Action Resource**                |
| <sub>/api/JobService/Jobs(%d) </sub> | DELETE | [Delete]
| <sub>/api/JobService/Jobs </sub> | POST | [Create]
| <sub>/api/JobService/Jobs(%d) </sub> | GET | [Create]
| <sub>/api/JobService/Jobs(%d) </sub> | GET | [Create]
| **Discovery Resource**                |
| <sub>/api/DiscoveryConfigService/DiscoveryConfigGroups(%d) </sub> | GET | [Read, Import, Create, Update]
| <sub>/api/DiscoveryConfigService/DiscoveryConfigGroups </sub> | POST | [Create, Update]
| <sub>/api/DiscoveryConfigService/Actions/DiscoveryConfigService.RemoveDiscoveryGroup </sub> | POST | [Delete]
| <sub>/api/JobService/Jobs(%d) </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs(%d)/LastExecutionDetail </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs(%d)/ExecutionHistories(%d)/ExecutionHistoryDetails </sub> | GET | [Create, Update]
| **Static Group Resource**                |
| <sub>/api/GroupService/Groups(%d) </sub> | GET | [Create, Read, Import, Update]
| <sub>/api/GroupService/Groups(%d) </sub> | DELETE | [Delete]
| <sub>/api/GroupService/Actions/GroupService.CreateGroup </sub> | POST | [Create]
| <sub>/api/GroupService/Actions/GroupService.UpdateGroup </sub> | POST | [Update]
| <sub>/api/GroupService/Actions/GroupService.AddMemberDevices </sub> | POST | [Create, Update]
| <sub>/api/GroupService/Actions/GroupService.RemoveMemberDevices </sub> | POST | [Update]
| <sub>/api/GroupService/Groups(%d)/Devices </sub> | GET | [Create, Read, Import, Update]
| **User Resource**                |
| /api/AccountService/Accounts </sub> | POST | [Create]
| /api/AccountService/Accounts({id}) </sub> | PUT | [Update]
| /api/AccountService/Accounts({id}) </sub> | GET | [Read, Import]
| /api/AccountService/Accounts({id}) </sub> | DELETE | [Delete]
| **Configuration Compliance Resource**                |
| <sub>/api/TemplateService/Baselines(%d)/DeviceConfigComplianceReports </sub> | GET | [Read]
| <sub>/api/TemplateService/Baselines </sub> | GET | [Create, Update]
| <sub>/api/TemplateService/Baselines(%d) </sub> | GET | [Create, Update, Read]
| <sub>/api/TemplateService/Actions/TemplateService.Remediate </sub> | POST | [Create, Update]
| <sub>/api/DeviceService/Devices </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs(%d) </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs(%d)/LastExecutionDetail </sub> | GET | [Create, Update]
| **Configuartion Baseline Resource**                |
| <sub>/api/TemplateService/Templates(%d) </sub> | GET | [Create, Update]
| <sub>/api/TemplateService/Templates </sub> | GET | [Create, Update]
| <sub>/api/DeviceService/Devices </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/TemplateService/Baselines </sub> | POST | [Create]
| <sub>/api/TemplateService/Baselines(%d) </sub> | PUT | [Update]
| <sub>/api/TemplateService/Baselines(%d) </sub> | GET | [Create, Read, Update]
| <sub>/api/JobService/Jobs(%d) </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs(%d)/LastExecutionDetail </sub> | GET | [Create, Update]
| **Deploy Resource**                |
| <sub>/api/TemplateService/Templates(%d) </sub> | GET | [Create]
| <sub>/api/TemplateService/Templates </sub> | GET | [Create, Import]
| <sub>/api/TemplateService/Actions/TemplateService.Deploy </sub> | POST | [Create, Update]
| <sub>/api/DeviceService/Devices </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/JobService/Jobs(%d) </sub> | GET | [Create, Update, Delete]
| <sub>/api/JobService/Jobs(%d)/LastExecutionDetail </sub> | GET | [Create, Update, Delete]
| <sub>/api/ProfileService/Profiles </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/ProfileService/Actions/ProfileService.UnassignProfiles </sub> | GET | [Update, Delete]
| **Template Resource**                |
| <sub>/api/TemplateService/TemplateViewTypes </sub> | GET | [Create]
| <sub>/api/TemplateService/TemplateTypes </sub> | GET | [Create]
| <sub>/api/TemplateService/Templates </sub> | GET | [Create, Import]
| <sub>/api/TemplateService/Templates </sub> | DELETE | [Create, Delete]
| <sub>/api/TemplateService/Templates </sub> | PUT | [Update]
| <sub>/api/TemplateService/Templates </sub> | POST | [Create]
| <sub>/api/TemplateService/Templates(%d) </sub> | GET | [Create, Read, Update]
| <sub>/api/TemplateService/Templates(%d)/AttributeDetails </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/TemplateService/Templates(%d)/Views(4)/AttributeViewDetails </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/TemplateService/Actions/TemplateService.UpdateNetworkConfig </sub> | POST | [Update]
| <sub>/api/TemplateService/Actions/TemplateService.Clone </sub> | POST | [Create]
| <sub>/api/TemplateService/Actions/TemplateService.Import </sub> | POST | [Create]
| <sub>/api/DeviceService/Devices </sub> | GET | [Create]
| <sub>/api/JobService/Jobs(%d) </sub> | GET | [Create]
| <sub>/api/JobService/Jobs(%d)/LastExecutionDetail </sub> | GET | [Create]
| <sub>/api/IdentityPoolService/IdentityPools(%d) </sub> | GET | [Read]
| <sub>/api/IdentityPoolService/IdentityPools </sub> | GET | [Update]
| <sub>/api/NetworkConfigurationService/Networks </sub> | GET | [Update]

