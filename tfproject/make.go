package tfproject

// We use make to interact with terraform

import (
	"os/exec"

	"github.com/prometheus/common/log"
)

// just because in general that's best practice IMO
// but some systems don't have make and/or terraform
// so we provide the ability to run via docker

type maker interface {
	getCommand(args string) (exec.Cmd, error)
}

var makecommand maker = standardmake{}

// UseDockerForMake call this to have the tool call off to make
// and terraform through a docker command
func UseDockerForMake() {
	makecommand = dockermake{}
}

// const DockerArgsKey = "usedocker"

type standardmake struct{}

func (s standardmake) getCommand(args string) (exec.Cmd, error) {
	bin, err := exec.LookPath("make")
	if err != nil {
		log.Warn("Could not find make in path:", err)
		return exec.Cmd{}, err
	}
	return exec.Cmd{
		Path: bin,
		Args: []string{args},
	}, nil
}

type dockermake struct{}

//TODO make this be configurable?
const defaultDockerArgs = string(`run -it --rm \
   -v $(pwd):/workspace \
   -e AWS_DEFAULT_REGION=$AWS_DEFAULT_REGION \
   -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
   -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
     genesysarch/cloud-workstation make`)

func (s dockermake) getCommand(args string) (exec.Cmd, error) {
	bin, err := exec.LookPath("docker")
	if err != nil {
		log.Warn("Could not find docker in path:", err)
		return exec.Cmd{}, err
	}
	allArgs := []string{defaultDockerArgs}
	allArgs = append(allArgs, args)
	return exec.Cmd{
		Path: bin,
		Args: allArgs,
	}, nil
}
