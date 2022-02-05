// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package rooms is a generated GoMock package.
package rooms

import (
	context "context"
	reflect "reflect"

	cheers "github.com/aljorhythm/yumseng/cheers"
	gomock "github.com/golang/mock/gomock"
)

// MockRoomServicer is a mock of RoomServicer interface.
type MockRoomServicer struct {
	ctrl     *gomock.Controller
	recorder *MockRoomServicerMockRecorder
}

// MockRoomServicerMockRecorder is the mock recorder for MockRoomServicer.
type MockRoomServicerMockRecorder struct {
	mock *MockRoomServicer
}

// NewMockRoomServicer creates a new mock instance.
func NewMockRoomServicer(ctrl *gomock.Controller) *MockRoomServicer {
	mock := &MockRoomServicer{ctrl: ctrl}
	mock.recorder = &MockRoomServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoomServicer) EXPECT() *MockRoomServicerMockRecorder {
	return m.recorder
}

// AddCheer mocks base method.
func (m *MockRoomServicer) AddCheer(room *Room, cheer *cheers.Cheer, user User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCheer", room, cheer, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCheer indicates an expected call of AddCheer.
func (mr *MockRoomServicerMockRecorder) AddCheer(room, cheer, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCheer", reflect.TypeOf((*MockRoomServicer)(nil).AddCheer), room, cheer, user)
}

// AddCheerAddedListener mocks base method.
func (m *MockRoomServicer) AddCheerAddedListener(room *Room, user User, clientId string, callback Callback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCheerAddedListener", room, user, clientId, callback)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCheerAddedListener indicates an expected call of AddCheerAddedListener.
func (mr *MockRoomServicerMockRecorder) AddCheerAddedListener(room, user, clientId, callback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCheerAddedListener", reflect.TypeOf((*MockRoomServicer)(nil).AddCheerAddedListener), room, user, clientId, callback)
}

// AddCheerImage mocks base method.
func (m *MockRoomServicer) AddCheerImage(ctx context.Context, roomId string, user User, url string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCheerImage", ctx, roomId, user, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCheerImage indicates an expected call of AddCheerImage.
func (mr *MockRoomServicerMockRecorder) AddCheerImage(ctx, roomId, user, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCheerImage", reflect.TypeOf((*MockRoomServicer)(nil).AddCheerImage), ctx, roomId, user, url)
}

// GetCheerImages mocks base method.
func (m *MockRoomServicer) GetCheerImages(ctx context.Context, roomId string, user User) ([]*CheerImage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCheerImages", ctx, roomId, user)
	ret0, _ := ret[0].([]*CheerImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCheerImages indicates an expected call of GetCheerImages.
func (mr *MockRoomServicerMockRecorder) GetCheerImages(ctx, roomId, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCheerImages", reflect.TypeOf((*MockRoomServicer)(nil).GetCheerImages), ctx, roomId, user)
}

// GetLeaderboard mocks base method.
func (m *MockRoomServicer) GetLeaderboard(roomId string) []*UserInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLeaderboard", roomId)
	ret0, _ := ret[0].([]*UserInfo)
	return ret0
}

// GetLeaderboard indicates an expected call of GetLeaderboard.
func (mr *MockRoomServicerMockRecorder) GetLeaderboard(roomId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLeaderboard", reflect.TypeOf((*MockRoomServicer)(nil).GetLeaderboard), roomId)
}

// GetOrCreateRoom mocks base method.
func (m *MockRoomServicer) GetOrCreateRoom(name string) *Room {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrCreateRoom", name)
	ret0, _ := ret[0].(*Room)
	return ret0
}

// GetOrCreateRoom indicates an expected call of GetOrCreateRoom.
func (mr *MockRoomServicerMockRecorder) GetOrCreateRoom(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrCreateRoom", reflect.TypeOf((*MockRoomServicer)(nil).GetOrCreateRoom), name)
}

// GetRoom mocks base method.
func (m *MockRoomServicer) GetRoom(name string) *Room {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoom", name)
	ret0, _ := ret[0].(*Room)
	return ret0
}

// GetRoom indicates an expected call of GetRoom.
func (mr *MockRoomServicerMockRecorder) GetRoom(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoom", reflect.TypeOf((*MockRoomServicer)(nil).GetRoom), name)
}

// RemoveOutdatedCheers mocks base method.
func (m *MockRoomServicer) RemoveOutdatedCheers() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveOutdatedCheers")
}

// RemoveOutdatedCheers indicates an expected call of RemoveOutdatedCheers.
func (mr *MockRoomServicerMockRecorder) RemoveOutdatedCheers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOutdatedCheers", reflect.TypeOf((*MockRoomServicer)(nil).RemoveOutdatedCheers))
}

// StopListeningCheers mocks base method.
func (m *MockRoomServicer) StopListeningCheers(room *Room, clientId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopListeningCheers", room, clientId)
}

// StopListeningCheers indicates an expected call of StopListeningCheers.
func (mr *MockRoomServicerMockRecorder) StopListeningCheers(room, clientId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopListeningCheers", reflect.TypeOf((*MockRoomServicer)(nil).StopListeningCheers), room, clientId)
}

// UploadCheerImage mocks base method.
func (m *MockRoomServicer) UploadCheerImage(ctx context.Context, roomId string, user User, data []byte) (*CheerImage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadCheerImage", ctx, roomId, user, data)
	ret0, _ := ret[0].(*CheerImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadCheerImage indicates an expected call of UploadCheerImage.
func (mr *MockRoomServicerMockRecorder) UploadCheerImage(ctx, roomId, user, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadCheerImage", reflect.TypeOf((*MockRoomServicer)(nil).UploadCheerImage), ctx, roomId, user, data)
}

// UserJoinsRoom mocks base method.
func (m *MockRoomServicer) UserJoinsRoom(ctx context.Context, room *Room, user User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserJoinsRoom", ctx, room, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserJoinsRoom indicates an expected call of UserJoinsRoom.
func (mr *MockRoomServicerMockRecorder) UserJoinsRoom(ctx, room, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserJoinsRoom", reflect.TypeOf((*MockRoomServicer)(nil).UserJoinsRoom), ctx, room, user)
}
