package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var w sync.WaitGroup

/*
 *TODO Edit this constant section!
 */
const defaultProjectName string = "myproject"

// Namespaces
const devNamespace string = "--namespace development"
const rcNamespace string = "--namespace release-candidate"

// Enviroment variables pointing at k8s config files on the local machine
const devKubeConfig string = "$KUBECONFIG_DEVELOPMENT"
const rcKubeConfig string = "$KUBECONFIG_RC"
const prodKubeConfig string = "$KUBECONFIG_PRODUCTION"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	projectName := defaultProjectName
	if args := os.Args; len(args) > 1 {
		projectName = os.Args[1]
	}

	w.Add(3)

	m := make(map[string]string)
	m["dev"] = fmt.Sprintf("kubectl --kubeconfig %s %s", devKubeConfig, devNamespace)
	m["rc"] = fmt.Sprintf("kubectl --kubeconfig %s %s", rcKubeConfig, rcNamespace)
	m["prod"] = fmt.Sprintf("kubectl --kubeconfig %s", prodKubeConfig)

	for env, commandString := range m {
		go fetch(env, commandString, projectName)
	}
	w.Wait()
}

func fetch(env, command, projectName string) {
	defer w.Done()
	cmd := exec.Command("zsh", "-c", command+" get pods -L version")
	stdOutErr, err := cmd.Output()
	if err != nil {
		println(err.Error())
	}
	n := len(stdOutErr)
	output := string(stdOutErr[:n])
	parts := strings.Split(output, "\n")
	for _, podLine := range parts {
		if strings.Contains(podLine, projectName) {

			podLineParts := strings.Fields(podLine)
			version := podLineParts[len(podLineParts)-1]

			status := "❓"
			if strings.Contains(podLine, "Running") {
				if strings.Contains(podLine, "1/1") {
					status = "✅"
				} else {
					status = fmt.Sprintf("⏳ (%s) restarted %s times", podLineParts[2], podLineParts[3])
				}
			}

			fmt.Printf("%s (%s)\t%s\t%s \n", projectName, env, version, status)
		}
	}
}
