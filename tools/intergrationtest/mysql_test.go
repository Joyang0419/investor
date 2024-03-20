package intergrationtest

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMySQLContainer(t *testing.T) {
	pool, resource, dbConn := CreateMySQLContainer("mysqldb")
	defer func() {
		if errPurge := pool.Purge(resource); errPurge != nil {
			log.Fatalf("Could not purge resource: %s", errPurge)
		}
	}()

	assert.Equal(t, "mysql", dbConn.Name())
}
