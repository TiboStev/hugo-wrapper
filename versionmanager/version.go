package versionmanager

import (
	"fmt"
	"io/ioutil"
	"net/http"
	osFile "os"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type versionPrecision int

const (
	major = versionPrecision(1)
	minor = versionPrecision(2)
	patch = versionPrecision(3)
)

type coreVersion struct {
	major int
	minor int
	patch int
}

type Version struct {
	*coreVersion
	extended bool
	finder   assetFinder
}

func NewVersion(desiredVersion string) (*Version, error) {
	finder := newAssetFinder()
	return newVersion(finder, desiredVersion)
}

func newVersion(finder assetFinder, desiredVersion string) (selectedVersion *Version, err error) {
	selectedVersion = new(Version)
	selectedVersion.finder = finder

	desiredVersion, isExtended, err := extractExtension(desiredVersion)
	selectedVersion.extended = isExtended

	if desiredVersion == "latest" {
		selectedVersion.coreVersion, err = finder.findLatestVersion()
		fmt.Println(desiredVersion)

		return selectedVersion, err
	}

	coreVersion, precision, err := parseCoreVersion(desiredVersion)
	if err != nil {
		return nil, err
	}

	selectedVersion.coreVersion, err = finder.resolveVersion(coreVersion, precision)
	return
}

func parseCoreVersion(version string) (*coreVersion, versionPrecision, error) {
	if version[0] == 'v' {
		version = version[1:]
	}
	coreVersion, precision, err := extractCoreIdentifiers(version)
	if err != nil {
		return nil, -1, err
	}
	if precision == minor && coreVersion.major == 0 && coreVersion.minor <= 53 {
		precision = patch
	}
	return coreVersion, precision, nil
}

func extractCoreIdentifiers(version string) (coreVer *coreVersion, precision versionPrecision, err error) {
	splitVersion := strings.Split(version, ".")
	if len(splitVersion) > 3 {
		return nil, -1, errors.New("the version must be in form of latest[-extended] or [v]int[.int[.int]][-extended]")
	}
	precision = versionPrecision(len(splitVersion))
	coreVer = new(coreVersion)
	switch precision {
	case patch:
		coreVer.patch, err = strconv.Atoi(splitVersion[2])
		fallthrough
	case minor:
		coreVer.minor, err = strconv.Atoi(splitVersion[1])
		fallthrough
	case major:
		coreVer.major, err = strconv.Atoi(splitVersion[0])
	default:
		err = errors.New("the version must be in form of latest[-extended] or [v]int[.int[.int]][-extended]")
	}
	return
}

func (version *Version) String() string {
	fmt.Printf("%+v", version.major)
	versionString := fmt.Sprintf("v%d.%d.%d", version.major, version.minor, version.patch)
	if version.extended {
		versionString += "-extended"
	}
	return versionString
}

func extractExtension(version string) (versionCore string, isExtented bool, err error) {
	splitVersion := strings.Split(version, "-")
	if len(splitVersion) > 2 {
		return "", false, errors.New("the version must be in form of latest[-extended] or [v]int[.int[.int]][-extended]")
	}
	if len(splitVersion) == 1 {
		return splitVersion[0], false, nil
	}
	if splitVersion[1] != "extended" {
		return "", false, errors.New("the version must be in form of latest[-extended] or [v]int[.int[.int]][-extended]")
	}
	return splitVersion[0], true, nil
}

func (version *coreVersion) Higher(other *coreVersion, precision versionPrecision) bool {
	switch precision {
	case patch:
		if version.major == other.major && version.minor == other.minor {
			return (version.patch - other.patch) > 0
		}
		fallthrough
	case minor:
		if version.major == other.major {
			return (version.minor - other.minor) > 0
		}
		fallthrough
	case major:
		return (version.major - other.major) > 0
	default:
		panicError := fmt.Sprintf("fatal the comparatorKey: %d doesn't exist", precision)
		panic(panicError)
	}
}

func (version *coreVersion) Equal(other *coreVersion, precision versionPrecision) bool {
	switch precision {
	case patch:
		if version.major == other.major && version.minor == other.minor {
			return (version.patch - other.patch) == 0
		}
		return false
	case minor:
		if version.major == other.major {
			return (version.minor - other.minor) == 0
		}
		return false
	case major:
		return (version.major - other.major) == 0
	default:
		panicError := fmt.Sprintf("fatal the comparatorKey: %d doesn't exist", precision)
		panic(panicError)
	}
}

func getTemporaryArchiveName() string {
	extension := ".tar.gz"
	if runtime.GOOS == "windows" {
		extension = ".zip"
	}
	return fmt.Sprintf("hugoArchive*.%s", extension)
}

func (version *Version) GetAsset() (f *osFile.File, err error) {
	url, err := version.finder.findAssetURL(version)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tmpfile, err := ioutil.TempFile("", getTemporaryArchiveName())
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if _, err := tmpfile.Write(content); err != nil {
		return nil, err
	}
	return tmpfile, nil
}
