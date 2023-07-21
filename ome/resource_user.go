package ome

import (
	"context"
	"reflect"
	"strings"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &userResource{}
)

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	p *omeProvider
}

// Configure implements resource.ResourceWithConfigure
func (r *userResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing user on OpenManage Enterprise.",
		Version:             1,
		Attributes:          UserSchema(),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_user create : Started")
	//Get Plan Data
	var plan, state models.OmeUser
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	omeClient, d := r.p.createOMESession(ctx, "resource_configuration_User Create")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	up := getUserPayload(ctx, &plan) 
	
	tflog.Trace(ctx, "resource_user create Creating User")
	tflog.Debug(ctx, "resource_user create Creating User", map[string]interface{}{
		"Create User Request": up,
	})

	cUser, err := omeClient.CreateUser(up)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateUser, err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_configuration_User : create Finished creating User")
	tflog.Trace(ctx, "resource_user create: updating state finished, saving ...")
	// Save into State
	state = saveState(cUser)
	state.Password = plan.Password
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_user create: finish")
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "resource_user read: started")
	var state models.OmeUser
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_configuration_User Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	user, err := omeClient.GetUserByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrReadUser, err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_user read: finished reading state")
	//Save into State
	istate := saveState(user)
	istate.Password = state.Password
	diags = resp.State.Set(ctx, &istate)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_user read: finished")
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get state Data
	tflog.Trace(ctx, "resource_user update: started")
	var state, plan models.OmeUser
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get plan Data
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and defer the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_configuration_baseline Update")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	if !reflect.DeepEqual(state, plan) {
		updatePayload := models.User{
			ID:                 state.ID.ValueString(),
			UserTypeID:         int(plan.UserTypeID.ValueInt64()),
			DirectoryServiceID: int(plan.DirectoryServiceID.ValueInt64()),
			Description:        plan.Description.ValueString(),
			Password:           plan.Password.ValueString(),
			UserName:           plan.UserName.ValueString(),
			RoleID:             plan.RoleID.ValueString(),
			Locked:             plan.Locked.ValueBool(),
			Enabled:            plan.Enabled.ValueBool(),
		}
		user, err := omeClient.UpdateUser(updatePayload)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrGnrUpdateUser, err.Error(),
			)
			return
		}
		state = saveState(user)
		state.Password = plan.Password
		tflog.Trace(ctx, "resource_configuration_baseline : update Finished creating Baseline")
	}
	tflog.Trace(ctx, "resource_user update: finished state update")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_user update: finished")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_user delete: started")
	// Get State Data
	var state models.OmeUser
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Session and differ the remove session
	omeClient, d := r.p.createOMESession(ctx, "resource_configuration_baseline Delete")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	status, err := omeClient.DeleteUser(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrDeleteUser,
			err.Error(),
		)
	}

	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "resource_user delete: finished "+status)
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	parser := req.ID
	items := strings.SplitN(parser, ",", 2)
	if len(items) < 2 {
		resp.Diagnostics.AddError(
			"Error while user import",
			"Error while user import",
		)
		return
	}
	id := items[0]
	password := items[1]
	idAttrPath := path.Root("id")
	passwordAttrPath := path.Root("password")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, idAttrPath, id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, passwordAttrPath, password)...)
}

func getUserPayload(ctx context.Context, plan *models.OmeUser) (models.UserPayload) {
	user := models.UserPayload{
		UserTypeID:         int(plan.UserTypeID.ValueInt64()),
		DirectoryServiceID: int(plan.DirectoryServiceID.ValueInt64()),
		Description:        plan.Description.ValueString(),
		Password:           plan.Password.ValueString(),
		UserName:           plan.UserName.ValueString(),
		RoleID:             plan.RoleID.ValueString(),
		Locked:             plan.Locked.ValueBool(),
		Enabled:            plan.Enabled.ValueBool(),
	}
	return user
}

func saveState(resp models.User) (state models.OmeUser) {
	state.Description = types.StringValue(resp.Description)
	state.ID = types.StringValue(resp.ID)
	state.UserName = types.StringValue(resp.UserName)
	state.RoleID = types.StringValue(resp.RoleID)
	state.Enabled = types.BoolValue(resp.Enabled)
	state.Locked = types.BoolValue(resp.Locked)
	state.UserTypeID = types.Int64Value(int64(resp.UserTypeID))
	state.DirectoryServiceID = types.Int64Value(int64(resp.DirectoryServiceID))
	return
}
