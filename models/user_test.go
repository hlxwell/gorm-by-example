package models

import (
	"encoding/json"
	"testing"

	"github.com/hlxwell/gorm-by-example/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Should be able to get User with Roles
func TestGetUserWithRoles(t *testing.T) {
	prepareTestDB()

	var user User
	// Load specific relationship
	db.Conn.Preload("Roles").First(&user)
	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.Roles)

	// Load all relationships
	db.Conn.Preload(clause.Associations).First(&user)
	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.Roles)
}

// Should be able to Check User.Roles and Features with one query
func TestLoadUserWithNestedRelation(t *testing.T) {
	prepareTestDB()

	var user User
	db.Conn.Preload("Roles.Features").First(&user)
	assert.NotEmpty(t, user.Roles[0].Features)
}

// Should be able to Add Exists Role to User.Roles and update the RoleCount
func TestUpdateRoleCountbyAddRole(t *testing.T) {
	prepareTestDB()

	var user User
	db.Conn.Preload("Roles").First(&user)
	roleCount := len(user.Roles)

	// Create Role for user
	var role Role
	db.Conn.Last(&role)
	user.AddRole(role)

	// Reload user and role.
	db.Conn.Preload("Roles").First(&user)
	assert.EqualValues(t, user.RoleCount, roleCount+1)
}

// Should be able to Uniqueness Validation for User.Roles
func TestUniqueRoleValidation(t *testing.T) {
	prepareTestDB()

	var user User
	db.Conn.Preload("Roles").First(&user)
	err := user.AddRole(*user.Roles[0])
	assert.Error(t, err)
}

// Should be able to Presence Validation for User.Name before save
func TestUserNamePresenceValidation(t *testing.T) {
	prepareTestDB()

	var user User
	db.Conn.First(&user)
	result := db.Conn.Create(&User{
		Name: user.Name,
	})
	assert.Error(t, result.Error)
}

// It will use transaction also.
// Should be able to Create User with roles
func TestCreateUserWithRole(t *testing.T) {
	prepareTestDB()

	db.SetScope(1234)
	result := db.ScopedConn().Create(&User{
		Name: "hlxwell",
		Roles: []*Role{
			{Name: "admin"},
			{Name: "writer"},
			{Name: "reader"},
		},
	})

	assert.NoError(t, result.Error)

	var user User
	db.Conn.Preload(clause.Associations).Last(&user)
	assert.Equal(t, 3, len(user.Roles))
}

func TestUpdateUser(t *testing.T) {
	prepareTestDB()

	result := db.Conn.Create(&User{
		Name: "old-hlxwell",
		Roles: []*Role{
			{Name: "admin"},
			{Name: "writer"},
			{Name: "reader"},
		},
	})
	assert.NoError(t, result.Error)

	var user User
	db.Conn.Last(&user)

	user.Name = "new-hlxwell"
	db.Conn.Save(&user)
}

// Should be able to Create Role and User in one request from Nested Form
func TestCreateUserAndRoleInOneRequest(t *testing.T) {
	prepareTestDB()

	err := db.Conn.Transaction(func(tx *gorm.DB) error {
		role := &Role{Name: "god"}
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		if err := tx.Create(&User{Name: "Michael He", Roles: []*Role{role}}).Error; err != nil {
			return err
		}
		return nil
	})

	var role Role
	db.Conn.Where("name = ?", "god").First(&role)
	assert.EqualValues(t, 0, role.ID) // Should not create role
	assert.Error(t, err)
}

// Should be able to find normal role users
func TestFindUserWithNormalRole(t *testing.T) {
	prepareTestDB()

	var users []User
	db.Conn.
		Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.name = ?", "normal").Find(&users)
	assert.NotEmpty(t, users)
}

// Should be able to use Scope to filter users with specific Role
func TestScopeFilterForUser(t *testing.T) {
	prepareTestDB()
	var users []User
	db.Conn.Scopes(WithRole("normal")).Find(&users)
	assert.NotEmpty(t, users)

	db.Conn.Scopes(WithRole("super")).Find(&users)
	assert.Empty(t, users)
}

// Should be able to batch updating user.Roles
func TestBatchUpdatingUserRole(t *testing.T) {
	prepareTestDB()

	var user User
	var role Role

	db.Conn.Preload("Roles").First(&user)
	assert.NotEmpty(t, user.Roles)

	// Add new role & Delete old roles
	// db.Conn.Delete(&user.Roles)
	// equals to
	// db.Conn.Model(&user).Association("Roles").Delete(&user.Roles)
	db.Conn.Last(&role)
	db.Conn.Model(&user).Association("Roles").Replace([]*Role{&role})

	db.Conn.Preload("Roles").First(&user)
	assert.EqualValues(t, 1, len(user.Roles))
}

// With user has JSON attributes
// Should be able to find by attributes.age
func TestJSONAttributes(t *testing.T) {
	prepareTestDB()

	// Create attributes
	var attrs UserAttributes
	attrs.Name = "Michael-He"
	attrs.Age = 18
	attrs.Tags = []string{"niu", "bi"}
	attrsJSON, _ := json.Marshal(attrs)

	// Create user
	var user User
	user.Name = "test user"
	user.Attributes = attrsJSON
	db.Conn.Create(&user)

	// Fetch user
	var newUser1 User
	assert.NoError(t, db.Conn.First(&newUser1, datatypes.JSONQuery("attributes").Equals(18, "age")).Error)
	assert.NotZero(t, newUser1.ID)

	var newUser2 User
	assert.Error(t, db.Conn.First(&newUser2, datatypes.JSONQuery("attributes").Equals(19, "age")).Error)
	assert.Zero(t, newUser2.ID)

	// Should be able to find by attributes.tags
	var newUser3 User
	assert.NoError(t, db.Conn.Where("JSON_CONTAINS(`attributes`, ?, '$.tags')", `["niu"]`).First(&newUser3).Error)
	assert.NotZero(t, newUser3.ID)
}

// Should be able to Load Users
func TestLoadUsers(t *testing.T) {
	prepareTestDB()

	user := User{}
	result := db.Conn.First(&user)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Michael He", user.Name)
}

// Should be able to Create User
func TestCreateUser(t *testing.T) {
	prepareTestDB()

	result := db.Conn.Create(&User{
		Name: "Michael",
	})
	assert.NoError(t, result.Error)
}

func TestDeleteRoleFromUser(t *testing.T) {
	prepareTestDB()

	var user User
	db.Conn.Preload("Roles").First(&user)
	user.DeleteRole(*user.Roles[0])

	// Reload
	var newUser User
	db.Conn.Preload("Roles").First(&newUser)
	assert.Equal(t, uint(len(newUser.Roles)), newUser.RoleCount)
	assert.Equal(t, uint(0), newUser.RoleCount)
}

// Should be able to Add New Role to User.Roles and update the RoleCount
func TestAddNewRole(t *testing.T) {
	prepareTestDB()

	var user User
	db.Conn.Preload("Roles").First(&user)

	// Create Role for user
	var role Role
	db.Conn.First(&role, 2)
	err := user.AddRole(role)
	assert.Nil(t, err)

	// Reload user and role.
	var r Role
	db.Conn.Preload("Users").Last(&r)
	assert.Equal(t, 1, len(r.Users))
	assert.Equal(t, 2, int(r.Users[0].RoleCount))
}
