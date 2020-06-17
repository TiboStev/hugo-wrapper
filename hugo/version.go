package hugo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Version struct {
	Extended bool
	Major    int
	Minor    int
	Patch    int
}

func (version *Version) toGithubTag() string {
	return fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch)
}

func (version *Version) String() (versionString string) {
	versionString = fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch)
	if version.Extended {
		versionString += "-extended"
	}
	return
}

func extractVersionCore(version string) (versionCore string, isExtented bool, err error) {
	splitVersion := strings.Split(version, "-")
	if len(splitVersion) > 2 {
		return "", false, errors.New("the version must be in form of latest[-extended] or [v]int[.int[.int]][-extended]")
	}
	if splitVersion[0][0] == 'v' {
		splitVersion[0] = splitVersion[0][1:]
	}
	if len(splitVersion) == 1 {
		return splitVersion[0], false, nil
	}
	if splitVersion[1] != "extended" {
		return "", false, errors.New("the version must be in form of latest[-extended] or [v]int[.int[.int]][-extended]")
	}
	return splitVersion[0], true, nil
}

func isLatest(version string, isExtented bool) (bool, *Version, error) {
	if version == "latest" {
		version, err := repository.getLatestVersion()
		if err != nil {
			return false, nil, errors.Cause(err)
		}
		version.Extended = isExtented
		return true, version, err
	}
	return false, nil, nil
}

func NewVersion(v string) (*Version, error) {
	var err error
	var isExtended bool

	v, isExtended, err = extractVersionCore(v)

	isLatest, version, err := isLatest(v, isExtended)
	if err != nil {
		return nil, errors.Cause(err)
	}
	if isLatest {
		fmt.Println("it's the latest")
		return version, nil
	}

	splitVersion := strings.Split(v, ".")
	if len(splitVersion) > 3 {
		return nil, errors.New("the version must be in form of latest[-extended] or [v]int[.int[.int]][-extended]")
	}

	version = new(Version)
	version.Extended = isExtended

	version.Major, version.Minor, version.Patch, err = extractIdentifiersFromVersionCore(splitVersion)

	if err != nil {
		return nil, errors.Cause(err)
	}

	switch len(splitVersion) {
	case 1:
		majorVersion, err := repository.getHighestMinor(version.Major)
		majorVersion.Extended = isExtended
		return majorVersion, err
	case 2:
		if version.Major == 0 && version.Minor <= 53 {
			return version, nil
		}
		minorVersion, err := repository.getHighestPatch(version.Major, version.Minor)
		minorVersion.Extended = isExtended
		return minorVersion, err
	case 3:
		return version, nil
	}
	return nil, nil
}

func extractIdentifiersFromVersionCore(splitVersion []string) (major int, minor int, fix int, err error) {
	switch len(splitVersion) {
	case 3:
		fix, err = strconv.Atoi(splitVersion[2])
		fallthrough
	case 2:
		minor, err = strconv.Atoi(splitVersion[1])
		fallthrough
	case 1:
		major, err = strconv.Atoi(splitVersion[0])
	default:
		return 0, 0, 0, errors.New("the version must be in form of latest[-extended] or [v]int[.int[.int]][-extended]")
	}
	return
}

type comparatorKey int

const (
	compareOnMajor = comparatorKey(1)
	compareOnMinor = comparatorKey(2)
	compareOnFix   = comparatorKey(3)
)

func (version *Version) isHigherOrEqual(other *Version, compareOn comparatorKey) bool {
	switch compareOn {
	case compareOnFix:
		if version.Major == other.Major && version.Minor == other.Minor {
			return (version.Patch - other.Patch) >= 0
		}
		fallthrough
	case compareOnMinor:
		if version.Major == other.Major {
			return (version.Minor - other.Minor) >= 0
		}
		fallthrough
	case compareOnMajor:
		return (version.Major - other.Major) >= 0
	default:
		panicError := fmt.Sprintf("fatal the comparatorKey: %d doesn't exist", compareOn)
		panic(panicError)
	}
}

func (version *Version) isEqual(other *Version, compareOn comparatorKey) bool {
	switch compareOn {
	case compareOnFix:
		if version.Major == other.Major && version.Minor == other.Minor {
			return (version.Patch - other.Patch) == 0
		}
		return false
	case compareOnMinor:
		if version.Major == other.Major {
			return (version.Minor - other.Minor) == 0
		}
		return false
	case compareOnMajor:
		return (version.Major - other.Major) == 0
	default:
		panicError := fmt.Sprintf("fatal the comparatorKey: %d doesn't exist", compareOn)
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

func (version *Version) GetAsset() (f *os.File, err error) {
	normalizedAssetName, err := version.toAssetName()
	if err != nil {
		return nil, err
	}
	normalizedAssetName = strings.ToLower(normalizedAssetName)
	release, err := repository.getRelease(version)
	if err != nil {
		return nil, err
	}
	for _, asset := range release.Assets {
		if normalizedAssetName == strings.ToLower(asset.GetName()) {
			resp, err := http.Get(asset.GetBrowserDownloadURL())
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
	}
	return nil, fmt.Errorf("no asset has been found for version %s", version.String())
}

func (version *Version) toAssetName() (assetName string, err error) {
	var builder strings.Builder
	builder.WriteString("hugo_")
	if version.Extended {
		builder.WriteString("extended_")
	}
	fmt.Fprintf(&builder, "%d.%d", version.Major, version.Minor)
	if !(version.Major == 0 && version.Minor <= 53 && version.Patch == 0) {
		fmt.Fprintf(&builder, ".%d", version.Patch)
	}
	arch, err := getArch()

	fmt.Fprintf(&builder, "_%s-%s.%s", runtime.GOOS, arch, getExtension())
	return builder.String(), err
}

func getExtension() string {
	if runtime.GOOS == "windows" {
		return "zip"
	} else {
		return "tar.gz"
	}
}

func getArch() (arch string, err error) {
	switch runtime.GOARCH {
	case "386":
		arch = "32bit"
	case "amd64":
		arch = "64bit"
	case "arm":
		arch = "arm"
	case "arm64":
		arch = "arm64"
	default:
		err = errors.New("There is no hugo for this architecture")
	}
	return
}
