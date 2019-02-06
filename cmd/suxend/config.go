package main

import (
	"flag"
	"fmt"
	"os"
)

var version string

type configuration struct {
	version     bool
	environment string
	host        string
	port        int
	static      string
	gcloud      struct {
		project string
	}
	logger struct {
		environment string
		level       string
	}
	nexus struct {
		address    string
		repository string
		searchPath string
		svc        struct {
			address   string
			authToken string
		}
		registry struct {
			address string
		}
	}
}

func (c *configuration) init() {
	if c == nil {
		*c = configuration{}
	}

	flag.BoolVar(&c.version, "version", false, "print version and exit")
	flag.StringVar(&c.environment, "env", "", "environment name")
	flag.StringVar(&c.host, "host", "127.0.0.1", "host")
	flag.StringVar(&c.static, "static", "./ui/dist", "path to static files")
	flag.IntVar(&c.port, "port", 8080, "port")
	// LOGGER
	flag.StringVar(&c.logger.environment, "log.environment", "production", "logger environment config")
	flag.StringVar(&c.logger.level, "log.level", "info", "logger level")
	// GOOGLE CLOUD
	flag.StringVar(&c.gcloud.project, "gcloud.project", "", "google cloud project id")
	// NEXUS REST API
	flag.StringVar(&c.nexus.address, "nexus.address", "nexus:8080", "nexus api address")
	flag.StringVar(&c.nexus.svc.address, "nexus.svc.address", "nexus:8080", "k8s endpoint of the nexus svc")
	flag.StringVar(&c.nexus.svc.authToken, "nexus.svc.authToken", "", "if required, add basic auth token")
	flag.StringVar(&c.nexus.repository, "nexus.repository", "docker-hosted", "nexus repository to search within")
	flag.StringVar(&c.nexus.searchPath, "nexus.searchPath", "service/rest/v1/search", "uri for searching repositories")
	flag.StringVar(&c.nexus.registry.address, "nexus.registry.address", "containers.example.com", "hostname from docker can pull images")
}

func (c *configuration) parse() {
	if !flag.Parsed() {
		flag.Parse()
		if c.version {
			fmt.Printf("%s", version)
			os.Exit(0)
		}
	}
}
