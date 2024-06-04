package mongo

import (
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ory/dockertest"

	"tools/infra"
)

func CreateMongoDBContainer(name string) (*dockertest.Pool, *dockertest.Resource, *mongo.Client) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	options := &dockertest.RunOptions{
		Name:       name,
		Repository: "mongo",
		Tag:        "latest",
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=root",
		},
	}
	resource, err := pool.RunWithOptions(options)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)

	}

	host, port := infra.GetHostPort(resource, "27017/tcp")

	var dbConn *mongo.Client
	if err = pool.Retry(func() error {
		if dbConn = SetupConn(
			Config{
				Host:            host,
				Port:            port,
				Username:        "root",
				Password:        "root",
				Database:        "admin",
				ConnectTimeout:  20 * time.Second,
				MaxPoolSize:     20,
				MaxConnIdleTime: 15 * time.Minute,
			},
		); dbConn == nil {
			return errors.New("dbConn is nil")
		}
		return nil

	}); err != nil {
		// You can't defer this because os.Exit doesn't care for defer
		if errPurge := pool.Purge(resource); errPurge != nil {
			log.Fatalf("Could not purge resource: %s", errPurge)
		}
		log.Fatalf("retry failed: %s", err)
	}

	return pool, resource, dbConn
}
