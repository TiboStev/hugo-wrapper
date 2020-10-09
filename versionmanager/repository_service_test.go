package versionmanager

import (
	context "context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	github "github.com/google/go-github/v31/github"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestRelease(t *testing.T) {
	repo, mockService, _ := getRepo(t)

	mockService.EXPECT().GetLatestRelease(context.TODO(), repo.organisation, repo.repository)

	repo.GetLatestRelease()
}

func TestGetReleaseByTag(t *testing.T) {
	repo, mockService, _ := getRepo(t)
	expectedName := "anyTag"

	mockService.EXPECT().GetReleaseByTag(context.TODO(), repo.organisation, repo.repository, expectedName)

	repo.GetReleaseByTag(expectedName)
}

func TestGetReleaseByTagErrorInGithubService(t *testing.T) {
	repo, mockService, _ := getRepo(t)
	expectedError := errors.New("expected error has occured")

	mockService.EXPECT().GetReleaseByTag(context.TODO(), repo.organisation, repo.repository, gomock.Any()).Return(nil, nil, expectedError)

	_, err := repo.GetReleaseByTag("anyTag")

	assert.Equal(t, expectedError, err, "error returned by githubrepo should be passed as it is")
}

func TestGetPreviousReleaseReleaseShouldCallGetPreviousReleaseWithPointerRelease(t *testing.T) {
	repo, mockService, pager := getRepo(t)
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease"})

	mockService.EXPECT().GetReleaseByTag(context.TODO(), "hugo", "gohugo", gomock.Any()).Return(pointerRelease, nil, nil)
	pager.EXPECT().locateRelease(&githubRelease{pointerRelease}).Return(-1, errors.New("no need to go further"))

	repo.GetPreviousRelease(*pointerRelease.Name)
}

func TestGetPreviousReleaseReleaseWithTagReturnError(t *testing.T) {
	repo, mockService, _ := getRepo(t)
	expectedError := errors.New("error while fetching release")

	mockService.EXPECT().GetReleaseByTag(context.TODO(), "hugo", "gohugo", gomock.Any()).Return(nil, nil, expectedError)

	_, actualError := repo.GetPreviousRelease("any")

	assert.Equal(t, expectedError, actualError, "error returned by githubrepo should be passed as it is")
}

func TestGetPreviousReleaseWithReleasePagerFailedWhileLocatingPreviousRelease(t *testing.T) {
	repo, _, mockPager := getRepo(t)
	dummyRelease := createDummyGithubRelease(releaseDescriptor{name: "dummyRelease"})
	expectedError := errors.New("error while locating release")

	mockPager.EXPECT().locateRelease(gomock.Any()).Return(-1, expectedError)

	_, actualError := repo.getPreviousReleaseWithRelease(dummyRelease)

	assert.Equal(t, expectedError, actualError, "error returned by githubrepo should be passed as it is")
}

func TestGetPreviousReleaseWithReleaseTagIsTheLastReleaseOfTheLatestPage(t *testing.T) {
	repo, _, mockPager := getRepo(t)

	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointer"})
	expectedError := fmt.Errorf("No previous release found for %s", *pointerRelease.Name)

	expecterPointerReleaseToBeLastOnThePage(mockPager, pointerRelease)
	mockPager.EXPECT().hasMore().Return(false)

	_, actualError := repo.getPreviousReleaseWithRelease(pointerRelease)

	assert.Equal(t, expectedError, actualError)
}

func TestGetPreviousReleaseWithReleasePreviousReleaseIsOnTheNextPage(t *testing.T) {
	repo, _, mockPager := getRepo(t)
	expectedPreviousRelease := createDummyGithubRelease(releaseDescriptor{name: "expectedRelease"})
	releasesOnNextPage := []*github.RepositoryRelease{expectedPreviousRelease}
	pointerRelease := createDummyGithubRelease(releaseDescriptor{})

	expecterPointerReleaseToBeLastOnThePage(mockPager, pointerRelease)
	mockPager.EXPECT().hasMore().Return(true)
	mockPager.EXPECT().getNextPage()
	mockPager.EXPECT().releasesOnPage().Return(releasesOnNextPage)

	actualRelease, err := repo.getPreviousReleaseWithRelease(pointerRelease)

	assert.Nil(t, err, "if a release has been found then no error should be returned")
	assert.Equal(t, expectedPreviousRelease.GetName(), actualRelease.GetName())
}

func expecterPointerReleaseToBeLastOnThePage(mockPager *Mockpager, pointerRelease *github.RepositoryRelease) {
	releasesOnPage := []*github.RepositoryRelease{pointerRelease}
	mockPager.EXPECT().locateRelease(gomock.Any()).Return(0, nil)
	mockPager.EXPECT().releasesOnPage().Return(releasesOnPage)
}

func TestGetPreviousReleaseWithReleasePreviousReleaseIsOnThePage(t *testing.T) {
	repo, _, mockPager := getRepo(t)
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease"})
	previousRelease := createDummyGithubRelease(releaseDescriptor{name: "previousRelease"})
	releasesOnPage := []*github.RepositoryRelease{pointerRelease, previousRelease}

	mockPager.EXPECT().locateRelease(gomock.Any()).Return(0, nil)
	mockPager.EXPECT().releasesOnPage().Return(releasesOnPage)

	actualRelease, err := repo.getPreviousReleaseWithRelease(pointerRelease)

	assert := assert.New(t)
	assert.Nil(err, "if a release has been found then no error should be returnde")
	assert.Equal(previousRelease, actualRelease)
}

func TestGetPreviousReleaseGetNextPageFailed(t *testing.T) {
	repo, _, mockPager := getRepo(t)
	pointerRelease := createDummyGithubRelease(releaseDescriptor{})
	expectedError := errors.New("getNextPage() failed")

	expecterPointerReleaseToBeLastOnThePage(mockPager, pointerRelease)
	mockPager.EXPECT().hasMore().Return(true)
	mockPager.EXPECT().getNextPage().Return(expectedError)

	_, actualError := repo.getPreviousReleaseWithRelease(pointerRelease)

	assert.Equal(t, expectedError, actualError)
}

func TestGetNameForRelease(t *testing.T) {
	expectedName := "dummyRelease"
	dummyRelease := &githubRelease{createDummyGithubRelease(releaseDescriptor{name: expectedName})}

	actualName := dummyRelease.GetName()

	assert.Equal(t, expectedName, actualName)
}
func TestGetNameForAsset(t *testing.T) {
	asset := createDummyAsset()

	actualName := asset.GetName()

	assert.Equal(t, dummyAssetName, actualName)
}

func TestGetDownloadUrl(t *testing.T) {
	asset := createDummyAsset()

	actualUrl := asset.GetDownloadUrl()

	assert.Equal(t, dummyAssetURL, actualUrl)
}

func TestPagerToStart(t *testing.T) {
	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo)
	dummyGithubRelease := createDummyGithubRelease(releaseDescriptor{name: "dummyRelease"})

	opt := &github.ListOptions{Page: 1}
	expectedReleasesOnPage := []*github.RepositoryRelease{dummyGithubRelease}
	expectedResponse := &github.Response{}

	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, opt).Return(expectedReleasesOnPage, expectedResponse, nil)

	err := pager.toStart()

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(expectedReleasesOnPage, pager.(*releasePager).currentReleases)
	assert.Equal(expectedResponse, pager.(*releasePager).currentResponse)
}

func TestPagerToStartListReleaseFail(t *testing.T) {
	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo)
	opt := &github.ListOptions{Page: 1}
	expectedError := errors.New("fail to Start")
	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, opt).Return(nil, nil, expectedError)

	err := pager.toStart()

	assert := assert.New(t)
	assert.Equal(expectedError, err)
}

func TestHasMore(t *testing.T) {
	givenPagerHasNotBeenInitialisedThenShouldPanic(t)
	givenThereIsMoreThenTrueShouldBeReturned(t)
	givenThereIsNoMoreThenFalseShouldBeReturn(t)
}

func givenPagerHasNotBeenInitialisedThenShouldPanic(t *testing.T) {
	repo, _, _ := getRepo(t)
	pager := newReleasePager(repo)

	as := assert.New(t)
	as.Panics(func() { pager.hasMore() })
}

func givenThereIsMoreThenTrueShouldBeReturned(t *testing.T) {
	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo)

	opt := &github.ListOptions{Page: 1}
	expectedResponse := &github.Response{
		NextPage: 1,
	}
	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, opt).Return(nil, expectedResponse, nil)

	pager.toStart()
	actualBool := pager.hasMore()

	assert := assert.New(t)
	assert.True(actualBool)
}

func givenThereIsNoMoreThenFalseShouldBeReturn(t *testing.T) {
	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo)

	opt := &github.ListOptions{Page: 1}
	expectedResponse := &github.Response{
		NextPage: 0,
	}
	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, opt).Return(nil, expectedResponse, nil)

	pager.toStart()
	actualBool := pager.hasMore()

	assert := assert.New(t)
	assert.False(actualBool)
}

func TestGetNextPage(t *testing.T) {
	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	expectedNextPage := 3
	pager.opt = &github.ListOptions{Page: 1}
	pager.currentResponse = &github.Response{
		NextPage: expectedNextPage,
	}

	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, &github.ListOptions{Page: expectedNextPage}).Return(nil, nil, nil)

	err := pager.getNextPage()

	assert := assert.New(t)
	assert.Nil(err)
}

func TestMoveToPageWhenPointerReleaseMightBeOnFirstPage(t *testing.T) {
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-02T22:08:41+00:00"})
	olderRelease := createDummyGithubRelease(releaseDescriptor{name: "olderRelease", date: "2012-11-01T22:08:41+00:00"})

	repo, _, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	pager.currentReleases = []*github.RepositoryRelease{olderRelease}

	err := pager.moveToPageContainingRelease(pointerRelease)

	assert.Nil(t, err)
	assert.Equal(t, olderRelease, pager.currentReleases[0])
}

func TestMoveToPageWhenPointerReleaseMightBeOnOtherPageThanFirstPage(t *testing.T) {
	olderRelease := createDummyGithubRelease(releaseDescriptor{name: "olderRelease", date: "2012-11-01T22:08:41+00:00"})
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-02T22:08:41+00:00"})
	youngerRelease := createDummyGithubRelease(releaseDescriptor{name: "youngerRelease", date: "2012-11-03T22:08:41+00:00"})

	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	pager.currentReleases = []*github.RepositoryRelease{youngerRelease}
	pager.currentResponse = &github.Response{
		NextPage: 2,
	}
	pager.opt = &github.ListOptions{Page: 1}
	nextPageReleases := []*github.RepositoryRelease{olderRelease}

	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, gomock.Any()).Return(nextPageReleases, nil, nil)

	err := pager.moveToPageContainingRelease(pointerRelease)

	assert.Nil(t, err)
	assert.Equal(t, olderRelease, pager.currentReleases[0])
}

func TestMoveToPageWhenPointerReleaseIsTheLastOnAPage(t *testing.T) {
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-02T22:08:41+00:00"})

	repo, _, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)

	pager.currentReleases = []*github.RepositoryRelease{pointerRelease}

	err := pager.moveToPageContainingRelease(pointerRelease)

	assert.Nil(t, err)
	assert.Equal(t, pointerRelease, pager.currentReleases[0])
}

func TestMoveToPageWhenPointerReleaseCannotBeInThePages(t *testing.T) {
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-02T22:08:41+00:00"})
	youngerRelease := createDummyGithubRelease(releaseDescriptor{name: "youngerRelease", date: "2012-11-03T22:08:41+00:00"})

	repo, _, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	pager.currentResponse = &github.Response{
		NextPage: 0,
	}
	pager.currentReleases = []*github.RepositoryRelease{youngerRelease}
	pager.opt = &github.ListOptions{Page: 1}
	expectedError := fmt.Errorf("The release %s version has not been found", *pointerRelease.Name)

	actualError := pager.moveToPageContainingRelease(pointerRelease)

	assert.Equal(t, expectedError, actualError)
}

func TestMoveToPageContainingReleaseWhenHasNextReturnedAnError(t *testing.T) {
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-02T22:08:41+00:00"})
	youngerRelease := createDummyGithubRelease(releaseDescriptor{name: "youngerRelease", date: "2012-11-03T22:08:41+00:00"})
	youngestRelease := createDummyGithubRelease(releaseDescriptor{name: "olderRelease", date: "2012-11-04T22:08:41+00:00"})

	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	pager.currentReleases = []*github.RepositoryRelease{youngestRelease}
	pager.currentResponse = &github.Response{
		NextPage: 2,
	}
	pager.opt = &github.ListOptions{Page: 1}

	expectedError := errors.New("Error on hasNext")

	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, gomock.Any()).Return([]*github.RepositoryRelease{youngerRelease}, nil, expectedError)

	actualError := pager.moveToPageContainingRelease(pointerRelease)

	assert.Equal(t, expectedError, actualError)
}

func TestFindReleaseOnPageWhenPointerReleaseIsLast(t *testing.T) {
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-01T22:08:41+00:00"})
	youngerRelease := createDummyGithubRelease(releaseDescriptor{name: "yougerRelease", date: "2012-11-02T22:08:41+00:00"})

	repo, _, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	pager.currentReleases = []*github.RepositoryRelease{youngerRelease, pointerRelease}

	index, err := pager.findReleaseOnPage(pointerRelease)

	assert.Nil(t, err)
	assert.Equal(t, 1, index)
}

func TestFindReleaseOnPageWhenPointerReleaseIsNotFirstNorLast(t *testing.T) {
	olderRelease := createDummyGithubRelease(releaseDescriptor{name: "olderRelease", date: "2012-11-01T22:08:41+00:00"})
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-02T22:08:41+00:00"})
	youngerRelease := createDummyGithubRelease(releaseDescriptor{name: "youngerRelease", date: "2012-11-03T22:08:41+00:00"})

	repo, _, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	pager.currentReleases = []*github.RepositoryRelease{youngerRelease, pointerRelease, olderRelease}

	index, err := pager.findReleaseOnPage(pointerRelease)

	assert.Nil(t, err)
	assert.Equal(t, 1, index)
}

func TestFindReleaseOnPageWhenPointerReleaseIsFirst(t *testing.T) {
	oldestRelease := createDummyGithubRelease(releaseDescriptor{name: "oldestRelease", date: "2012-11-01T22:08:41+00:00"})
	olderRelease := createDummyGithubRelease(releaseDescriptor{name: "olderRelease", date: "2012-11-02T22:08:41+00:00"})
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-03T22:08:41+00:00"})

	repo, _, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	pager.currentReleases = []*github.RepositoryRelease{pointerRelease, olderRelease, oldestRelease}

	index, err := pager.findReleaseOnPage(pointerRelease)

	assert.Nil(t, err)
	assert.Equal(t, 0, index)
}
func TestFindReleaseOnPageWhenPointerReleaseIsNotPresent(t *testing.T) {
	otherRelease := createDummyGithubRelease(releaseDescriptor{name: "otherRelease", date: "2012-11-01T22:08:41+00:00"})
	otherRelease1 := createDummyGithubRelease(releaseDescriptor{name: "otherRelease1", date: "2012-11-03T22:08:41+00:00"})
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-02T22:08:41+00:00"})

	expectedError := fmt.Errorf("The version %s version has not been found", *pointerRelease.Name)

	repo, _, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	pager.currentReleases = []*github.RepositoryRelease{otherRelease1, otherRelease}

	index, actualError := pager.findReleaseOnPage(pointerRelease)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, -1, index)
}

func TestLocateRelease(t *testing.T) {
	dummyRelease1 := createDummyGithubRelease(releaseDescriptor{name: "first", date: "2012-11-01T22:08:41+00:00"})
	dummyRelease2 := createDummyGithubRelease(releaseDescriptor{name: "second", date: "2012-11-02T22:08:41+00:00"})

	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	pager.currentReleases = []*github.RepositoryRelease{dummyRelease1}

	currentReleases := []*github.RepositoryRelease{dummyRelease2}

	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, gomock.Any()).Return(currentReleases, nil, nil)

	index, err := pager.locateRelease(dummyRelease2)

	assert.Nil(t, err)
	assert.Equal(t, 0, index)
}

func TestLocateReleaseToStartFailed(t *testing.T) {
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-02T22:08:41+00:00"})
	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)
	expectedError := errors.New("Pager failed to start")

	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, gomock.Any()).Return(nil, nil, expectedError)

	_, actualError := pager.locateRelease(pointerRelease)

	assert.Equal(t, expectedError, actualError)
}

func TestLocateReleasemoveToPageContainingReleaseFailed(t *testing.T) {
	pointerRelease := createDummyGithubRelease(releaseDescriptor{name: "pointerRelease", date: "2012-11-02T22:08:41+00:00"})
	youngerRelease := createDummyGithubRelease(releaseDescriptor{name: "youngerRelease", date: "2012-11-03T22:08:41+00:00"})

	repo, mockService, _ := getRepo(t)
	pager := newReleasePager(repo).(*releasePager)

	responseReceivedWithByStart := &github.Response{
		NextPage: 0,
	}
	releasesReceivedByToStart := []*github.RepositoryRelease{youngerRelease}

	expectedError := fmt.Errorf("The release %s version has not been found", *pointerRelease.Name)
	//toStart
	mockService.EXPECT().ListReleases(context.TODO(), repo.organisation, repo.repository, gomock.Any()).Return(releasesReceivedByToStart, responseReceivedWithByStart, nil)

	_, actualError := pager.locateRelease(pointerRelease)

	assert.Equal(t, expectedError, actualError)
}

func TestReleasesOnPage(t *testing.T) {
	pager := newReleasePager(nil).(*releasePager)
	dummyRelease := createDummyGithubRelease(releaseDescriptor{name: "dummyRelease"})
	expectedReleases := []*github.RepositoryRelease{dummyRelease}
	pager.currentReleases = expectedReleases

	actualReleases := pager.releasesOnPage()

	assert.Equal(t, expectedReleases, actualReleases)
}

func TestGetAssetByName(t *testing.T) {
	expectedName := "asset1"
	dummyRelease := githubRelease{createDummyGithubRelease(releaseDescriptor{name: "dummyRelease", assets: []string{expectedName, "asset2"}})}

	asset, err := dummyRelease.GetAssetByName(expectedName)

	assert.Nil(t, err)
	assert.Equal(t, expectedName, asset.GetName())
}

func TestGetAssetByNameAssetNotFound(t *testing.T) {
	lookingForAsset := "asset3"
	inRelease := "dummyRelease"
	expectedError := fmt.Errorf("asset %s not found in release %s", lookingForAsset, inRelease)
	dummyRelease := githubRelease{createDummyGithubRelease(releaseDescriptor{name: inRelease, assets: []string{"asset1", "asset2"}})}

	_, actualError := dummyRelease.GetAssetByName(lookingForAsset)

	assert.Equal(t, expectedError, actualError)
}

func getRepo(t *testing.T) (*githubRepository, *MockgithubRepositoryServiceInterface, *Mockpager) {
	ctrl := gomock.NewController(t)
	serviceMock := NewMockgithubRepositoryServiceInterface(ctrl)
	repo := githubRepository{
		service:      serviceMock,
		organisation: "hugo",
		repository:   "gohugo",
	}
	pagerMock := NewMockpager(ctrl)
	repo.pager = pagerMock
	return &repo, serviceMock, pagerMock
}

var dummyAssetName = "dummyAsset"
var dummyAssetURL = "dummyURL"

func createDummyAsset() *githubAsset {
	asset := &github.ReleaseAsset{
		Name:               &dummyAssetName,
		BrowserDownloadURL: &dummyAssetURL,
	}
	return &githubAsset{asset}
}

type releaseDescriptor struct {
	name   string
	date   string
	assets []string
}

func createDummyGithubRelease(descriptor releaseDescriptor) *github.RepositoryRelease {
	name := descriptor.name
	if name == "" {
		name = "dummyRelease"
	}
	date, _ := time.Parse(time.RFC3339, descriptor.date)
	assets := make([]*github.ReleaseAsset, len(descriptor.assets))
	for index, _ := range descriptor.assets {
		assets[index] = &github.ReleaseAsset{Name: &descriptor.assets[index]}
	}
	return &github.RepositoryRelease{Name: &name, CreatedAt: &github.Timestamp{date}, Assets: assets}
}
