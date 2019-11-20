package action

import "testing"

func TestDockerAction(t *testing.T) {
	err := Run("actions/hello-world-docker-action@master", map[string]string{
		"who-to-greet": "Mona the Octocat",
	})
	if err != nil {
		t.Error("Error running Dockerfile action:", err)
	}
}
