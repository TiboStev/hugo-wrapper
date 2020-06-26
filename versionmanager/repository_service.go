package versionmanager

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-github/v31/github"
	"github.com/pkg/errors"
)

type RepositoryClient interface {
	GetLatestRelease() (Release, error)
	GetReleaseByTag(tag string) (Release, error)
	GetPreviousRelease(tag string) (Release, error)
}

type Release interface {
	GetName() string
	GetAssetByName(name string) (Asset, error)
}

type Asset interface {
	GetName() string
	GetDownloadUrl() string
}

type RepositoryType int

var Github = RepositoryType(1)

func NewRepositoryService(repoType RepositoryType, organisation string, repository string, username string, password string) RepositoryClient {
	switch repoType {
	case Github:
		return newGithubRepository(&http.Client{}, organisation, repository)
	default:
		panic("no service for this repository type")
	}
}

type githubRepositoryInterface interface {
}

type githubRepositoryServiceInterface interface {
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
	GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
}

type githubRepository struct {
	service      githubRepositoryServiceInterface
	organisation string
	repository   string
}

type githubRelease struct {
	*github.RepositoryRelease
}

type githubAsset struct {
	*github.ReleaseAsset
}

func newGithubRepository(client *http.Client, organisation string, repository string) (repo *githubRepository) {
	return &githubRepository{
		organisation: organisation,
		repository:   repository,
		service:      github.NewClient(client).Repositories,
	}
}

func (repo *githubRepository) GetLatestRelease() (Release, error) {
	release, _, err := repo.service.GetLatestRelease(context.TODO(), repo.organisation, repo.repository)
	return &githubRelease{release}, err
}

func (repo *githubRepository) GetReleaseByTag(tag string) (Release, error) {
	release, _, err := repo.service.GetReleaseByTag(context.TODO(), repo.organisation, repo.repository, tag)
	return &githubRelease{release}, err
}

func (repo *githubRepository) GetPreviousRelease(tag string) (Release, error) {
	pager, err := repo.newReleasePager()
	release, err := repo.GetReleaseByTag(tag)
	if err != nil {
		return nil, errors.Cause(err)
	}
	pointer := release.(*githubRelease)
	err = pager.moveToPageContainingRelease(pointer)
	if err != nil {
		return nil, errors.Cause(err)
	}
	pointerIndex, err := pager.findReleaseOnPage(pointer)
	if err != nil {
		return nil, errors.Cause(err)
	}
	if pointerIndex == (len(pager.currentReleases) - 1) {
		if !pager.hasMore() {
			return nil, fmt.Errorf("No previous release found for %s", pointer.GetName())
		}
		if err = pager.getNextPage(); err != nil {
			return nil, err
		}
		return &githubRelease{pager.currentReleases[0]}, nil
	}
	return &githubRelease{pager.currentReleases[pointerIndex+1]}, nil
}

func (release *githubRelease) GetName() string {
	return release.RepositoryRelease.GetName()
}

func (release *githubRelease) GetAssetByName(name string) (Asset, error) {
	for _, asset := range release.Assets {
		if asset.GetName() == name {
			return githubAsset{asset}, nil
		}
	}
	return nil, fmt.Errorf("asset %s not found in release %s", name, release.GetName())
}

func (asset githubAsset) GetName() string {
	return asset.ReleaseAsset.GetName()
}
func (asset githubAsset) GetDownloadUrl() string {
	return asset.ReleaseAsset.GetBrowserDownloadURL()
}

type releasePager struct {
	*githubRepository
	currentReleases []*github.RepositoryRelease
	currentResponse *github.Response
	opt             *github.ListOptions
}

func (repo *githubRepository) newReleasePager() (*releasePager, error) {
	pager := new(releasePager)
	pager.githubRepository = repo
	err := pager.toStart()
	return pager, err
}

func (pager *releasePager) hasMore() bool {
	if pager.currentResponse == nil {
		panic("the pager should have been initialize with newReleasePager")
	}
	return pager.currentResponse.NextPage != 0
}

func (pager *releasePager) toStart() (err error) {
	pager.opt = &github.ListOptions{Page: 1}
	pager.currentReleases, pager.currentResponse, err = pager.service.ListReleases(context.TODO(), pager.organisation, pager.repository, pager.opt)
	return
}

func (pager *releasePager) getNextPage() (err error) {
	if pager.currentResponse == nil {
		panic("the pager should have been initialize with newReleasePager")
	}
	if !pager.hasMore() {
		panic("no more releases to be found")
	}
	pager.opt.Page++
	pager.currentReleases, pager.currentResponse, err = pager.service.ListReleases(context.TODO(), pager.organisation, pager.repository, pager.opt)
	return
}

func (pager *releasePager) moveToPageContainingRelease(release *githubRelease) (err error) {
	pager.toStart()
	isFound := false
	for isFound == false {
		olderReleaseOnPage := pager.currentReleases[len(pager.currentReleases)-1]
		isFound = release.GetCreatedAt().After(olderReleaseOnPage.GetCreatedAt().Time)
		if !pager.hasMore() {
			return fmt.Errorf("The release %s version has not been found", release.GetName())
		}
		if err = pager.getNextPage(); err != nil {
			return err
		}
	}
	return nil
}

func (pager *releasePager) findReleaseOnPage(release *githubRelease) (index int, err error) {
	i := sort.Search(len(pager.currentReleases)-1, func(i int) bool {
		return release.GetCreatedAt().After(pager.currentReleases[i].GetCreatedAt().Time)
	})
	if i < len(pager.currentReleases) && release.GetCreatedAt().Equal(pager.currentReleases[i].GetCreatedAt()) {
		return i, nil
	}
	return -1, fmt.Errorf("The version %s version has not been found", release.GetName())
}
