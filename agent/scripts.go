package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/HackUCF/Quincy/common/types"
	"github.com/google/uuid"
)

type scriptPath string

var (
	scriptCache      = make(map[types.CheckName]scriptPath)
	scriptCacheMutex sync.RWMutex
)

func (cfg *agentConfig) getScript(id types.CheckName) (scriptPath, error) {
	scriptCacheMutex.RLock()
	cachedPath, ok := scriptCache[id]
	scriptCacheMutex.RUnlock()

	if ok {
		// early return. the script is cached
		return cachedPath, nil
	}

	// list all files in the check directory
	scripts, err := os.ReadDir(cfg.checksDir)
	if err != nil {
		err = fmt.Errorf("couldn't list directory '%s': %s", cfg.checksDir, err.Error())
		return "", err
	}

	// loop through them all
	for _, script := range scripts {

		// match 'SSH', 'ssh.py', 'ssh.sh', 'ssh.x86-64.exe', etc
		// but not 'ssh-windows.sh' or 'ssh2.py'
		scriptID := strings.Split(script.Name(), ".")[0]
		if strings.EqualFold(scriptID, string(id)) {

			// script is found
			// construct full path
			joinedPath := filepath.Join(cfg.checksDir, script.Name())
			{
				// add to cache
				scriptCacheMutex.Lock()
				scriptCache[id] = scriptPath(joinedPath)
				scriptCacheMutex.Unlock()
			}
			return scriptCache[id], nil
		}
	}

	// script not in cache, and not found in checks directory
	err = fmt.Errorf("couldn't find script for check '%s'", id)
	return "", err
}

type scriptOutput struct {
	status bool
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (cfg *agentConfig) runScript(script scriptPath, c types.Service) (scriptOutput, error) {
	var output scriptOutput

	// create a temporary file called {{uuid}}.json
	uid, err := uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("failed to generate a random uuid. this should never happen: %w", err)
		return output, err
	}
	tempFileName := uid.String() + ".json"
	tempFilePath := filepath.Join(cfg.tempDir, tempFileName)

	// convert service to json
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		err = fmt.Errorf("failed to marshal json for check '%s': %w", c.CheckName, err)
		return output, err
	}

	// write json to randomly named file
	err = os.WriteFile(tempFilePath, bytes, 0666)
	if err != nil {
		err = fmt.Errorf("failed to write temp file '%s': %w", tempFilePath, err)
		return output, err
	}

	// nuke that shit from orbit!!
	defer os.Remove(tempFilePath)

	// create timeout context to kill commands that run too long
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create command
	cmd := exec.CommandContext(ctx, string(script), tempFilePath)
	cmd.Stdout = &output.stdout
	cmd.Stderr = &output.stderr

	// start it and block until completion
	err = cmd.Run()

	output.status = (err == nil)

	return output, nil
}
