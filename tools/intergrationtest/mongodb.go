package intergrationtest

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"tools/infra_conn"

	"github.com/ory/dockertest"
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

	host, port := GetHostPort(resource, "27017/tcp")

	var dbConn *mongo.Client
	if err = pool.Retry(func() error {
		var retryErr error
		if dbConn, retryErr = infra_conn.SetupMongoDB(
			infra_conn.MongoDBCfg{
				Host:            host,
				Port:            port,
				Username:        "root",
				Password:        "root",
				Database:        "admin",
				ConnectTimeout:  20 * time.Second,
				MaxPoolSize:     20,
				MaxConnIdleTime: 15 * time.Minute,
			},
		); retryErr != nil {
			return retryErr
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
