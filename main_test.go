package main

import (
	// "go-fiber-test-job/src/config"
	"go-fiber-test-job/test"
	// testDatabase "go-fiber-test-job/test/database"
	accountTests "go-fiber-test-job/test/tests/account"
	cronTests "go-fiber-test-job/test/tests/cron"
	"testing"
)

// TestMain runs before and after all test cases
func TestMain(m *testing.M) {
	test.InitApp()
	// defer testDatabase.DropDatabase(config.AppConfig.TestDatabase.DbName)

	// Running integration
	m.Run()
}

func TestAllRoutes(t *testing.T) {
	t.Run("TestAccountRoute", accountTests.TestAccountRoute)
	t.Run("TestCronRoute", cronTests.TestCronRoute)
}
