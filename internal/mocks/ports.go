// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/ports/ports.go

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/posilva/simplechat/internal/core/domain"
	ports "github.com/posilva/simplechat/internal/core/ports"
)

// MockChatService is a mock of ChatService interface.
type MockChatService struct {
	ctrl     *gomock.Controller
	recorder *MockChatServiceMockRecorder
}

// MockChatServiceMockRecorder is the mock recorder for MockChatService.
type MockChatServiceMockRecorder struct {
	mock *MockChatService
}

// NewMockChatService creates a new mock instance.
func NewMockChatService(ctrl *gomock.Controller) *MockChatService {
	mock := &MockChatService{ctrl: ctrl}
	mock.recorder = &MockChatServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatService) EXPECT() *MockChatServiceMockRecorder {
	return m.recorder
}

// DeRegister mocks base method.
func (m *MockChatService) DeRegister(ep ports.Endpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeRegister", ep)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeRegister indicates an expected call of DeRegister.
func (mr *MockChatServiceMockRecorder) DeRegister(ep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeRegister", reflect.TypeOf((*MockChatService)(nil).DeRegister), ep)
}

// History mocks base method.
func (m *MockChatService) History(dst string, since time.Duration) ([]*domain.ModeratedMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "History", dst, since)
	ret0, _ := ret[0].([]*domain.ModeratedMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// History indicates an expected call of History.
func (mr *MockChatServiceMockRecorder) History(dst, since interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "History", reflect.TypeOf((*MockChatService)(nil).History), dst, since)
}

// Register mocks base method.
func (m *MockChatService) Register(ep ports.Endpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ep)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockChatServiceMockRecorder) Register(ep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockChatService)(nil).Register), ep)
}

// Send mocks base method.
func (m *MockChatService) Send(arg0 domain.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockChatServiceMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockChatService)(nil).Send), arg0)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Fetch mocks base method.
func (m *MockRepository) Fetch(key string, since time.Duration) ([]*domain.ModeratedMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", key, since)
	ret0, _ := ret[0].([]*domain.ModeratedMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockRepositoryMockRecorder) Fetch(key, since interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockRepository)(nil).Fetch), key, since)
}

// Store mocks base method.
func (m_2 *MockRepository) Store(m domain.ModeratedMessage) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Store", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockRepositoryMockRecorder) Store(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockRepository)(nil).Store), m)
}

// MockNotifier is a mock of Notifier interface.
type MockNotifier struct {
	ctrl     *gomock.Controller
	recorder *MockNotifierMockRecorder
}

// MockNotifierMockRecorder is the mock recorder for MockNotifier.
type MockNotifierMockRecorder struct {
	mock *MockNotifier
}

// NewMockNotifier creates a new mock instance.
func NewMockNotifier(ctrl *gomock.Controller) *MockNotifier {
	mock := &MockNotifier{ctrl: ctrl}
	mock.recorder = &MockNotifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotifier) EXPECT() *MockNotifierMockRecorder {
	return m.recorder
}

// Broadcast mocks base method.
func (m_2 *MockNotifier) Broadcast(m domain.ModeratedMessage) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Broadcast", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Broadcast indicates an expected call of Broadcast.
func (mr *MockNotifierMockRecorder) Broadcast(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*MockNotifier)(nil).Broadcast), m)
}

// Subscribe mocks base method.
func (m *MockNotifier) Subscribe(ep ports.Endpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ep)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockNotifierMockRecorder) Subscribe(ep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockNotifier)(nil).Subscribe), ep)
}

// Unsubscribe mocks base method.
func (m *MockNotifier) Unsubscribe(ep ports.Endpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", ep)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockNotifierMockRecorder) Unsubscribe(ep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockNotifier)(nil).Unsubscribe), ep)
}

// MockModerator is a mock of Moderator interface.
type MockModerator struct {
	ctrl     *gomock.Controller
	recorder *MockModeratorMockRecorder
}

// MockModeratorMockRecorder is the mock recorder for MockModerator.
type MockModeratorMockRecorder struct {
	mock *MockModerator
}

// NewMockModerator creates a new mock instance.
func NewMockModerator(ctrl *gomock.Controller) *MockModerator {
	mock := &MockModerator{ctrl: ctrl}
	mock.recorder = &MockModeratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModerator) EXPECT() *MockModeratorMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m_2 *MockModerator) Check(m domain.Message) (*domain.ModeratedMessage, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Check", m)
	ret0, _ := ret[0].(*domain.ModeratedMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockModeratorMockRecorder) Check(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockModerator)(nil).Check), m)
}

// MockRegistry is a mock of Registry interface.
type MockRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockRegistryMockRecorder
}

// MockRegistryMockRecorder is the mock recorder for MockRegistry.
type MockRegistryMockRecorder struct {
	mock *MockRegistry
}

// NewMockRegistry creates a new mock instance.
func NewMockRegistry(ctrl *gomock.Controller) *MockRegistry {
	mock := &MockRegistry{ctrl: ctrl}
	mock.recorder = &MockRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegistry) EXPECT() *MockRegistryMockRecorder {
	return m.recorder
}

// DeRegister mocks base method.
func (m *MockRegistry) DeRegister(ep ports.Endpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeRegister", ep)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeRegister indicates an expected call of DeRegister.
func (mr *MockRegistryMockRecorder) DeRegister(ep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeRegister", reflect.TypeOf((*MockRegistry)(nil).DeRegister), ep)
}

// Notify mocks base method.
func (m_2 *MockRegistry) Notify(m domain.ModeratedMessage) {
	m_2.ctrl.T.Helper()
	m_2.ctrl.Call(m_2, "Notify", m)
}

// Notify indicates an expected call of Notify.
func (mr *MockRegistryMockRecorder) Notify(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Notify", reflect.TypeOf((*MockRegistry)(nil).Notify), m)
}

// Register mocks base method.
func (m *MockRegistry) Register(ep ports.Endpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ep)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockRegistryMockRecorder) Register(ep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockRegistry)(nil).Register), ep)
}

// MockReceiver is a mock of Receiver interface.
type MockReceiver struct {
	ctrl     *gomock.Controller
	recorder *MockReceiverMockRecorder
}

// MockReceiverMockRecorder is the mock recorder for MockReceiver.
type MockReceiverMockRecorder struct {
	mock *MockReceiver
}

// NewMockReceiver creates a new mock instance.
func NewMockReceiver(ctrl *gomock.Controller) *MockReceiver {
	mock := &MockReceiver{ctrl: ctrl}
	mock.recorder = &MockReceiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReceiver) EXPECT() *MockReceiverMockRecorder {
	return m.recorder
}

// Receive mocks base method.
func (m *MockReceiver) Receive(arg0 domain.ModeratedMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Receive", arg0)
}

// Receive indicates an expected call of Receive.
func (mr *MockReceiverMockRecorder) Receive(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Receive", reflect.TypeOf((*MockReceiver)(nil).Receive), arg0)
}

// Recover mocks base method.
func (m *MockReceiver) Recover() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Recover")
}

// Recover indicates an expected call of Recover.
func (mr *MockReceiverMockRecorder) Recover() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recover", reflect.TypeOf((*MockReceiver)(nil).Recover))
}

// MockEndpoint is a mock of Endpoint interface.
type MockEndpoint struct {
	ctrl     *gomock.Controller
	recorder *MockEndpointMockRecorder
}

// MockEndpointMockRecorder is the mock recorder for MockEndpoint.
type MockEndpointMockRecorder struct {
	mock *MockEndpoint
}

// NewMockEndpoint creates a new mock instance.
func NewMockEndpoint(ctrl *gomock.Controller) *MockEndpoint {
	mock := &MockEndpoint{ctrl: ctrl}
	mock.recorder = &MockEndpointMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEndpoint) EXPECT() *MockEndpointMockRecorder {
	return m.recorder
}

// ID mocks base method.
func (m *MockEndpoint) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockEndpointMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockEndpoint)(nil).ID))
}

// Receive mocks base method.
func (m *MockEndpoint) Receive(arg0 domain.ModeratedMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Receive", arg0)
}

// Receive indicates an expected call of Receive.
func (mr *MockEndpointMockRecorder) Receive(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Receive", reflect.TypeOf((*MockEndpoint)(nil).Receive), arg0)
}

// Recover mocks base method.
func (m *MockEndpoint) Recover() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Recover")
}

// Recover indicates an expected call of Recover.
func (mr *MockEndpointMockRecorder) Recover() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recover", reflect.TypeOf((*MockEndpoint)(nil).Recover))
}

// Room mocks base method.
func (m *MockEndpoint) Room() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Room")
	ret0, _ := ret[0].(string)
	return ret0
}

// Room indicates an expected call of Room.
func (mr *MockEndpointMockRecorder) Room() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Room", reflect.TypeOf((*MockEndpoint)(nil).Room))
}

// MockPresence is a mock of Presence interface.
type MockPresence struct {
	ctrl     *gomock.Controller
	recorder *MockPresenceMockRecorder
}

// MockPresenceMockRecorder is the mock recorder for MockPresence.
type MockPresenceMockRecorder struct {
	mock *MockPresence
}

// NewMockPresence creates a new mock instance.
func NewMockPresence(ctrl *gomock.Controller) *MockPresence {
	mock := &MockPresence{ctrl: ctrl}
	mock.recorder = &MockPresenceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPresence) EXPECT() *MockPresenceMockRecorder {
	return m.recorder
}

// IsPresent mocks base method.
func (m *MockPresence) IsPresent(ep ports.Endpoint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPresent", ep)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPresent indicates an expected call of IsPresent.
func (mr *MockPresenceMockRecorder) IsPresent(ep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPresent", reflect.TypeOf((*MockPresence)(nil).IsPresent), ep)
}

// Join mocks base method.
func (m *MockPresence) Join(ep ports.Endpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Join", ep)
	ret0, _ := ret[0].(error)
	return ret0
}

// Join indicates an expected call of Join.
func (mr *MockPresenceMockRecorder) Join(ep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Join", reflect.TypeOf((*MockPresence)(nil).Join), ep)
}

// Leave mocks base method.
func (m *MockPresence) Leave(ep ports.Endpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Leave", ep)
	ret0, _ := ret[0].(error)
	return ret0
}

// Leave indicates an expected call of Leave.
func (mr *MockPresenceMockRecorder) Leave(ep interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Leave", reflect.TypeOf((*MockPresence)(nil).Leave), ep)
}

// Presents mocks base method.
func (m *MockPresence) Presents(room string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Presents", room)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Presents indicates an expected call of Presents.
func (mr *MockPresenceMockRecorder) Presents(room interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Presents", reflect.TypeOf((*MockPresence)(nil).Presents), room)
}
