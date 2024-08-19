package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const INSTANCE_NAME string = "Portkey"
const INSTANCE_TITLE string = "Portkey boot"

const (
	HOST_INDEX   = 0
	TYPE_INDEX   = 1
	OBJ_INDEX    = 2
	RESPAWN_WAIT = 5
)

const baseCommand string = "ssh"

var baseArgs = []string{
	"-N",
	"-o",
	"ExitOnForwardFailure=yes",
	"-o",
	"StreamLocalBindUnlink=yes",
	"-o",
	"ServerAliveInterval=5",
	"-o",
	"ServerAliveCountMax=1",
	"-L",
}

func setup_exec(target string, metadata []string) (*exec.Cmd, error) {

	var command *exec.Cmd
	var commandErr error
	var commandLine []string = nil

	switch metadata[TYPE_INDEX] {

	// tcp socket
	case "tcp":
		commandLine = append(baseArgs, fmt.Sprintf("%s:localhost:%s",
			metadata[OBJ_INDEX], metadata[OBJ_INDEX]))
	// unix socket
	case "unix":
		commandLine = append(baseArgs, fmt.Sprintf("%s:%s",
			metadata[OBJ_INDEX], metadata[OBJ_INDEX]))

	default:
		fmt.Println(metadata[TYPE_INDEX])
		return nil, commandErr
	}

	commandLine = append(commandLine, metadata[HOST_INDEX])

	command = exec.Command(baseCommand, commandLine...)
	commandErr = command.Start()

	if commandErr != nil {
		return nil, commandErr
	}
	return command, commandErr
}

func nights_watch(target string, metadata []string) {

	var command *exec.Cmd
	var commandErr error

	for true {

		// exec and start the process
		// sleep for 5 sec
		command, commandErr = setup_exec(target, metadata)

		if command == nil || commandErr != nil {
			fmt.Printf("Oops cannot spawn process \n")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		command.Wait()
		time.Sleep(time.Duration(5) * time.Second)

	}

}

func main() {

	var configBox map[string]interface{}
	// converting map of interface into string slices
	var portConfig = map[string][]string{}
	var start string
	var first bool = false

	portFile, err := os.ReadFile("./services.json")

	if err != nil {
		fmt.Printf("There is an error while reading the file \n")
		os.Exit(1)
	}

	err = json.Unmarshal(portFile, &configBox)
	if err != nil {
		fmt.Printf("Error in unmarshalling! \n")
		os.Exit(2)
	}

	for i := range configBox {

		switch interfaceType := configBox[i].(type) {

		case []interface{}:
			for j := 0; j < len(interfaceType); j++ {

				switch valueType := interfaceType[j].(type) {
				case string:
					portConfig[i] = append(portConfig[i], valueType)
				default:
					continue
				}
			}

		default:
			continue
		}
	}

	for i := range portConfig {
		if !first {
			first = true
			start = i
			continue
		}
		go nights_watch(i, portConfig[i])
	}
	nights_watch(start, portConfig[start])

}
