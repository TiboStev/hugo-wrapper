// Code generated by MockGen. DO NOT EDIT.
// Source: repository_service.go

// Package versionmanager is a generated GoMock package.
package versionmanager

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/v31/github"
	reflect "reflect"
)

// MockRepositoryClient is a mock of RepositoryClient interface
type MockRepositoryClient struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryClientMockRecorder
}

// MockRepositoryClientMockRecorder is the mock recorder for MockRepositoryClient
type MockRepositoryClientMockRecorder struct {
	mock *MockRepositoryClient
}

// NewMockRepositoryClient creates a new mock instance
func NewMockRepositoryClient(ctrl *gomock.Controller) *MockRepositoryClient {
	mock := &MockRepositoryClient{ctrl: ctrl}
	mock.recorder = &MockRepositoryClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepositoryClient) EXPECT() *MockRepositoryClientMockRecorder {
	return m.recorder
}

// GetLatestRelease mocks base method
func (m *MockRepositoryClient) GetLatestRelease() (Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestRelease")
	ret0, _ := ret[0].(Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestRelease indicates an expected call of GetLatestRelease
func (mr *MockRepositoryClientMockRecorder) GetLatestRelease() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestRelease", reflect.TypeOf((*MockRepositoryClient)(nil).GetLatestRelease))
}

// GetReleaseByTag mocks base method
func (m *MockRepositoryClient) GetReleaseByTag(tag string) (Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReleaseByTag", tag)
	ret0, _ := ret[0].(Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReleaseByTag indicates an expected call of GetReleaseByTag
func (mr *MockRepositoryClientMockRecorder) GetReleaseByTag(tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReleaseByTag", reflect.TypeOf((*MockRepositoryClient)(nil).GetReleaseByTag), tag)
}

// GetPreviousRelease mocks base method
func (m *MockRepositoryClient) GetPreviousRelease(tag string) (Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPreviousRelease", tag)
	ret0, _ := ret[0].(Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPreviousRelease indicates an expected call of GetPreviousRelease
func (mr *MockRepositoryClientMockRecorder) GetPreviousRelease(tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPreviousRelease", reflect.TypeOf((*MockRepositoryClient)(nil).GetPreviousRelease), tag)
}

// MockRelease is a mock of Release interface
type MockRelease struct {
	ctrl     *gomock.Controller
	recorder *MockReleaseMockRecorder
}

// MockReleaseMockRecorder is the mock recorder for MockRelease
type MockReleaseMockRecorder struct {
	mock *MockRelease
}

// NewMockRelease creates a new mock instance
func NewMockRelease(ctrl *gomock.Controller) *MockRelease {
	mock := &MockRelease{ctrl: ctrl}
	mock.recorder = &MockReleaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRelease) EXPECT() *MockReleaseMockRecorder {
	return m.recorder
}

// GetName mocks base method
func (m *MockRelease) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockReleaseMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockRelease)(nil).GetName))
}

// GetAssetByName mocks base method
func (m *MockRelease) GetAssetByName(name string) (Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssetByName", name)
	ret0, _ := ret[0].(Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssetByName indicates an expected call of GetAssetByName
func (mr *MockReleaseMockRecorder) GetAssetByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssetByName", reflect.TypeOf((*MockRelease)(nil).GetAssetByName), name)
}

// MockAsset is a mock of Asset interface
type MockAsset struct {
	ctrl     *gomock.Controller
	recorder *MockAssetMockRecorder
}

// MockAssetMockRecorder is the mock recorder for MockAsset
type MockAssetMockRecorder struct {
	mock *MockAsset
}

// NewMockAsset creates a new mock instance
func NewMockAsset(ctrl *gomock.Controller) *MockAsset {
	mock := &MockAsset{ctrl: ctrl}
	mock.recorder = &MockAssetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAsset) EXPECT() *MockAssetMockRecorder {
	return m.recorder
}

// GetName mocks base method
func (m *MockAsset) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockAssetMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockAsset)(nil).GetName))
}

// GetDownloadUrl mocks base method
func (m *MockAsset) GetDownloadUrl() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDownloadUrl")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDownloadUrl indicates an expected call of GetDownloadUrl
func (mr *MockAssetMockRecorder) GetDownloadUrl() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDownloadUrl", reflect.TypeOf((*MockAsset)(nil).GetDownloadUrl))
}

// MockgithubRepositoryInterface is a mock of githubRepositoryInterface interface
type MockgithubRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockgithubRepositoryInterfaceMockRecorder
}

// MockgithubRepositoryInterfaceMockRecorder is the mock recorder for MockgithubRepositoryInterface
type MockgithubRepositoryInterfaceMockRecorder struct {
	mock *MockgithubRepositoryInterface
}

// NewMockgithubRepositoryInterface creates a new mock instance
func NewMockgithubRepositoryInterface(ctrl *gomock.Controller) *MockgithubRepositoryInterface {
	mock := &MockgithubRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockgithubRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockgithubRepositoryInterface) EXPECT() *MockgithubRepositoryInterfaceMockRecorder {
	return m.recorder
}

// MockgithubRepositoryServiceInterface is a mock of githubRepositoryServiceInterface interface
type MockgithubRepositoryServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockgithubRepositoryServiceInterfaceMockRecorder
}

// MockgithubRepositoryServiceInterfaceMockRecorder is the mock recorder for MockgithubRepositoryServiceInterface
type MockgithubRepositoryServiceInterfaceMockRecorder struct {
	mock *MockgithubRepositoryServiceInterface
}

// NewMockgithubRepositoryServiceInterface creates a new mock instance
func NewMockgithubRepositoryServiceInterface(ctrl *gomock.Controller) *MockgithubRepositoryServiceInterface {
	mock := &MockgithubRepositoryServiceInterface{ctrl: ctrl}
	mock.recorder = &MockgithubRepositoryServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockgithubRepositoryServiceInterface) EXPECT() *MockgithubRepositoryServiceInterfaceMockRecorder {
	return m.recorder
}

// GetLatestRelease mocks base method
func (m *MockgithubRepositoryServiceInterface) GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestRelease", ctx, owner, repo)
	ret0, _ := ret[0].(*github.RepositoryRelease)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLatestRelease indicates an expected call of GetLatestRelease
func (mr *MockgithubRepositoryServiceInterfaceMockRecorder) GetLatestRelease(ctx, owner, repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestRelease", reflect.TypeOf((*MockgithubRepositoryServiceInterface)(nil).GetLatestRelease), ctx, owner, repo)
}

// GetReleaseByTag mocks base method
func (m *MockgithubRepositoryServiceInterface) GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReleaseByTag", ctx, owner, repo, tag)
	ret0, _ := ret[0].(*github.RepositoryRelease)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetReleaseByTag indicates an expected call of GetReleaseByTag
func (mr *MockgithubRepositoryServiceInterfaceMockRecorder) GetReleaseByTag(ctx, owner, repo, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReleaseByTag", reflect.TypeOf((*MockgithubRepositoryServiceInterface)(nil).GetReleaseByTag), ctx, owner, repo, tag)
}

// ListReleases mocks base method
func (m *MockgithubRepositoryServiceInterface) ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListReleases", ctx, owner, repo, opts)
	ret0, _ := ret[0].([]*github.RepositoryRelease)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListReleases indicates an expected call of ListReleases
func (mr *MockgithubRepositoryServiceInterfaceMockRecorder) ListReleases(ctx, owner, repo, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListReleases", reflect.TypeOf((*MockgithubRepositoryServiceInterface)(nil).ListReleases), ctx, owner, repo, opts)
}