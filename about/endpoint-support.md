
## Dell OME

|  **Endpoints**  |  **Verb**  |  **Terraform Actions**  |
| :-----------------------------------------------------------------: | :----: | :----------------: |
| **Certificate (Datasource)**                |
| <sub>/api/ApplicationService/Certificate </sub> | GET | [Read]
| **Baseline Report (Datasource)**                |
| <sub>/api/TemplateService/Baseline </sub> | GET | [Read]
| <sub>/api/TemplateService/Baselines({id})/DeviceConfigComplianceReports </sub> | GET | [Read]
| <sub>/api/TemplateService/Baselines({id})/DeviceConfigComplianceReports({id})/DeviceComplianceDetails </sub> | GET | [Read]
| **Device (Datasource)**                |
| <sub>/api/DeviceService/Devices </sub> | GET | [Read]
| <sub>/api/DeviceService/Devices({id})/InventoryDetails </sub> | GET | [Read]
| <sub>/api/DeviceService/Devices({id})/InventoryDetails('%s') </sub> | GET | [Read]
| **Group Device (Datasource)**                |
| <sub>/api/GroupService/Groups </sub> | GET | [Read]
| <sub>/api/GroupService/Groups({id})/Devices </sub> | GET | [Read]
| **Template (Datasource)**                |
| <sub>/api/TemplateService/Templates </sub> | GET | [Read]
| <sub>/api/TemplateService/Templates({id})/AttributeDetails </sub> | GET | [Read]
| <sub>/api/TemplateService/Template({id})/Views(4)/AttributeViewDetails </sub> | GET | [Read]
| **VLAN (Datasource)**                |
| <sub>/api/NetworkConfigurationService/Networks </sub> | GET | [Read]
| **Certificate (Resource)**                |
| <sub>/api/ApplicationService/Actions/ApplicationService.UploadCertificate </sub> | POST | [Create]
| **CSR (Resource)**                |
| <sub>/api/ApplicationService/Actions/ApplicationService.GenerateCSR </sub> | POST | [Create]
| **Devices (Resource)**                |
| <sub>/api/DeviceService/Devices </sub> | GET | [Create, Read, Import]
| <sub>/api/DeviceService/Actions/DeviceService.RemoveDevices </sub> | POST | [Update]
| **Device Action (Resource)**                |
| <sub>/api/JobService/Jobs({id}) </sub> | DELETE | [Delete]
| <sub>/api/JobService/Jobs </sub> | POST | [Create]
| <sub>/api/JobService/Jobs({id}) </sub> | GET | [Create]
| <sub>/api/JobService/Jobs({id}) </sub> | GET | [Create]
| **Discovery (Resource)**                |
| <sub>/api/DiscoveryConfigService/DiscoveryConfigGroups({id}) </sub> | GET | [Read, Import, Create, Update]
| <sub>/api/DiscoveryConfigService/DiscoveryConfigGroups </sub> | POST | [Create, Update]
| <sub>/api/DiscoveryConfigService/Actions/DiscoveryConfigService.RemoveDiscoveryGroup </sub> | POST | [Delete]
| <sub>/api/JobService/Jobs({id}) </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs({id})/LastExecutionDetail </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs({id})/ExecutionHistories({id})/ExecutionHistoryDetails </sub> | GET | [Create, Update]
| **Static Group (Resource)**                |
| <sub>/api/GroupService/Groups({id}) </sub> | GET | [Create, Read, Import, Update]
| <sub>/api/GroupService/Groups({id}) </sub> | DELETE | [Delete]
| <sub>/api/GroupService/Actions/GroupService.CreateGroup </sub> | POST | [Create]
| <sub>/api/GroupService/Actions/GroupService.UpdateGroup </sub> | POST | [Update]
| <sub>/api/GroupService/Actions/GroupService.AddMemberDevices </sub> | POST | [Create, Update]
| <sub>/api/GroupService/Actions/GroupService.RemoveMemberDevices </sub> | POST | [Update]
| <sub>/api/GroupService/Groups({id})/Devices </sub> | GET | [Create, Read, Import, Update]
| **User (Resource)**                |
| <sub>/api/AccountService/Accounts </sub> | POST | [Create]
| <sub>/api/AccountService/Accounts({id}) </sub> | PUT | [Update]
| <sub>/api/AccountService/Accounts({id}) </sub> | GET | [Read, Import]
| <sub>/api/AccountService/Accounts({id}) </sub> | DELETE | [Delete]
| **Configuration Compliance (Resource)**                |
| <sub>/api/TemplateService/Baselines({id})/DeviceConfigComplianceReports </sub> | GET | [Read]
| <sub>/api/TemplateService/Baselines </sub> | GET | [Create, Update]
| <sub>/api/TemplateService/Baselines({id}) </sub> | GET | [Create, Update, Read]
| <sub>/api/TemplateService/Actions/TemplateService.Remediate </sub> | POST | [Create, Update]
| <sub>/api/DeviceService/Devices </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs({id}) </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs({id})/LastExecutionDetail </sub> | GET | [Create, Update]
| **Configuration Baseline (Resource)**                |
| <sub>/api/TemplateService/Templates({id}) </sub> | GET | [Create, Update]
| <sub>/api/TemplateService/Templates </sub> | GET | [Create, Update]
| <sub>/api/DeviceService/Devices </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/TemplateService/Baselines </sub> | POST | [Create]
| <sub>/api/TemplateService/Baselines({id}) </sub> | PUT | [Update]
| <sub>/api/TemplateService/Baselines({id}) </sub> | GET | [Create, Read, Update]
| <sub>/api/JobService/Jobs({id}) </sub> | GET | [Create, Update]
| <sub>/api/JobService/Jobs({id})/LastExecutionDetail </sub> | GET | [Create, Update]
| **Deploy (Resource)**                |
| <sub>/api/TemplateService/Templates({id}) </sub> | GET | [Create]
| <sub>/api/TemplateService/Templates </sub> | GET | [Create, Import]
| <sub>/api/TemplateService/Actions/TemplateService.Deploy </sub> | POST | [Create, Update]
| <sub>/api/DeviceService/Devices </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/JobService/Jobs({id}) </sub> | GET | [Create, Update, Delete]
| <sub>/api/JobService/Jobs({id})/LastExecutionDetail </sub> | GET | [Create, Update, Delete]
| <sub>/api/ProfileService/Profiles </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/ProfileService/Actions/ProfileService.UnassignProfiles </sub> | GET | [Update, Delete]
| **Template (Resource)**                |
| <sub>/api/TemplateService/TemplateViewTypes </sub> | GET | [Create]
| <sub>/api/TemplateService/TemplateTypes </sub> | GET | [Create]
| <sub>/api/TemplateService/Templates </sub> | GET | [Create, Import]
| <sub>/api/TemplateService/Templates </sub> | DELETE | [Create, Delete]
| <sub>/api/TemplateService/Templates </sub> | PUT | [Update]
| <sub>/api/TemplateService/Templates </sub> | POST | [Create]
| <sub>/api/TemplateService/Templates({id}) </sub> | GET | [Create, Read, Update]
| <sub>/api/TemplateService/Templates({id})/AttributeDetails </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/TemplateService/Templates({id})/Views(4)/AttributeViewDetails </sub> | GET | [Create, Read, Update, Import]
| <sub>/api/TemplateService/Actions/TemplateService.UpdateNetworkConfig </sub> | POST | [Update]
| <sub>/api/TemplateService/Actions/TemplateService.Clone </sub> | POST | [Create]
| <sub>/api/TemplateService/Actions/TemplateService.Import </sub> | POST | [Create]
| <sub>/api/DeviceService/Devices </sub> | GET | [Create]
| <sub>/api/JobService/Jobs({id}) </sub> | GET | [Create]
| <sub>/api/JobService/Jobs({id})/LastExecutionDetail </sub> | GET | [Create]
| <sub>/api/IdentityPoolService/IdentityPools({id}) </sub> | GET | [Read]
| <sub>/api/IdentityPoolService/IdentityPools </sub> | GET | [Update]
| <sub>/api/NetworkConfigurationService/Networks </sub> | GET | [Update]
| **Network Settings (Resource)**                |
| <sub>/api/ApplicationService/Network/AdapterConfigurations('{interface_name}') </sub> | GET | [Create, Read, Update]
| <sub>/api/ApplicationService/Actions/Network.ConfigureNetworkAdapter | POST | [Create, Update] 
| <sub>/api/JobService/Jobs | GET | [Create, Update]
| <sub>/api/SessionService/SessionConfiguration | GET | [Create, Read, Update]
| <sub>/api/SessionService/Actions/SessionService.SessionConfigurationUpdate | POST | [Create, Update]
| <sub>/api/ApplicationService/Network/TimeConfiguration | GET | [Create, Read]
| <sub>/api/ApplicationService/Network/TimeConfiguration | PUT | [Create, Update]
| <sub>/api/ApplicationService/Network/ProxyConfiguration | GET | [Create, Read]
| <sub>/api/ApplicationService/Network/ProxyConfiguration | PUT | [Create, Update]

