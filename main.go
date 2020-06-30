package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adborbas/latest_version_step/version"
)

func main() {

	stepID, err := readStepID()
	if err != nil {
		failf("Could not read step ID. %s", err)
	}

	repo, err := cloneSteplib()
	defer os.RemoveAll(repo)
	if err != nil {
		failf("Could not read step ID. %s", err)
	}

	selectedStep := filepath.Join(repo, "steps", stepID)
	infos, err := ioutil.ReadDir(selectedStep)
	if err != nil {
		failf("Could not walk directories. %s", err)
	}

	sort.Slice(infos, func(i, j int) bool {
		iVersion := version.New(infos[i].Name())
		jVersion := version.New(infos[j].Name())
		if iVersion == nil {
			return true
		}
		if jVersion == nil {
			return false
		}

		return !iVersion.IsNewer(*jVersion)
	})

	lastInfo := infos[len(infos)-1]
	fmt.Printf("Latest version of %s is: %s \n", stepID, lastInfo.Name())
}

func directories(root string) ([]string, error) {
	dirs := []string{}
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return dirs, nil
}

func failf(msg string, args ...interface{}) {
	fmt.Printf(msg+" ", args...)
	os.Exit(1)
}

func cloneSteplib() (string, error) {
	tmpDir, err := tempDir()
	if err != nil {
		return "", fmt.Errorf("Could not create tempDir. %s", err)
	}

	cmd := exec.Command("git", "clone", "https://github.com/bitrise-io/bitrise-steplib.git", tmpDir)
	if err = cmd.Run(); err != nil {
		return "", fmt.Errorf("Could not clone git repo. %s", err)
	}

	return tmpDir, nil
}

func readStepID() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Which Bitrise step's latest version are you looking for: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Replace(text, "\n", "", -1), nil
}

func tempDir() (string, error) {
	dir, err := ioutil.TempDir("", "git")
	if err != nil {
		return "", err
	}
	return dir, nil
}
