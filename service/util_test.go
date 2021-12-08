package service

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetContainer(t *testing.T) {
	_, err := GetContainer("should_not_exists")

	if err == nil {
		t.Error("Expected not found error, got nil")
	}

	containerName := "testcontainer_" + strconv.Itoa(int(time.Now().UnixNano()))

	id, err := exec.Command("docker", "run", "-td", "--name", containerName, "alpine:latest").CombinedOutput()
	if err != nil {
		t.Errorf("Error creating container: %s", err)
	}

	container, err := GetContainer(containerName)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if container.ID != strings.TrimSpace(string(id)) {
		t.Errorf("Expected container id %s, got %s", string(id), container.ID)
	}

	out, err := exec.Command("docker", "run", "-td", "--name", containerName+"2nd", "alpine:latest").CombinedOutput()
	if err != nil {
		t.Errorf("Error creating second container: %s.\n%s", err, string(out))
	}

	_, err = GetContainer(containerName)

	if err == nil {
		t.Errorf("Expected ambigous container names error, got nil")
	}

	// clean up
	go exec.Command("docker", "rm", "-f", containerName).Run()
	go exec.Command("docker", "rm", "-f", containerName+"2nd").Run()
}

func TestRenameContainer(t *testing.T) {
	currentName := nextContainerName()
	newName := nextContainerName()

	id, err := exec.Command("docker", "run", "-td", "--name", currentName, "alpine:latest").CombinedOutput()
	if err != nil {
		t.Errorf("Error creating container: %s", err)
	}

	err = RenameContainer(strings.TrimSpace(string(id)), newName)

	if err != nil {
		t.Errorf("Error renaming container: %s", err)
	}

	_, err = GetContainer(currentName)
	if err == nil {
		t.Errorf("Expected not found error when searching for the container with its old name, got nil")
	}

	_, err = GetContainer(newName)
	if err != nil {
		t.Errorf("Expected no error when searching for the container with its new name, got %s", err)
	}

	// clean up
	go exec.Command("docker", "rm", "-f", newName).Run()
}

func TestStopContainerByName(t *testing.T) {
	name := nextContainerName()
	_, err := exec.Command("docker", "run", "-td", "--name", name, "alpine:latest").CombinedOutput()
	if err != nil {
		t.Errorf("Error creating container: %s", err)
	}

	err = StopContainerByName(name)

	if err != nil {
		t.Errorf("Error stopping container: %s", err)
	}

	_, err = GetContainer(name)
	if err == nil {
		t.Errorf("Expected not found error when searching for a deleted container, got nil")
	}
}

func nextContainerName() string {
	return "testcontainer_" + strconv.Itoa(int(time.Now().UnixNano()))
}
