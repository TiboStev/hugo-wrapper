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

type githubRepositoryServiceInterface interface {
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
	GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
}

type githubRepository struct {
	service      githubRepositoryServiceInterface
	organisation string
	repository   string
	pager        pager
}

type githubRelease struct {
	*github.RepositoryRelease
}

type githubAsset struct {
	*github.ReleaseAsset
}

func newGithubRepository(client *http.Client, organisation string, repository string) (repo *githubRepository) {
	repo = &githubRepository{
		organisation: organisation,
		repository:   repository,
		service:      github.NewClient(client).Repositories,
	}
	repo.pager = newReleasePager(repo)
	return
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
	release, err := repo.GetReleaseByTag(tag)
	if err != nil {
		return nil, errors.Cause(err)
	}
	gitRelease, err := repo.getPreviousReleaseWithRelease(release.(*githubRelease).RepositoryRelease)
	return &githubRelease{gitRelease}, err
}

func (repo *githubRepository) getPreviousReleaseWithRelease(pointerRelease *github.RepositoryRelease) (*github.RepositoryRelease, error) {
	pointerIndex, err := repo.pager.locateRelease(pointerRelease)
	if err != nil {
		return nil, err
	}
	releasesOnPage := repo.pager.releasesOnPage()
	if pointerIndex == (len(releasesOnPage) - 1) {
		if !repo.pager.hasMore() {
			return nil, fmt.Errorf("No previous release found for %s", pointerRelease.GetName())
		}
		if err = repo.pager.getNextPage(); err != nil {
			return nil, err
		}
		return repo.pager.releasesOnPage()[0], nil
	}
	return releasesOnPage[pointerIndex+1], nil
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

type pager interface {
	toStart() error
	moveToPageContainingRelease(release *github.RepositoryRelease) (err error)
	findReleaseOnPage(release *github.RepositoryRelease) (index int, err error)
	releasesOnPage() []*github.RepositoryRelease
	hasMore() bool
	getNextPage() error
	locateRelease(pointerRelease *github.RepositoryRelease) (int, error)
}

type releasePager struct {
	*githubRepository
	currentReleases []*github.RepositoryRelease
	currentResponse *github.Response
	opt             *github.ListOptions
}

func (pager *releasePager) releasesOnPage() []*github.RepositoryRelease {
	return pager.currentReleases
}

func newReleasePager(repo *githubRepository) (pager pager) {
	releasePager := new(releasePager)
	releasePager.githubRepository = repo
	return releasePager
}

func (pager *releasePager) hasMore() bool {
	return pager.currentResponse.NextPage != 0
}

func (pager *releasePager) toStart() (err error) {
	pager.opt = &github.ListOptions{Page: 1}
	return pager.listRelease()
}

func (pager *releasePager) getNextPage() (err error) {
	pager.opt.Page = pager.currentResponse.NextPage
	return pager.listRelease()
}

func (pager *releasePager) listRelease() (err error) {
	pager.currentReleases, pager.currentResponse, err = pager.service.ListReleases(context.TODO(), pager.organisation, pager.repository, pager.opt)
	return
}

func (pager *releasePager) locateRelease(pointerRelease *github.RepositoryRelease) (int, error) {
	err := pager.toStart()
	if err != nil {
		return -1, errors.Cause(err)
	}
	err = pager.moveToPageContainingRelease(pointerRelease)
	if err != nil {
		return -1, errors.Cause(err)
	}
	return pager.findReleaseOnPage(pointerRelease)
}

func (pager *releasePager) moveToPageContainingRelease(release *github.RepositoryRelease) (err error) {
	isFound := false
	for isFound == false {
		olderReleaseOnPage := pager.currentReleases[len(pager.currentReleases)-1]
		isFound = release.CreatedAt.Time.After(olderReleaseOnPage.CreatedAt.Time) || release.CreatedAt.Time.Equal(olderReleaseOnPage.CreatedAt.Time)
		if isFound {
			break
		}
		if !pager.hasMore() {
			return fmt.Errorf("The release %s version has not been found", release.GetName())
		}
		if err = pager.getNextPage(); err != nil {
			return err
		}
	}
	return nil
}

func (pager *releasePager) findReleaseOnPage(release *github.RepositoryRelease) (index int, err error) {
	i := sort.Search(len(pager.currentReleases), func(i int) bool {
		return (release.CreatedAt.Time.After(pager.currentReleases[i].CreatedAt.Time) || release.CreatedAt.Time.Equal(pager.currentReleases[i].CreatedAt.Time))
	})
	if i < len(pager.currentReleases) && release.CreatedAt.Time.Equal(pager.currentReleases[i].CreatedAt.Time) {
		return i, nil
	}
	return -1, fmt.Errorf("The version %s version has not been found", release.GetName())
}
