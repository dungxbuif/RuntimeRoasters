package tests

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"RuntimeRoasters/apps/farm-service/internal/domain"
	"RuntimeRoasters/apps/farm-service/internal/infrastructure/repository"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/database"
	"github.com/stretchr/testify/assert"
)

func setupTestEnforcer() *casbin.SyncedEnforcer {
	modelText := `
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _, _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, "admin") || (g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act)
`
	m, _ := model.NewModelFromString(modelText)
	e, _ := casbin.NewSyncedEnforcer(m)
	return e
}

func TestFarmRepository_Integration(t *testing.T) {
	// 1. Setup real DB connection
	dsn := "postgres://user:password@localhost:54321/farm_db?sslmode=disable"
	db, err := database.NewPostgres(database.PostgresConfig{URL: dsn})
	if err != nil {
		t.Skip("Skipping integration test: database not available")
		return
	}

	// Clean up and recreate table for testing isolation
	db.Migrator().DropTable(&repository.FarmModel{})
	err = db.AutoMigrate(&repository.FarmModel{})
	assert.NoError(t, err)

	enforcer := setupTestEnforcer()
	repo := repository.NewFarmRepository(db, enforcer)

	ctx := context.Background()
	ownerID := "farm-manager-1"
	ctx = identity.InjectContext(ctx, identity.Claims{Subject: ownerID, Role: domain.RoleManager})

	// Setup Policies
	enforcer.AddPolicy(domain.RoleManager, "farm", "read")
	enforcer.AddPolicy(domain.RoleManager, "farm", "write")
	enforcer.AddPolicy(domain.RoleManager, "farm", "delete")
	enforcer.AddGroupingPolicy(ownerID, domain.RoleManager)
	enforcer.AddGroupingPolicy("admin-user", "admin")

	// 2. Test Create
	farm := &domain.Farm{
		Name:       "Cau Dat Specialty",
		Location:   domain.LocationCauDat,
		Area:       50.0,
		CoffeeType: domain.CoffeeTypeArabica,
		OwnerID:    ownerID,
	}
	err = repo.Create(ctx, farm)
	assert.NoError(t, err)
	assert.NotZero(t, farm.ID) // ID should be auto-generated

	createdID := farm.ID

	// 3. Test List (Isolation)
	farms, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, farms, 1)
	assert.Equal(t, "Cau Dat Specialty", farms[0].Name)

	// 4. Test List (Global for Admin)
	adminCtx := identity.InjectContext(context.Background(), identity.Claims{Subject: "admin-user", Role: domain.RoleAdmin})
	allFarms, err := repo.List(adminCtx)
	assert.NoError(t, err)
	assert.Len(t, allFarms, 1)

	// 5. Test Isolation - Another farm manager should see 0
	otherManagerCtx := identity.InjectContext(context.Background(), identity.Claims{Subject: "farm-manager-2", Role: domain.RoleManager})
	enforcer.AddGroupingPolicy("farm-manager-2", domain.RoleManager)
	otherFarms, err := repo.List(otherManagerCtx)
	assert.NoError(t, err)
	assert.Len(t, otherFarms, 0)

	// 6. Test Update
	farm.Area = 60.0
	err = repo.Update(ctx, farm)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, createdID)
	assert.NoError(t, err)
	assert.Equal(t, 60.0, updated.Area)

	// 7. Test Delete
	err = repo.Delete(ctx, createdID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, createdID)
	assert.Error(t, err)
}
