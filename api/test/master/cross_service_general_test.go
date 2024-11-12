package test

import (
	"after-sales/api/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Define base URL for the role endpoint
const baseURL = "https://testing-backendims.indomobil.co.id/general-service/v1/role"

// Role represents the structure for role data
type Role struct {
	RoleCode string `json:"role_code"`
	RoleName string `json:"role_name"`
}

// RoleResponse represents the structure for a single role response
type RoleResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       Role   `json:"data"`
}

// DeleteResponse represents the structure for a delete response
type DeleteResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// TestGetAllRoles tests retrieving all roles
func TestGetAllRoles(t *testing.T) {
	var result []Role
	err := utils.Get(baseURL, &result, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// TestCreateRole tests creating a new role
func TestCreateRole(t *testing.T) {
	role := Role{
		RoleCode: "test_code",
		RoleName: "test_name",
	}

	// Step 1: Create the Role
	err := utils.Post(baseURL, role, nil)
	assert.NoError(t, err, "Failed to create role")

	// Step 2: Fetch the Created Role to Confirm
	var roles []Role
	err = utils.Get(baseURL, &roles, nil)
	assert.NoError(t, err, "Failed to retrieve roles")

	// Step 3: Check if the created role exists in the retrieved roles
	var createdRole *Role
	for _, r := range roles {
		if r.RoleCode == "test_code" && r.RoleName == "test_name" {
			createdRole = &r
			break
		}
	}
	assert.NotNil(t, createdRole, "Created role not found in role list")
	assert.Equal(t, "test_code", createdRole.RoleCode)
	assert.Equal(t, "test_name", createdRole.RoleName)
}

// TestUpdateRole tests updating an existing role by role_id
func TestUpdateRole(t *testing.T) {
	roleID := "6" // Replace with a valid role_id for testing
	updateURL := fmt.Sprintf("%s/%s", baseURL, roleID)

	role := Role{
		RoleCode: "updated_code",
		RoleName: "updated_name",
	}
	var result RoleResponse
	err := utils.Put(updateURL, role, &result)
	assert.NoError(t, err)
	assert.Equal(t, "updated_code", result.Data.RoleCode)
	assert.Equal(t, "updated_name", result.Data.RoleName)
}

// TestDeleteRole tests deleting a role by role_id
func TestDeleteRole(t *testing.T) {
	roleID := "6" // Replace with a valid role_id for testing
	deleteURL := fmt.Sprintf("%s/%s", baseURL, roleID)

	var result DeleteResponse
	err := utils.Delete(deleteURL, nil, &result)
	assert.NoError(t, err)
	assert.Equal(t, 200, result.StatusCode)
}
