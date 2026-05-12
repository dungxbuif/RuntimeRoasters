package tests

import (
	"context"
	"testing"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/infrastructure/repository"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
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

	// Clean up and migrate
	db.Exec("DELETE FROM farms")
	err = db.AutoMigrate(&repository.FarmModel{})
	assert.NoError(t, err)

	enforcer := setupTestEnforcer()
	repo := repository.NewFarmRepository(db, enforcer)

	ctx := context.Background()
	ownerID := "farmer-1"
	ctx = identity.InjectContext(ctx, identity.Claims{Subject: ownerID, Role: "FARMER"})

	// Setup Policies
	enforcer.AddPolicy("FARMER", "farm", "read")
	enforcer.AddPolicy("FARMER", "farm", "write")
	enforcer.AddPolicy("FARMER", "farm", "delete")
	enforcer.AddGroupingPolicy(ownerID, "FARMER")
	enforcer.AddGroupingPolicy("admin-user", "admin")

	// 2. Test Create
	farm := &domain.Farm{
		ID:         "farm-1",
		Name:       "Cau Dat Specialty",
		Location:   "Cầu Đất, Đà Lạt",
		Area:       50.0,
		CoffeeType: "Arabica",
		OwnerID:    ownerID,
	}
	err = repo.Create(ctx, farm)
	assert.NoError(t, err)

	// 3. Test List (Isolation)
	farms, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, farms, 1)
	assert.Equal(t, "Cau Dat Specialty", farms[0].Name)

	// 4. Test List (Global for Admin)
	adminCtx := identity.InjectContext(context.Background(), identity.Claims{Subject: "admin-user", Role: "ADMIN"})
	allFarms, err := repo.List(adminCtx)
	assert.NoError(t, err)
	assert.Len(t, allFarms, 1)

	// 5. Test Isolation - Another farmer should see 0
	otherFarmerCtx := identity.InjectContext(context.Background(), identity.Claims{Subject: "farmer-2", Role: "FARMER"})
	enforcer.AddGroupingPolicy("farmer-2", "FARMER")
	otherFarms, err := repo.List(otherFarmerCtx)
	assert.NoError(t, err)
	assert.Len(t, otherFarms, 0)

	// 6. Test Update
	farm.Area = 60.0
	err = repo.Update(ctx, farm)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, "farm-1")
	assert.NoError(t, err)
	assert.Equal(t, 60.0, updated.Area)

	// 7. Test Delete
	err = repo.Delete(ctx, "farm-1")
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, "farm-1")
	assert.Error(t, err)
}
