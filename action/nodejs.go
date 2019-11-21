package action

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func nodeEnvVars(stepInputs StepInputs) []string {
	vars := []string{}
	for k, v := range stepInputs {
		bind := "INPUT_" + strings.ToUpper(strings.Replace(k, " ", "_", -1)) + "=" + v
		vars = append(vars, bind)
	}
	return vars
}

func RunNodeAction(action Action, stepInputs StepInputs, dir string) error {
	cmd := exec.Command("node", action.Runs.Main)
	// TODO: capture "::set-output" et al. in stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), nodeEnvVars(stepInputs)...)
	fmt.Println("env:", cmd.Env)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
