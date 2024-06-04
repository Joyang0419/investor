package mysql

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
	"gorm.io/gorm"

	"tools/infra"
)

func CreateMySQLContainer(name string) (*dockertest.Pool, *dockertest.Resource, *gorm.DB) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	options := &dockertest.RunOptions{
		Name:       name,
		Repository: "mysql",
		Tag:        "latest",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=root",
			"MYSQL_DATABASE=dev",
			"MYSQL_USER=joy",
			"MYSQL_PASSWORD=joy",
		},
	}
	resource, err := pool.RunWithOptions(options)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)

	}

	host, port := infra.GetHostPort(resource, "3306/tcp")
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = mysql.SetLogger(log.New(io.Discard, "", log.LstdFlags)); err != nil {
		log.Fatalf("set logger failed: %v", err)
	}

	if err != nil {
		log.Fatalf("strconv.ParseInt failed: %v", err)
	}

	var dbConn *gorm.DB
	if err = pool.Retry(func() error {
		if dbConn = SetupConn(
			Config{
				Host:            host,
				Port:            port,
				Username:        "joy",
				Password:        "joy",
				Database:        "dev",
				MaxIdleConns:    10,
				MaxOpenConns:    10,
				ConnMaxLifeTime: 60 * time.Second,
			}, nil,
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