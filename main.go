package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type version struct {
	major, minor, patch int
}

func parse(raw string) *version {
	values := strings.Split(raw, ".")

	if len(values) != 3 {
		return nil
	}

	major, err := strconv.Atoi(values[0])
	minor, err := strconv.Atoi(values[1])
	patch, err := strconv.Atoi(values[2])
	if err != nil {
		return nil
	}

	return &version{
		major: major,
		minor: minor,
		patch: patch}
}

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

	fmt.Println(selectedStep)
	for _, info := range infos {
		fmt.Println(info.Name())
	}

	sort.Slice(infos, func(i, j int) bool {
		iVersion := parse(infos[i].Name())
		jVersion := parse(infos[j].Name())
		if iVersion == nil {
			return true
		}
		if jVersion == nil {
			return false
		}

		if iVersion.major != jVersion.major {
			return iVersion.major < jVersion.major
		}
		if iVersion.minor != jVersion.minor {
			return iVersion.major < jVersion.major
		}
	})
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
