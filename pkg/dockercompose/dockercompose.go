package dockercompose

import (
	"log"

	"golang.org/x/net/context"
	_ "github.com/lib/pq"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
	"github.com/docker/libcompose/lookup"
	"os"
	"github.com/phayes/freeport"
	"strconv"
)

var runningProject project.APIProject
var runningContext context.Context

func StartDockerCompose(composeFile string, projectName string, portsToFind []string) {
	for _, key := range portsToFind {
		freePort, err := freeport.GetFreePort()
		if err != nil {
			panic(err)
		}
		os.Setenv(key, strconv.Itoa(freePort))
	}

	dockerCompose, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles:      []string{composeFile},
			ProjectName:       projectName,
			EnvironmentLookup: &lookup.OsEnvLookup{},
		},
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	runningProject = dockerCompose

	runningContext = context.Background()
	err = dockerCompose.Up(runningContext, options.Up{})

	if err != nil {
		log.Fatal(err)
	}
}

func StopDockerCompose() {
	err := runningProject.Stop(runningContext, 10)
	if err != nil {
		panic(err)
	}
}
