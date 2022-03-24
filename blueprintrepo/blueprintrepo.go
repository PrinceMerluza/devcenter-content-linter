package blueprintrepo

import (
	"errors"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/PrinceMerluza/devcenter-content-linter/logger"
)

var (
	remoteUrl   string // For results. If remote repo, it will be the remote path instead of the cloned dir
	workingPath string
)

// Use a repository to lint (local or remote)
// If local, doesn't actually need to be a git repository
func UseRepo(repoPath string, isRemote bool) {
	if isRemote {
		tmpPath, err := cloneRepoTemp(repoPath)
		if err != nil {
			logger.Fatal(err)
		}

		remoteUrl = repoPath
		workingPath = tmpPath
		return
	}

	workingPath = repoPath
}

// Get the working path of the repo
func GetWorkingPath() string {
	return workingPath
}

// Depending on if the repo is local or remote, rebuild the path/URL to point to a relative path
func GetOriginalRelPath(localPath string) string {
	if remoteUrl == "" {
		return localPath
	}

	relPath, err := filepath.Rel(workingPath, localPath)
	if err != nil {
		logger.Fatal(err)
	}

	u, err := url.Parse(remoteUrl)
	if err != nil {
		logger.Fatal(err)
	}

	u.Path = path.Join(u.Path, relPath)
	s := u.String()

	return s
}

func cloneRepoTemp(repoUrl string) (string, error) {
	tmpPath, err := os.MkdirTemp("", "gc-content")
	if err != nil {
		logger.Warn("Error creating temp dir:", err)
		return "", err
	}

	logger.Info("Cloning blueprint...")

	// Clone the blueprint into the temporary directory
	_, err = exec.Command("git", "-C", tmpPath, "clone", repoUrl).Output()
	if err != nil {
		logger.Warn("Error cloning repo:", err)
		return "", err
	}

	files, err := os.ReadDir(tmpPath)
	if err != nil {
		return "", err
	}

	if len(files) < 1 {
		err = errors.New("can't find cloned repo directory")
		return "", err
	}

	logger.Info("Successfully cloned blueprint")
	dirPath := filepath.Join(tmpPath, files[0].Name())

	return dirPath, nil
}
