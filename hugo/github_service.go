package hugo

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type repositoryDetails struct {
	organisation string
	repository   string
}

var repository *repositoryDetails

var client *github.Client
var latestRelease *github.RepositoryRelease
var latestVersion *Version

func init() {
	client = github.NewClient(nil)
	repository = &repositoryDetails{organisation: "gohugoio", repository: "hugo"}
}

func (repo *repositoryDetails) getRelease(version *Version) (*github.RepositoryRelease, error) {
	if _, err := repo.getLatestVersion(); err == nil && version == latestVersion {
		return latestRelease, nil
	}
	release, _, err := client.Repositories.GetReleaseByTag(context.TODO(), repo.organisation, repo.repository, version.toGithubTag())
	if err != nil {
		return nil, errors.Wrap(err, "not being able to get repository")
	}
	return release, nil
}
func (repo *repositoryDetails) getLatestRelease() (*github.RepositoryRelease, error) {
	if latestRelease == nil {
		latestRelease, _, err := client.Repositories.GetLatestRelease(context.TODO(), repo.organisation, repo.repository)
		return latestRelease, err
	}
	return latestRelease, nil
}

func (repo *repositoryDetails) getLatestVersion() (version *Version, err error) {
	release, err := repo.getLatestRelease()
	if err != nil {
		return nil, errors.Cause(err)
	}
	return NewVersion(release.GetName())
}

type releasePager struct {
	repository      *repositoryDetails
	currentReleases []*github.RepositoryRelease
	currentResponse *github.Response
	opt             *github.ListOptions
}

func (repo *repositoryDetails) newReleasePager() (*releasePager, error) {
	pager := new(releasePager)
	pager.opt = &github.ListOptions{Page: 1}
	releases, response, err := client.Repositories.ListReleases(context.TODO(), repo.organisation, repo.repository, pager.opt)
	if err != nil {
		return nil, err
	}
	pager.currentReleases = releases
	pager.currentResponse = response
	pager.repository = repo
	return pager, nil
}

func (pager *releasePager) hasMore() bool {
	if pager.currentResponse == nil {
		panic("the pager should have been initialize with newReleasePager")
	}
	return pager.currentResponse.NextPage != 0
}

func (pager *releasePager) getNextPage() error {
	if pager.currentResponse == nil {
		panic("the pager should have been initialize with newReleasePager")
	}
	if !pager.hasMore() {
		panic("no more releases to be found")
	}
	pager.opt.Page += 1
	releases, response, err := client.Repositories.ListReleases(context.TODO(), pager.repository.organisation, pager.repository.repository, pager.opt)
	if err != nil {
		return err
	}
	pager.currentReleases = releases
	pager.currentResponse = response
	return nil
}

func (repo *repositoryDetails) findHighestRelease(version *Version, compareOn comparatorKey) (*Version, error) {
	pager, err := repo.newReleasePager()
	if err != nil {
		return nil, errors.Cause(err)
	}
	for {
		latestVersionOnPage, err := NewVersion(pager.currentReleases[len(pager.currentReleases)-1].GetName())
		if err != nil {
			return nil, err
		}
		if version.isEqual(latestVersionOnPage, compareOn) {
			return latestVersionOnPage, nil
		}
		if version.isHigherOrEqual(latestVersionOnPage, compareOn) {
			break
		}
		if !pager.hasMore() {
			return nil, fmt.Errorf("The version %s version has not been found", version)
		}
		if err := pager.getNextPage(); err != nil {
			return nil, errors.Cause(err)
		}
	}
	i := sort.Search(len(pager.currentReleases)-1, func(i int) bool {
		other, err := NewVersion(pager.currentReleases[i].GetName())
		if err != nil {
			panic("not able to parse the version")
		}
		return version.isHigherOrEqual(other, compareOn)
	})
	if i < len(pager.currentReleases) {
		found, err := NewVersion(pager.currentReleases[i].GetName())
		if err != nil {
			return nil, err
		}
		if version.isEqual(found, compareOn) {
			return found, nil
		}
		return nil, fmt.Errorf("The version %s version has not been found", version)
	}
	return nil, fmt.Errorf("The version %s version has not been found", version)
}

func (repo *repositoryDetails) getHighestMinor(major int) (*Version, error) {
	latestVersion, err := repo.getLatestVersion()
	if err != nil {
		return nil, errors.Cause(err)
	}
	if major == latestVersion.Major {
		return latestVersion, nil
	}
	basedVersion := &Version{Major: major, Minor: 0, Patch: 0}
	return repo.findHighestRelease(basedVersion, compareOnMajor)
}

func (repo *repositoryDetails) getHighestPatch(major int, minor int) (*Version, error) {
	latestVersion, err := repo.getLatestVersion()
	if err != nil {
		return nil, errors.Cause(err)
	}
	if major == latestVersion.Major && minor == latestVersion.Minor {
		return latestVersion, nil
	}
	basedVersion := &Version{Major: major, Minor: minor, Patch: 0}
	return repo.findHighestRelease(basedVersion, compareOnMinor)
}
