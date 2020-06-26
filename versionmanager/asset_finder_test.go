package versionmanager

import (
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type assetNameTestingItem struct {
	os            string
	arch          string
	version       Version
	desiredOutput string
}

func TestAssetName(t *testing.T) {
	assetNameCoreVersion := coreVersion{major: 0, minor: 73, patch: 0}
	assetNameCoreOldVersion := coreVersion{major: 0, minor: 53, patch: 0}
	assetNameTestingVersion := Version{coreVersion: &assetNameCoreVersion, extended: false, finder: nil}
	assetNameTestingExtentedVersion := Version{coreVersion: &assetNameCoreVersion, extended: true, finder: nil}
	assetNameTestingOldVersion := Version{coreVersion: &assetNameCoreOldVersion, extended: false, finder: nil}

	assetNameTestingList := []assetNameTestingItem{
		assetNameTestingItem{os: "darwin", arch: "amd64", desiredOutput: "hugo_0.73.0_macOS-64bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "dragonfly", arch: "amd64", desiredOutput: "hugo_0.73.0_DragonFlyBSD-64bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "freebsd", arch: "amd64", desiredOutput: "hugo_0.73.0_FreeBSD-64bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "linux", arch: "amd64", desiredOutput: "hugo_0.73.0_Linux-64bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "netbsd", arch: "amd64", desiredOutput: "hugo_0.73.0_NetBSD-64bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "openbsd", arch: "amd64", desiredOutput: "hugo_0.73.0_OpenBSD-64bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "windows", arch: "amd64", desiredOutput: "hugo_0.73.0_Windows-64bit.zip", version: assetNameTestingVersion},

		assetNameTestingItem{os: "darwin", arch: "386", desiredOutput: "hugo_0.73.0_macOS-32bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "freebsd", arch: "386", desiredOutput: "hugo_0.73.0_FreeBSD-32bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "linux", arch: "386", desiredOutput: "hugo_0.73.0_Linux-32bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "netbsd", arch: "386", desiredOutput: "hugo_0.73.0_NetBSD-32bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "openbsd", arch: "386", desiredOutput: "hugo_0.73.0_OpenBSD-32bit.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "windows", arch: "386", desiredOutput: "hugo_0.73.0_Windows-32bit.zip", version: assetNameTestingVersion},

		assetNameTestingItem{os: "freebsd", arch: "arm", desiredOutput: "hugo_0.73.0_FreeBSD-ARM.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "linux", arch: "arm", desiredOutput: "hugo_0.73.0_Linux-ARM.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "netbsd", arch: "arm", desiredOutput: "hugo_0.73.0_NetBSD-ARM.tar.gz", version: assetNameTestingVersion},
		assetNameTestingItem{os: "openbsd", arch: "arm", desiredOutput: "hugo_0.73.0_OpenBSD-ARM.tar.gz", version: assetNameTestingVersion},

		assetNameTestingItem{os: "linux", arch: "arm64", desiredOutput: "hugo_0.73.0_Linux-ARM64.tar.gz", version: assetNameTestingVersion},

		assetNameTestingItem{os: "darwin", arch: "amd64", desiredOutput: "hugo_extended_0.73.0_macOS-64bit.tar.gz", version: assetNameTestingExtentedVersion},
		assetNameTestingItem{os: "linux", arch: "amd64", desiredOutput: "hugo_extended_0.73.0_Linux-64bit.tar.gz", version: assetNameTestingExtentedVersion},
		assetNameTestingItem{os: "windows", arch: "amd64", desiredOutput: "hugo_extended_0.73.0_Windows-64bit.zip", version: assetNameTestingExtentedVersion},

		assetNameTestingItem{os: "darwin", arch: "amd64", desiredOutput: "hugo_0.53_macOS-64bit.tar.gz", version: assetNameTestingOldVersion},
		assetNameTestingItem{os: "dragonfly", arch: "amd64", desiredOutput: "hugo_0.53_DragonFlyBSD-64bit.tar.gz", version: assetNameTestingOldVersion},
		assetNameTestingItem{os: "freebsd", arch: "amd64", desiredOutput: "hugo_0.53_FreeBSD-64bit.tar.gz", version: assetNameTestingOldVersion},
		assetNameTestingItem{os: "linux", arch: "amd64", desiredOutput: "hugo_0.53_Linux-64bit.tar.gz", version: assetNameTestingOldVersion},
		assetNameTestingItem{os: "netbsd", arch: "amd64", desiredOutput: "hugo_0.53_NetBSD-64bit.tar.gz", version: assetNameTestingOldVersion},
		assetNameTestingItem{os: "openbsd", arch: "amd64", desiredOutput: "hugo_0.53_OpenBSD-64bit.tar.gz", version: assetNameTestingOldVersion},
		assetNameTestingItem{os: "windows", arch: "amd64", desiredOutput: "hugo_0.53_Windows-64bit.zip", version: assetNameTestingOldVersion},
	}

	assert := assert.New(t)
	for _, item := range assetNameTestingList {
		setOs(item.os)
		setArch(item.arch)
		actualOutput, _ := assetName(&item.version)
		assert.Equal(item.desiredOutput, actualOutput)
	}

	resetOs()
	resetArch()
}

func resetOs() {
	monkey.Unpatch(goOS)
}

func resetArch() {
	monkey.Unpatch(goArch)
}

func setOs(os string) {
	monkey.Patch(goOS, func() string {
		return os
	})
}

func setArch(arch string) {
	monkey.Patch(goArch, func() string {
		return arch
	})
}

func TestFindAssetURL(t *testing.T) {

}
func TestResolveVersion(t *testing.T) {
	testResolveMajor()
	testResolveMinor()
	testResolvePatch()
}

func TestReleaseTag(t *testing.T) {
	assert := assert.New(t)

	// after 0.53 minor releases start at major.minor.0
	assert.Equal(releaseTag(&coreVersion{major: 1, minor: 42, patch: 0}), "v1.42.0")

	// prior to 0.53 minor releases start at major.minor
	assert.Equal(releaseTag(&coreVersion{major: 0, minor: 53, patch: 0}), "v0.53")
	assert.Equal(releaseTag(&coreVersion{major: 0, minor: 53, patch: 1}), "v0.53.1")
	assert.Equal(releaseTag(&coreVersion{major: 0, minor: 42, patch: 0}), "v0.42")
}

func TestGetExtension(t *testing.T) {
	assert := assert.New(t)

	osToDesiredExtension := map[string]string{
		"darwin":    ".tar.gz",
		"dragonfly": ".tar.gz",
		"freebsd":   ".tar.gz",
		"linux":     ".tar.gz",
		"netbsd":    ".tar.gz",
		"openbsd":   ".tar.gz",
		"windows":   ".zip",
	}

	for os, desireExtension := range osToDesiredExtension {
		setOs(os)
		assert.Equal(desireExtension, getExtension())
	}

	resetOs()
	resetArch()
}

func testResolveMajor() {

}
func testResolveMinor() {

}
func testResolvePatch() {

}

func TestFindLatestVersion(t *testing.T) {
	whenLatestVersionHasNotBeenFetchedYet(t)
	//whenLatestVersionHasAlreadyBeenFetched(t)

}

func whenLatestVersionHasNotBeenFetchedYet(t *testing.T) {
	finder := new(finder)
	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	repository := NewMockRepositoryClient(ctrl)
	release := NewMockRelease(ctrl)

	repository.EXPECT().GetLatestRelease().Return(release, nil)
	release.EXPECT().GetName().Return("v0.72.3")

	finder.repository = repository
	finder.findLatestVersion()
	assert.Equal(t, 0, finder.latestSelectedVersion.major)
	assert.Equal(t, 72, finder.latestSelectedVersion.minor)
	assert.Equal(t, 3, finder.latestSelectedVersion.patch)
	if finder.latestSelectedRelease != release {
		t.Errorf("the latest release should be the one returned by the repository service")
	}
}

func whenLatestVersionHasAlreadyBeenFetched(t *testing.T) {
	finder := new(finder)
	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	finder.repository = NewMockRepositoryClient(ctrl)

}
