package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func ExecBashCommand(command string) (string, error) {
	log.Printf("[DEBUG] executing bash command: %s", command)

	cmd := exec.Command("bash", "-x", "-c", command)

	stderr, _ := cmd.StderrPipe()

	output, err := cmd.Output()

	if err != nil {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Println("[DEBUG] " + scanner.Text())
		}

		return string(output), err
	}

	if output == nil {
		return "", nil
	}

	return string(output), nil
}

func ExecuteScript(scriptPath string, multiClusterCertsFolder string, skipReadyCheck bool, kubeconfig string) (string, error) {
	log.Printf("[DEBUG] executing script")

	command := "./" + scriptPath
	if multiClusterCertsFolder != "" {
		command = command + " -c " + multiClusterCertsFolder
	}

	if skipReadyCheck {
		command = command + " --skip-ready-check"
	}

	output, err := ExecBashCommand(fmt.Sprintf("KUBECONFIG=%s %s", kubeconfig, command))
	if err != nil {
		return output, err
	}

	return output, nil
}

func MakeExecutable(scriptPath string) error {
	log.Printf("[DEBUG] making executable: " + scriptPath)

	err := os.Chmod(scriptPath, 0700)
	if err != nil {
		return err
	}

	return nil
}
