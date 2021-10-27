package utils

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

const kubectlApplyCommand = "kubectl apply -f "

func ExecBashCommand(command string) (string, error) {
	log.Printf("[DEBUG] executing bash command")

	cmd := exec.Command("bash", "-c", command)

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

func ExecuteScript(scriptPath string, multiClusterCertsFolder string) (string, error) {
	log.Printf("[DEBUG] executing script")

	//if runtime.GOOS == "windows" {
	//	return applyYaml(yamlPath)
	//}

	err := MakeExecutable(scriptPath)
	if err != nil {
		return "", err
	}

	command := "./" + scriptPath
	if multiClusterCertsFolder != "" {
		command = command + " -c " + multiClusterCertsFolder
	}

	_, err = ExecBashCommand(command)

	if err != nil {
		return "", err
	}

	return "", nil
}

func applyYaml(yamlPath string) (string, error) {
	log.Printf("[DEBUG] applying yaml: " + yamlPath)
	return ExecBashCommand(kubectlApplyCommand + yamlPath)
}

func MakeExecutable(scriptPath string) error {
	log.Printf("[DEBUG] making executable: " + scriptPath)

	err := os.Chmod(scriptPath, 0700)
	if err != nil {
		return err
	}

	return nil
}
