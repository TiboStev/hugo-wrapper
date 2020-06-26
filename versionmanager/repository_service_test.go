package versionmanager

import (
	context "context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	github "github.com/google/go-github/v31/github"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestRelease(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	serviceMock := NewMockgithubRepositoryServiceInterface(ctrl)

	repo := githubRepository{
		service:      serviceMock,
		organisation: "hugo",
		repository:   "gohugo",
	}
	serviceMock.EXPECT().GetLatestRelease(context.TODO(), repo.organisation, repo.repository).Return(new(github.RepositoryRelease), nil, nil)
	repo.GetLatestRelease()
}

func TestGetReleaseByTag(t *testing.T) {
	getReleaseByTagHappyFlow(t)
	getReleaseByTagError(t)
}

func getReleaseByTagHappyFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	tag := "tag"

	defer ctrl.Finish()
	serviceMock := NewMockgithubRepositoryServiceInterface(ctrl)

	repo := githubRepository{
		service:      serviceMock,
		organisation: "hugo",
		repository:   "gohugo",
	}
	serviceMock.EXPECT().GetReleaseByTag(context.TODO(), repo.organisation, repo.repository, tag)
	repo.GetReleaseByTag(tag)
}

func getReleaseByTagError(t *testing.T) {
	ctrl := gomock.NewController(t)
	tag := "tag"

	expectedError := errors.New("expected error has occured")

	defer ctrl.Finish()
	serviceMock := NewMockgithubRepositoryServiceInterface(ctrl)

	repo := githubRepository{
		service:      serviceMock,
		organisation: "hugo",
		repository:   "gohugo",
	}
	serviceMock.EXPECT().GetReleaseByTag(context.TODO(), repo.organisation, repo.repository, tag).Return(nil, nil, expectedError)

	_, err := repo.GetReleaseByTag(tag)
	assert := assert.New(t)
	assert.Equal(expectedError, err, "error returned by githubrepo should be passed as it is")
}

func TestGetPreviousRelease(t *testing.T) {
	whenReleaseWithTagReturnError_thenTheErrorShouldBeReturned(t)
	whenNextPageCantBeReached()
	whenErrorOnPage()
	whenTagIsTheLastReleaseOfThePageAndThereIsNoMorePages()
	whenPreviousReleaseIsOnTheNextPage()
	whenTagIsOnThePage()
}

func whenReleaseWithTagReturnError_thenTheErrorShouldBeReturned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serviceMock := NewMockgithubRepositoryServiceInterface(ctrl)
	repo := githubRepository{
		service:      serviceMock,
		organisation: "hugo",
		repository:   "gohugo",
	}
	pager, _ := repo.newReleasePager()
	fmt.Printf("%+v\n", pager)
}

func whenNextPageCantBeReached()                             {}
func whenErrorOnPage()                                       {}
func whenTagIsTheLastReleaseOfThePageAndThereIsNoMorePages() {}
func whenPreviousReleaseIsOnTheNextPage()                    {}
func whenTagIsOnThePage()                                    {}

func TestGetName_forRelease(t *testing.T) {}
func TestGetName_forAsset(t *testing.T)   {}

func TestGetDownloadUrl(t *testing.T) {}

func TesthasMore(t *testing.T) {}

func TesttoStart(t *testing.T) {}

func TestgetNextPage(t *testing.T) {}

func TestmoveToPageContainingRelease(t *testing.T) {}

func TestfindReleaseOnPage(t *testing.T) {}
