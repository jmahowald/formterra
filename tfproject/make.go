package tfproject

// We use make to interact with terraform

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/prometheus/common/log"
)

// just because in general that's best practice IMO
// but some systems don't have make and/or terraform
// so we provide the ability to run via docker

//exec.Cmd simple wrapper for exec.exec.Cmd so that we could add methods to it

type maker interface {
	getMakeCommand(TerraformLayer) (exec.Cmd, error)
}

var makecommand maker = standardmake{}

type makeable interface {
	GetMake() (exec.Cmd, error)
}

//PlanCommand executable command to see plan for layer
func (t TerraformLayer) PlanCommand() (exec.Cmd, error) {
	makeCmd, err := t.GetMake()
	return addArgs(makeCmd, "plan"), err
}

//ApplyCommand executable command to build the layer
func (t TerraformLayer) ApplyCommand() (exec.Cmd, error) {
	makeCmd, err := t.GetMake()
	return addArgs(makeCmd, "apply"), err
}

//GetMake we interact with terraform through make
func (t TerraformLayer) GetMake() (exec.Cmd, error) {
	//See if makefile exists in the layer's directory
	dir, _ := t.dir()
	_, err := os.Stat(filepath.Join(dir, "Makefile"))

	if err != nil {
		log.Warn("No makefile found in %s", t)
		return exec.Cmd{}, err
	}
	return makecommand.getMakeCommand(t)
}

func addArgs(cmd exec.Cmd, args ...string) exec.Cmd {
	newargs := make([]string, len(cmd.Args)+len(args))
	newargs = append(newargs, cmd.Args...)
	newargs = append(newargs, args...)
	cmd.Args = newargs
	return cmd
}

type standardmake struct{}

func (s standardmake) getMakeCommand(layer TerraformLayer) (exec.Cmd, error) {
	bin, err := exec.LookPath("make")
	if err != nil {
		log.Warn("Could not find make in path:", err)
		return exec.Cmd{}, err
	}
	dir, _ := layer.dir()
	return exec.Cmd{Path: bin, Dir: dir}, nil
}

// const DockerArgsKey = "usedocker"

// UseDockerForMake call this to have the tool call off to make
// and terraform through a docker command
func UseDockerForMake() {
	makecommand = dockermake{}
}

type dockermake struct{}

//TODO make this be configurable?
const defaultDockerArgs = string(`run -it --rm \
   -v $(pwd):/workspace \
   -e AWS_DEFAULT_REGION=$AWS_DEFAULT_REGION \
   -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
   -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
     genesysarch/cloud-workstation make`)

func (s dockermake) getMakeCommand(layer TerraformLayer) (exec.Cmd, error) {
	bin, err := exec.LookPath("docker")
	if err != nil {
		log.Warn("Could not find docker in path:", err)
		return exec.Cmd{}, err
	}
	allArgs := []string{defaultDockerArgs}
	dir, _ := layer.dir()

	return exec.Cmd{
		Path: bin,
		Args: allArgs,
		Dir:  dir,
	}, nil
}
