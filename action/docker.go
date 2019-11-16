package action

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func (action Action) reifyArgs(stepInputs StepInputs) []string {
	// TODO: make this logic faster and more correct
	args := action.Runs.Args
	for i := 0; i < len(args); i++ {
		for k, v := range stepInputs {
			args[i] = strings.Replace(args[i], "${{ inputs."+k+" }}", v, -1)
		}
	}
	return args
}

func RunDockerAction(action Action, stepInputs StepInputs, dir string) error {
	if action.Runs.Image != "Dockerfile" {
		return errors.Errorf("Can't use non-Dockerfile image %s", action.Runs.Image)
	}
	imageArgs := action.reifyArgs(stepInputs)
	cmd := exec.Command("docker", "build", "-t", "lights-camera-action-image", dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Building Docker image to run Action...")
	if err := cmd.Run(); err != nil {
		return err
	}
	dockerArgs := append([]string{"run", "lights-camera-action-image"}, imageArgs...)
	cmd = exec.Command("docker", dockerArgs...)
	// TODO: capture "::set-output" et al. in stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Running Docker image...")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
