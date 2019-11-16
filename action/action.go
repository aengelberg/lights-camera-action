package action

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

type StepInputs map[string]string

type Param struct {
	Description string
	Required    bool
	Default     string
}

type Action struct {
	Name        string
	Description string
	Inputs      map[string]Param
	Outputs     map[string]Param
	Runs        struct {
		Using string
		Main  string
		Image string
		Args  []string
	}
}

func tempDir() string {
	i, err := rand.Int(rand.Reader, big.NewInt(100000000))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("/tmp/action%d", i)
}

func cloneRepo(actionRef string, dir string) error {
	strs := strings.Split(actionRef, "@")
	if len(strs) != 2 {
		return fmt.Errorf("Expected org/repo@tag, got %s", actionRef)
	}
	repoName := strs[0]
	tag := strs[1]
	cmd := exec.Command("git", "clone", "https://github.com/"+repoName, dir, "--branch", tag, "--depth", "1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "error cloning the action repo: ")
	}
	return nil
}

func parseAction(dir string) (action Action, err error) {
	in, err := os.Open(dir + "/action.yml")
	if err != nil {
		return
	}
	defer in.Close()
	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		return
	}
	yaml.Unmarshal(bytes, &action)
	return
}

func Run(actionRef string, stepInputs StepInputs) error {
	fmt.Println("Running", actionRef, stepInputs)
	dir := tempDir()
	defer func() {
		fmt.Println("Cleaning up", dir)
		os.RemoveAll(dir)
	}()
	if err := cloneRepo(actionRef, dir); err != nil {
		return err
	}
	action, err := parseAction(dir)
	if err != nil {
		return err
	}
	// TODO: resolve required/optional parameters, add defaults for missing fields
	switch action.Runs.Using {
	case "docker":
		RunDockerAction(action, stepInputs, dir)
	}
	return nil
}
