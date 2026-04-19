package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/HackUCF/quincy/common/types"
)

// dumpService creates a json file with the service information.
// it returns the path of the file to be passed into the script.
// Get the path to this file with tempFile.Name().
// the caller needs to call os.Remove(tmpFile.Name()).
func dumpService(svc *types.Service) (*os.File, error) {

	// convert service to json
	bytes, err := json.MarshalIndent(svc, "", "  ")
	if err != nil {
		err = fmt.Errorf("failed to marshal json for service %q: %w", svc.CheckName, err)
		return nil, err
	}

	// create a temporary file
	tmpFile, err := os.CreateTemp("", "quincy-*.json")
	if err != nil {
		err = fmt.Errorf("failed to create temporary file: %w", err)
		return nil, err
	}
	tmpFilePath := tmpFile.Name()

	// always close
	defer tmpFile.Close()

	// write json to temporary file
	_, err = tmpFile.Write(bytes)
	if err != nil {
		err = fmt.Errorf("failed to write temp file %q: %w", tmpFilePath, err)
		return nil, err
	}

	return tmpFile, nil
}

// actually test the server (and capture output)
func runCheck(svc *types.Service, timeout time.Duration) (*scriptOutput, error) {

	// write temporary file
	tmpFile, err := dumpService(svc)
	if err != nil {
		err := fmt.Errorf("failed to dump service to temporary file: %w", err)
		return nil, err
	}
	defer os.Remove(tmpFile.Name()) // remove it on return

	// create output variable
	output := new(scriptOutput)

	// create a timeout context for the command
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// create command to run
	// call the check name with the temporary filepath as the only argument.
	cmd := exec.CommandContext(ctx, string(svc.CheckName), tmpFile.Name())
	cmd.Stdout = &output.Stdout
	cmd.Stderr = &output.Stderr

	// start and wait
	err = cmd.Run()

	// sample error for condition
	var exitError *exec.ExitError

	// check timing out is an organizer failure
	// script logic should determine it's own timeouts and fail gracefully
	select {
	case <-ctx.Done():
		err := fmt.Errorf("check timed out (ran past %0.2f seconds)", timeout.Seconds())
		return nil, err
	default:
	}

	if err == nil {
		// no failure!
		output.Status = Pass
	} else if errors.As(err, &exitError) {
		// exit code other than zero
		output.Status = Fail
	} else {
		// strange error running command
		// could be reading stdout or stderr
		// also could be bad command / args
		// no matter what this is organizer error
		err := fmt.Errorf("failed to run command: %w", err)
		return nil, err
	}

	return output, nil
}
