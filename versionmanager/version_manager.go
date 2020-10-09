package versionmanager

import (
	"fmt"
	"os"
	"path"
	"runtime"

	archiver "github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
)

type VersionManager struct {
	installDirectory string
}
type HugoInstaller struct {
	installDirectory string
	selectedVersion  *Version
	execPath         string
}

func NewVersionManager(installDirectory string) (*VersionManager, error) {
	if _, err := os.Stat(installDirectory); err != nil {
		return nil, errors.New("The installation directory doesn't exist")
	}
	return &VersionManager{installDirectory: installDirectory}, nil
}

func (manager *VersionManager) GetExecPath(desiredVersion string) (execPath string, version string, err error) {
	selectedVersion, err := NewVersion(desiredVersion)
	if err != nil {
		return
	}
	fmt.Printf("%+v", selectedVersion)
	version = selectedVersion.String()
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

func (manager *VersionManager) install(execPath string, version *Version) (err error) {
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
	}
	return "hugo"
}
