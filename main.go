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

	repo := make(chan string)
	fail := make(chan error)
	go cloneSteplib(repo, fail)

	stepID, err := readStepID()
	if err != nil {
		failf("Could not read step ID. %s", err)
	}

	rep, err := <-repo, <-fail
	if err != nil {
		failf("Could not clone repo. %s", err)
	}

	selectedStep := filepath.Join(rep, "steps", stepID)
	infos, err := ioutil.ReadDir(selectedStep)
	if err != nil {
		failf("Could not walk directories. %s", err)
	}

	versions := orderByVersion(infos)
	lastVersion := versions[len(versions)-1]
	fmt.Printf("Latest version of %s is: %s \n", stepID, lastVersion)
	defer os.RemoveAll(rep)
}

func orderByVersion(unordered []os.FileInfo) []version.Version {
	versions := []version.Version{}
	for _, fileinfo := range unordered {
		if version := version.New(fileinfo.Name()); version != nil {
			versions = append(versions, *version)
		}

	}
	sort.Slice(versions, func(i, j int) bool {
		iVersion := versions[i]
		jVersion := versions[j]

		return !iVersion.IsNewer(jVersion)
	})

	return versions
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

func cloneSteplib(repo chan string, err chan error) {
	defer func() {
		close(err)
		close(repo)
	}()

	tmpDir, fail := tempDir()
	if fail != nil {
		err <- fmt.Errorf("Could not create tempDir. %s", fail)
		return
	}

	cmd := exec.Command("git", "clone", "https://github.com/bitrise-io/bitrise-steplib.git", tmpDir)
	if fail = cmd.Run(); fail != nil {
		err <- fmt.Errorf("Could not clone git repo. %s", fail)
		return
	}

	repo <- tmpDir
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
