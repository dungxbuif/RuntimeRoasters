package integration

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
 * TC-1.4, TC-1.5, TC-1.6: Backend Data Consistency Audit
 * Chạy sau khi Seeding thành công để đảm bảo DB khớp với Spec.
 */

type TestConfig struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
}

func getTestConfig() TestConfig {
	return TestConfig{
		DBUser: "user",
		DBPass: "password",
		DBHost: "localhost",
		DBPort: "54321",
	}
}

func TestFlow1_BootstrapDataConsistency(t *testing.T) {
	cfg := getTestConfig()

	t.Run("TC-1.5: Farm Service Seed Check", func(t *testing.T) {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/farm_db?sslmode=disable", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort)
		db, err := sql.Open("postgres", dsn)
		require.NoError(t, err)
		defer db.Close()

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM farms").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 6, count, "Should have 6 seeded farms")
	})

	t.Run("TC-1.4: Retail Service Seed Check", func(t *testing.T) {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/retail_db?sslmode=disable", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort)
		db, err := sql.Open("postgres", dsn)
		require.NoError(t, err)
		defer db.Close()

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM stores").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 5, count, "Should have 5 seeded stores")
	})

	t.Run("TC-1.6: Logistics Service Seed Check", func(t *testing.T) {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/logistics_db?sslmode=disable", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort)
		db, err := sql.Open("postgres", dsn)
		require.NoError(t, err)
		defer db.Close()

		// Check Locations
		var locCount int
		err = db.QueryRow("SELECT COUNT(*) FROM locations").Scan(&locCount)
		require.NoError(t, err)
		assert.Equal(t, 14, locCount, "Should have 14 seeded locations (6 farms + 3 WH + 5 retail)")

		// Check Fleet
		var driverCount int
		err = db.QueryRow("SELECT COUNT(*) FROM drivers").Scan(&driverCount)
		require.NoError(t, err)
		assert.Equal(t, 3, driverCount, "Should have 3 seeded drivers")

		var vehicleCount int
		err = db.QueryRow("SELECT COUNT(*) FROM vehicles").Scan(&vehicleCount)
		require.NoError(t, err)
		assert.Equal(t, 3, vehicleCount, "Should have 3 seeded vehicles")
	})
}
