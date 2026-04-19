package agent

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// where scripts are searched for (in addition to PATH)
const scriptsDir = "scripts"

// Start runs the agent.
func Start(cfg *AgentConfig) {

	// find scripts directory
	scriptsPath, err := filepath.Abs(scriptsDir)
	if err != nil {
		fmt.Printf("failed to create fullpath for scripts directory: %v\n", err)
	}

	// add to path
	// this only matters for this process
	currentPath := os.Getenv("PATH")
	var newPath string
	if currentPath == "" {
		newPath = currentPath
	} else {
		sep := string(os.PathListSeparator) // : for linux, ; for windows
		newPath = scriptsPath + sep + currentPath
	}
	os.Setenv("PATH", newPath)

	// spawn requested number of threads
	for range cfg.NumThreads {
		go cfg.Loop()
	}

	// infinite loop
	// todo: implement graceful shutdown
	select {}
}

//go:embed default-scripts
var defaultScripts embed.FS

// DumpScripts creates the default scripts directory in the cwd.
func DumpScripts(forceOverwrite bool) {

	// make scripts directory
	err := os.MkdirAll(scriptsDir, 0755)
	if err != nil {
		fmt.Printf("failed to create scripts directory: %v\n", err)
		return
	}

	// get scripts from embed
	entries, err := defaultScripts.ReadDir("default-scripts")
	if err != nil {
		fmt.Printf("failed to find embedded scripts: %v\n", err)
		return
	}

	// store skipped writes
	var skipped = []string{}

	// loop through each file
	for _, entry := range entries {

		// file name and paths
		scriptName := entry.Name()
		embedPath := path.Join("default-scripts", scriptName) // path in the embed file system
		diskPath := path.Join(scriptsDir, scriptName)         // new path on disk

		// skip file if it already exists
		_, err := os.Stat(diskPath)
		if !errors.Is(err, os.ErrNotExist) && !forceOverwrite {

			// file exists, add to list of skipped
			skipped = append(skipped, scriptName)
			continue
		}

		// read file from embed
		script, err := defaultScripts.ReadFile(embedPath)
		if err != nil {
			fmt.Printf("failed to read script %q: %v\n", embedPath, err)
			os.Exit(1)
		}

		// write file to disk
		err = os.WriteFile(diskPath, script, 0755)
		if err != nil {
			fmt.Printf("failed to write file %q to disk: %v\n", scriptName, err)
			os.Exit(1)
		}
	}

	// print skipped writes
	if len(skipped) != 0 {
		fmt.Println("skipped the following files:")
		fmt.Printf("  %v\n", skipped)
		fmt.Println("use --force / -f to overwrite")
		os.Exit(1)
	}
}
