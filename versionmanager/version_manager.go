package versionmanager

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/TiboStev/hugo-wrapper/hugo"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

type VersionManager struct {
	installDirectory string
}
type HugoInstaller struct {
	installDirectory string
	selectedVersion  *hugo.Version
	execPath         string
}

func NewVersionManager(installDirectory string) (*VersionManager, error) {
	if _, err := os.Stat(installDirectory); err != nil {
		return nil, errors.New("The installation directory doesn't exist")
	}
	return &VersionManager{installDirectory: installDirectory}, nil
}

func (manager *VersionManager) GetExecPath(desiredVersion string) (execPath string, version string, err error) {
	selectedVersion, err := hugo.NewVersion(desiredVersion)
	version = selectedVersion.String()
	if err != nil {
		return
	}
	execPath = path.Join(manager.installDirectory, selectedVersion.String(), binaryName())
	if isAlreadyInstalled(execPath) {
		fmt.Println("found local installation")
		return
	}
	fmt.Println("local installation not found")
	fmt.Println("installation started ...")
	if err = manager.install(execPath, selectedVersion); err != nil {
		return
	}
	fmt.Println("installation ended")
	return
}

func (manager *VersionManager) install(execPath string, version *hugo.Version) (err error) {
	fmt.Println("download started")
	assetTmpFile, err := version.GetAsset()
	fmt.Println("download ended")
	defer assetTmpFile.Close()
	defer os.Remove(assetTmpFile.Name())
	if err != nil {
		return err
	}
	err = archiver.Unarchive(assetTmpFile.Name(), path.Dir(execPath))
	return
}

func isAlreadyInstalled(execPath string) bool {
	_, err := os.Stat(execPath)
	return err == nil
}

func binaryName() string {
	if runtime.GOOS == "windows" {
		return "hugo.exe"
	} else {
		return "hugo"
	}
}
