package versionmanager

import (
	"fmt"
	"runtime"
	"strings"
)

type assetFinder interface {
	findLatestVersion() (version *coreVersion, err error)
	findAssetURL(version *Version) (downloadUrl string, err error)
	resolveVersion(desiredVersion *coreVersion, compareOn versionPrecision) (*coreVersion, error)
}

type finder struct {
	repository            RepositoryClient
	latestRelease         Release
	latestVersion         *coreVersion
	latestSelectedRelease Release
	latestSelectedVersion *coreVersion
}

func newAssetFinder() (assetFinder assetFinder) {
	finder := new(finder)
	finder.repository = NewRepositoryService(Github, "gohugoio", "hugo", "", "")
	return finder
}

func (finder *finder) findLatestVersion() (version *coreVersion, err error) {
	if finder.latestVersion == nil {
		finder.latestSelectedRelease, err = finder.repository.GetLatestRelease()
		if err != nil {
			return nil, err
		}
		finder.latestSelectedVersion, _, err = parseCoreVersion(finder.latestSelectedRelease.GetName())
		if err != nil {
			return nil, err
		}
	}
	fmt.Printf("%v+", finder.latestSelectedVersion)
	return finder.latestSelectedVersion, nil
}

func (finder *finder) findAssetURL(version *Version) (downloadUrl string, err error) {
	if finder.latestSelectedVersion == nil || !finder.latestSelectedVersion.Equal(version.coreVersion, patch) {
		_, err = finder.resolveVersion(version.coreVersion, patch)
		if err != nil {
			return "", err
		}
	}
	assetName, err := assetName(version)
	if err != nil {
		return "", err
	}
	asset, err := finder.latestSelectedRelease.GetAssetByName(assetName)
	if err != nil {
		return "", err
	}
	return asset.GetDownloadUrl(), nil
}

func (finder *finder) resolveVersion(desiredVersion *coreVersion, precision versionPrecision) (*coreVersion, error) {
	latestVersion, err := finder.findLatestVersion()
	if err != nil {
		return nil, err
	}
	if desiredVersion.Higher(latestVersion, precision) {
		return nil, fmt.Errorf("the requested version is higher than the latest version available, latest available is %s", releaseTag(latestVersion))
	}
	if desiredVersion.Equal(latestVersion, precision) {
		return latestVersion, nil
	}
	switch precision {
	case major:
		nextMajorVersion := &coreVersion{major: desiredVersion.major + 1, minor: 0, patch: 0}
		finder.latestSelectedRelease, err = finder.repository.GetPreviousRelease(releaseTag(nextMajorVersion))
	case minor:
		nextMinorVersion := &coreVersion{major: desiredVersion.major, minor: desiredVersion.minor + 1, patch: 0}
		finder.latestSelectedRelease, err = finder.repository.GetPreviousRelease(releaseTag(nextMinorVersion))
	case patch:
		finder.latestSelectedRelease, err = finder.repository.GetReleaseByTag(releaseTag(desiredVersion))
	default:
		panic("Can't fetch version, the comparator key is not known")
	}
	if err != nil {
		return nil, err
	}
	finder.latestSelectedVersion, _, err = parseCoreVersion(finder.latestSelectedRelease.GetName())
	return finder.latestSelectedVersion, err
}

var osToAssetOs = map[string]string{
	"darwin":    "macOS",
	"dragonfly": "DragonFlyBSD",
	"freebsd":   "FreeBSD",
	"linux":     "Linux",
	"netbsd":    "NetBSD",
	"openbsd":   "OpenBSD",
	"windows":   "Windows",
}

var archToAssetArch = map[string]string{
	"386":   "32bit",
	"amd64": "64bit",
	"arm":   "ARM",
	"arm64": "ARM64",
}

func goOS() string {
	return runtime.GOOS
}

func goArch() string {
	return runtime.GOARCH
}

func assetName(version *Version) (assetName string, err error) {
	var builder strings.Builder
	builder.WriteString("hugo_")
	if version.extended {
		builder.WriteString("extended_")
	}
	fmt.Printf("%s", version)
	fmt.Printf("%d", version.coreVersion.major)
	fmt.Printf("%d", version.minor)
	fmt.Fprintf(&builder, "%d.%d", version.major, version.minor)
	if !(version.major == 0 && version.minor <= 53 && version.patch == 0) {
		fmt.Fprintf(&builder, ".%d", version.patch)
	}

	fmt.Fprintf(&builder, "_%s-%s%s", osToAssetOs[goOS()], archToAssetArch[goArch()], getExtension())
	return builder.String(), err
}

func releaseTag(version *coreVersion) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "v%d.%d", version.major, version.minor)
	if !(version.major == 0 && version.minor <= 53 && version.patch == 0) {
		fmt.Fprintf(&builder, ".%d", version.patch)
	}
	return builder.String()
}

func getExtension() string {
	if goOS() == "windows" {
		return ".zip"
	}
	return ".tar.gz"
}
