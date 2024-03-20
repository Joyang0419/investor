package intergrationtest

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMongoDBContainer(t *testing.T) {
	pool, resource, dbConn := CreateMongoDBContainer("mongoDB")
	defer func() {
		if errPurge := pool.Purge(resource); errPurge != nil {
			log.Fatalf("Could not purge resource: %s", errPurge)
		}
	}()

	assert.NoError(t, dbConn.Ping(context.TODO(), nil))
}
