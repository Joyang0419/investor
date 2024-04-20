package intergrationtest

import (
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/ory/dockertest"
)

func GetHostPort(resource *dockertest.Resource, id string) (host string, port int) {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		portInt, err := strconv.ParseInt(resource.GetPort(id), 10, 64)
		if err != nil {
			panic(err)
		}

		return "localhost", int(portInt)
	}

	u, err := url.Parse(dockerURL)
	if err != nil {
		log.Fatal(err)
	}

	portInt, err := strconv.ParseInt(u.Port(), 10, 64)
	if err != nil {
		panic(err)
	}

	return u.Hostname(), int(portInt)
}
