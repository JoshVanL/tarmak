// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/tarmak/interfaces/interfaces.go

package mocks

import (
	logrus "github.com/Sirupsen/logrus"
	gomock "github.com/golang/mock/gomock"
	kv "github.com/jetstack-experimental/vault-unsealer/pkg/kv"
	. "github.com/jetstack/tarmak/pkg/tarmak/interfaces"
	net "net"
)

// MockContext is a mock of Context interface
type MockContext struct {
	ctrl     *gomock.Controller
	recorder *MockContextMockRecorder
}

// MockContextMockRecorder is the mock recorder for MockContext
type MockContextMockRecorder struct {
	mock *MockContext
}

// NewMockContext creates a new mock instance
func NewMockContext(ctrl *gomock.Controller) *MockContext {
	mock := &MockContext{ctrl: ctrl}
	mock.recorder = &MockContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockContext) EXPECT() *MockContextMockRecorder {
	return _m.recorder
}

// Variables mocks base method
func (_m *MockContext) Variables() map[string]interface{} {
	ret := _m.ctrl.Call(_m, "Variables")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Variables indicates an expected call of Variables
func (_mr *MockContextMockRecorder) Variables() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Variables")
}

// Environment mocks base method
func (_m *MockContext) Environment() Environment {
	ret := _m.ctrl.Call(_m, "Environment")
	ret0, _ := ret[0].(Environment)
	return ret0
}

// Environment indicates an expected call of Environment
func (_mr *MockContextMockRecorder) Environment() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Environment")
}

// Name mocks base method
func (_m *MockContext) Name() string {
	ret := _m.ctrl.Call(_m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (_mr *MockContextMockRecorder) Name() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Name")
}

// Validate mocks base method
func (_m *MockContext) Validate() error {
	ret := _m.ctrl.Call(_m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate
func (_mr *MockContextMockRecorder) Validate() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Validate")
}

// Stacks mocks base method
func (_m *MockContext) Stacks() []Stack {
	ret := _m.ctrl.Call(_m, "Stacks")
	ret0, _ := ret[0].([]Stack)
	return ret0
}

// Stacks indicates an expected call of Stacks
func (_mr *MockContextMockRecorder) Stacks() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stacks")
}

// NetworkCIDR mocks base method
func (_m *MockContext) NetworkCIDR() *net.IPNet {
	ret := _m.ctrl.Call(_m, "NetworkCIDR")
	ret0, _ := ret[0].(*net.IPNet)
	return ret0
}

// NetworkCIDR indicates an expected call of NetworkCIDR
func (_mr *MockContextMockRecorder) NetworkCIDR() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NetworkCIDR")
}

// RemoteState mocks base method
func (_m *MockContext) RemoteState(stackName string) string {
	ret := _m.ctrl.Call(_m, "RemoteState", stackName)
	ret0, _ := ret[0].(string)
	return ret0
}

// RemoteState indicates an expected call of RemoteState
func (_mr *MockContextMockRecorder) RemoteState(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoteState", arg0)
}

// ConfigPath mocks base method
func (_m *MockContext) ConfigPath() string {
	ret := _m.ctrl.Call(_m, "ConfigPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// ConfigPath indicates an expected call of ConfigPath
func (_mr *MockContextMockRecorder) ConfigPath() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ConfigPath")
}

// BaseImage mocks base method
func (_m *MockContext) BaseImage() string {
	ret := _m.ctrl.Call(_m, "BaseImage")
	ret0, _ := ret[0].(string)
	return ret0
}

// BaseImage indicates an expected call of BaseImage
func (_mr *MockContextMockRecorder) BaseImage() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "BaseImage")
}

// SSHConfigPath mocks base method
func (_m *MockContext) SSHConfigPath() string {
	ret := _m.ctrl.Call(_m, "SSHConfigPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// SSHConfigPath indicates an expected call of SSHConfigPath
func (_mr *MockContextMockRecorder) SSHConfigPath() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SSHConfigPath")
}

// SSHHostKeysPath mocks base method
func (_m *MockContext) SSHHostKeysPath() string {
	ret := _m.ctrl.Call(_m, "SSHHostKeysPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// SSHHostKeysPath indicates an expected call of SSHHostKeysPath
func (_mr *MockContextMockRecorder) SSHHostKeysPath() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SSHHostKeysPath")
}

// SetImageID mocks base method
func (_m *MockContext) SetImageID(_param0 string) {
	_m.ctrl.Call(_m, "SetImageID", _param0)
}

// SetImageID indicates an expected call of SetImageID
func (_mr *MockContextMockRecorder) SetImageID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetImageID", arg0)
}

// ContextName mocks base method
func (_m *MockContext) ContextName() string {
	ret := _m.ctrl.Call(_m, "ContextName")
	ret0, _ := ret[0].(string)
	return ret0
}

// ContextName indicates an expected call of ContextName
func (_mr *MockContextMockRecorder) ContextName() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ContextName")
}

// Log mocks base method
func (_m *MockContext) Log() *logrus.Entry {
	ret := _m.ctrl.Call(_m, "Log")
	ret0, _ := ret[0].(*logrus.Entry)
	return ret0
}

// Log indicates an expected call of Log
func (_mr *MockContextMockRecorder) Log() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Log")
}

// MockEnvironment is a mock of Environment interface
type MockEnvironment struct {
	ctrl     *gomock.Controller
	recorder *MockEnvironmentMockRecorder
}

// MockEnvironmentMockRecorder is the mock recorder for MockEnvironment
type MockEnvironmentMockRecorder struct {
	mock *MockEnvironment
}

// NewMockEnvironment creates a new mock instance
func NewMockEnvironment(ctrl *gomock.Controller) *MockEnvironment {
	mock := &MockEnvironment{ctrl: ctrl}
	mock.recorder = &MockEnvironmentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockEnvironment) EXPECT() *MockEnvironmentMockRecorder {
	return _m.recorder
}

// Tarmak mocks base method
func (_m *MockEnvironment) Tarmak() Tarmak {
	ret := _m.ctrl.Call(_m, "Tarmak")
	ret0, _ := ret[0].(Tarmak)
	return ret0
}

// Tarmak indicates an expected call of Tarmak
func (_mr *MockEnvironmentMockRecorder) Tarmak() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Tarmak")
}

// Variables mocks base method
func (_m *MockEnvironment) Variables() map[string]interface{} {
	ret := _m.ctrl.Call(_m, "Variables")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Variables indicates an expected call of Variables
func (_mr *MockEnvironmentMockRecorder) Variables() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Variables")
}

// Provider mocks base method
func (_m *MockEnvironment) Provider() Provider {
	ret := _m.ctrl.Call(_m, "Provider")
	ret0, _ := ret[0].(Provider)
	return ret0
}

// Provider indicates an expected call of Provider
func (_mr *MockEnvironmentMockRecorder) Provider() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Provider")
}

// Validate mocks base method
func (_m *MockEnvironment) Validate() error {
	ret := _m.ctrl.Call(_m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate
func (_mr *MockEnvironmentMockRecorder) Validate() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Validate")
}

// Name mocks base method
func (_m *MockEnvironment) Name() string {
	ret := _m.ctrl.Call(_m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (_mr *MockEnvironmentMockRecorder) Name() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Name")
}

// BucketPrefix mocks base method
func (_m *MockEnvironment) BucketPrefix() string {
	ret := _m.ctrl.Call(_m, "BucketPrefix")
	ret0, _ := ret[0].(string)
	return ret0
}

// BucketPrefix indicates an expected call of BucketPrefix
func (_mr *MockEnvironmentMockRecorder) BucketPrefix() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "BucketPrefix")
}

// Contexts mocks base method
func (_m *MockEnvironment) Contexts() []Context {
	ret := _m.ctrl.Call(_m, "Contexts")
	ret0, _ := ret[0].([]Context)
	return ret0
}

// Contexts indicates an expected call of Contexts
func (_mr *MockEnvironmentMockRecorder) Contexts() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Contexts")
}

// SSHPrivateKeyPath mocks base method
func (_m *MockEnvironment) SSHPrivateKeyPath() string {
	ret := _m.ctrl.Call(_m, "SSHPrivateKeyPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// SSHPrivateKeyPath indicates an expected call of SSHPrivateKeyPath
func (_mr *MockEnvironmentMockRecorder) SSHPrivateKeyPath() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SSHPrivateKeyPath")
}

// SSHPrivateKey mocks base method
func (_m *MockEnvironment) SSHPrivateKey() interface{} {
	ret := _m.ctrl.Call(_m, "SSHPrivateKey")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// SSHPrivateKey indicates an expected call of SSHPrivateKey
func (_mr *MockEnvironmentMockRecorder) SSHPrivateKey() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SSHPrivateKey")
}

// Log mocks base method
func (_m *MockEnvironment) Log() *logrus.Entry {
	ret := _m.ctrl.Call(_m, "Log")
	ret0, _ := ret[0].(*logrus.Entry)
	return ret0
}

// Log indicates an expected call of Log
func (_mr *MockEnvironmentMockRecorder) Log() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Log")
}

// StateStack mocks base method
func (_m *MockEnvironment) StateStack() Stack {
	ret := _m.ctrl.Call(_m, "StateStack")
	ret0, _ := ret[0].(Stack)
	return ret0
}

// StateStack indicates an expected call of StateStack
func (_mr *MockEnvironmentMockRecorder) StateStack() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "StateStack")
}

// VaultStack mocks base method
func (_m *MockEnvironment) VaultStack() Stack {
	ret := _m.ctrl.Call(_m, "VaultStack")
	ret0, _ := ret[0].(Stack)
	return ret0
}

// VaultStack indicates an expected call of VaultStack
func (_mr *MockEnvironmentMockRecorder) VaultStack() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VaultStack")
}

// MockProvider is a mock of Provider interface
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return _m.recorder
}

// Name mocks base method
func (_m *MockProvider) Name() string {
	ret := _m.ctrl.Call(_m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (_mr *MockProviderMockRecorder) Name() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Name")
}

// Region mocks base method
func (_m *MockProvider) Region() string {
	ret := _m.ctrl.Call(_m, "Region")
	ret0, _ := ret[0].(string)
	return ret0
}

// Region indicates an expected call of Region
func (_mr *MockProviderMockRecorder) Region() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Region")
}

// Validate mocks base method
func (_m *MockProvider) Validate() error {
	ret := _m.ctrl.Call(_m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate
func (_mr *MockProviderMockRecorder) Validate() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Validate")
}

// RemoteStateBucketName mocks base method
func (_m *MockProvider) RemoteStateBucketName() string {
	ret := _m.ctrl.Call(_m, "RemoteStateBucketName")
	ret0, _ := ret[0].(string)
	return ret0
}

// RemoteStateBucketName indicates an expected call of RemoteStateBucketName
func (_mr *MockProviderMockRecorder) RemoteStateBucketName() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoteStateBucketName")
}

// RemoteStateBucketAvailable mocks base method
func (_m *MockProvider) RemoteStateBucketAvailable() (bool, error) {
	ret := _m.ctrl.Call(_m, "RemoteStateBucketAvailable")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoteStateBucketAvailable indicates an expected call of RemoteStateBucketAvailable
func (_mr *MockProviderMockRecorder) RemoteStateBucketAvailable() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoteStateBucketAvailable")
}

// RemoteState mocks base method
func (_m *MockProvider) RemoteState(contextName string, stackName string) string {
	ret := _m.ctrl.Call(_m, "RemoteState", contextName, stackName)
	ret0, _ := ret[0].(string)
	return ret0
}

// RemoteState indicates an expected call of RemoteState
func (_mr *MockProviderMockRecorder) RemoteState(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoteState", arg0, arg1)
}

// Environment mocks base method
func (_m *MockProvider) Environment() ([]string, error) {
	ret := _m.ctrl.Call(_m, "Environment")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Environment indicates an expected call of Environment
func (_mr *MockProviderMockRecorder) Environment() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Environment")
}

// Variables mocks base method
func (_m *MockProvider) Variables() map[string]interface{} {
	ret := _m.ctrl.Call(_m, "Variables")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Variables indicates an expected call of Variables
func (_mr *MockProviderMockRecorder) Variables() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Variables")
}

// QueryImage mocks base method
func (_m *MockProvider) QueryImage(tags map[string]string) (string, error) {
	ret := _m.ctrl.Call(_m, "QueryImage", tags)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryImage indicates an expected call of QueryImage
func (_mr *MockProviderMockRecorder) QueryImage(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "QueryImage", arg0)
}

// VaultKV mocks base method
func (_m *MockProvider) VaultKV() (kv.Service, error) {
	ret := _m.ctrl.Call(_m, "VaultKV")
	ret0, _ := ret[0].(kv.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VaultKV indicates an expected call of VaultKV
func (_mr *MockProviderMockRecorder) VaultKV() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VaultKV")
}

// ListHosts mocks base method
func (_m *MockProvider) ListHosts() ([]Host, error) {
	ret := _m.ctrl.Call(_m, "ListHosts")
	ret0, _ := ret[0].([]Host)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListHosts indicates an expected call of ListHosts
func (_mr *MockProviderMockRecorder) ListHosts() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListHosts")
}

// MockStack is a mock of Stack interface
type MockStack struct {
	ctrl     *gomock.Controller
	recorder *MockStackMockRecorder
}

// MockStackMockRecorder is the mock recorder for MockStack
type MockStackMockRecorder struct {
	mock *MockStack
}

// NewMockStack creates a new mock instance
func NewMockStack(ctrl *gomock.Controller) *MockStack {
	mock := &MockStack{ctrl: ctrl}
	mock.recorder = &MockStackMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockStack) EXPECT() *MockStackMockRecorder {
	return _m.recorder
}

// Variables mocks base method
func (_m *MockStack) Variables() map[string]interface{} {
	ret := _m.ctrl.Call(_m, "Variables")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Variables indicates an expected call of Variables
func (_mr *MockStackMockRecorder) Variables() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Variables")
}

// Name mocks base method
func (_m *MockStack) Name() string {
	ret := _m.ctrl.Call(_m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (_mr *MockStackMockRecorder) Name() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Name")
}

// Validate mocks base method
func (_m *MockStack) Validate() error {
	ret := _m.ctrl.Call(_m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate
func (_mr *MockStackMockRecorder) Validate() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Validate")
}

// Context mocks base method
func (_m *MockStack) Context() Context {
	ret := _m.ctrl.Call(_m, "Context")
	ret0, _ := ret[0].(Context)
	return ret0
}

// Context indicates an expected call of Context
func (_mr *MockStackMockRecorder) Context() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Context")
}

// RemoteState mocks base method
func (_m *MockStack) RemoteState() string {
	ret := _m.ctrl.Call(_m, "RemoteState")
	ret0, _ := ret[0].(string)
	return ret0
}

// RemoteState indicates an expected call of RemoteState
func (_mr *MockStackMockRecorder) RemoteState() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoteState")
}

// Log mocks base method
func (_m *MockStack) Log() *logrus.Entry {
	ret := _m.ctrl.Call(_m, "Log")
	ret0, _ := ret[0].(*logrus.Entry)
	return ret0
}

// Log indicates an expected call of Log
func (_mr *MockStackMockRecorder) Log() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Log")
}

// VerifyPost mocks base method
func (_m *MockStack) VerifyPost() error {
	ret := _m.ctrl.Call(_m, "VerifyPost")
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyPost indicates an expected call of VerifyPost
func (_mr *MockStackMockRecorder) VerifyPost() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VerifyPost")
}

// SetOutput mocks base method
func (_m *MockStack) SetOutput(_param0 map[string]interface{}) {
	_m.ctrl.Call(_m, "SetOutput", _param0)
}

// SetOutput indicates an expected call of SetOutput
func (_mr *MockStackMockRecorder) SetOutput(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetOutput", arg0)
}

// Output mocks base method
func (_m *MockStack) Output() map[string]interface{} {
	ret := _m.ctrl.Call(_m, "Output")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Output indicates an expected call of Output
func (_mr *MockStackMockRecorder) Output() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Output")
}

// MockTarmak is a mock of Tarmak interface
type MockTarmak struct {
	ctrl     *gomock.Controller
	recorder *MockTarmakMockRecorder
}

// MockTarmakMockRecorder is the mock recorder for MockTarmak
type MockTarmakMockRecorder struct {
	mock *MockTarmak
}

// NewMockTarmak creates a new mock instance
func NewMockTarmak(ctrl *gomock.Controller) *MockTarmak {
	mock := &MockTarmak{ctrl: ctrl}
	mock.recorder = &MockTarmakMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockTarmak) EXPECT() *MockTarmakMockRecorder {
	return _m.recorder
}

// Variables mocks base method
func (_m *MockTarmak) Variables() map[string]interface{} {
	ret := _m.ctrl.Call(_m, "Variables")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Variables indicates an expected call of Variables
func (_mr *MockTarmakMockRecorder) Variables() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Variables")
}

// Log mocks base method
func (_m *MockTarmak) Log() *logrus.Entry {
	ret := _m.ctrl.Call(_m, "Log")
	ret0, _ := ret[0].(*logrus.Entry)
	return ret0
}

// Log indicates an expected call of Log
func (_mr *MockTarmakMockRecorder) Log() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Log")
}

// RootPath mocks base method
func (_m *MockTarmak) RootPath() string {
	ret := _m.ctrl.Call(_m, "RootPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// RootPath indicates an expected call of RootPath
func (_mr *MockTarmakMockRecorder) RootPath() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RootPath")
}

// ConfigPath mocks base method
func (_m *MockTarmak) ConfigPath() string {
	ret := _m.ctrl.Call(_m, "ConfigPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// ConfigPath indicates an expected call of ConfigPath
func (_mr *MockTarmakMockRecorder) ConfigPath() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ConfigPath")
}

// Context mocks base method
func (_m *MockTarmak) Context() Context {
	ret := _m.ctrl.Call(_m, "Context")
	ret0, _ := ret[0].(Context)
	return ret0
}

// Context indicates an expected call of Context
func (_mr *MockTarmakMockRecorder) Context() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Context")
}

// Environments mocks base method
func (_m *MockTarmak) Environments() []Environment {
	ret := _m.ctrl.Call(_m, "Environments")
	ret0, _ := ret[0].([]Environment)
	return ret0
}

// Environments indicates an expected call of Environments
func (_mr *MockTarmakMockRecorder) Environments() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Environments")
}

// Terraform mocks base method
func (_m *MockTarmak) Terraform() Terraform {
	ret := _m.ctrl.Call(_m, "Terraform")
	ret0, _ := ret[0].(Terraform)
	return ret0
}

// Terraform indicates an expected call of Terraform
func (_mr *MockTarmakMockRecorder) Terraform() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Terraform")
}

// Packer mocks base method
func (_m *MockTarmak) Packer() Packer {
	ret := _m.ctrl.Call(_m, "Packer")
	ret0, _ := ret[0].(Packer)
	return ret0
}

// Packer indicates an expected call of Packer
func (_mr *MockTarmakMockRecorder) Packer() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Packer")
}

// SSH mocks base method
func (_m *MockTarmak) SSH() SSH {
	ret := _m.ctrl.Call(_m, "SSH")
	ret0, _ := ret[0].(SSH)
	return ret0
}

// SSH indicates an expected call of SSH
func (_mr *MockTarmakMockRecorder) SSH() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SSH")
}

// HomeDirExpand mocks base method
func (_m *MockTarmak) HomeDirExpand(in string) (string, error) {
	ret := _m.ctrl.Call(_m, "HomeDirExpand", in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HomeDirExpand indicates an expected call of HomeDirExpand
func (_mr *MockTarmakMockRecorder) HomeDirExpand(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "HomeDirExpand", arg0)
}

// HomeDir mocks base method
func (_m *MockTarmak) HomeDir() string {
	ret := _m.ctrl.Call(_m, "HomeDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// HomeDir indicates an expected call of HomeDir
func (_mr *MockTarmakMockRecorder) HomeDir() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "HomeDir")
}

// MockPacker is a mock of Packer interface
type MockPacker struct {
	ctrl     *gomock.Controller
	recorder *MockPackerMockRecorder
}

// MockPackerMockRecorder is the mock recorder for MockPacker
type MockPackerMockRecorder struct {
	mock *MockPacker
}

// NewMockPacker creates a new mock instance
func NewMockPacker(ctrl *gomock.Controller) *MockPacker {
	mock := &MockPacker{ctrl: ctrl}
	mock.recorder = &MockPackerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockPacker) EXPECT() *MockPackerMockRecorder {
	return _m.recorder
}

// MockTerraform is a mock of Terraform interface
type MockTerraform struct {
	ctrl     *gomock.Controller
	recorder *MockTerraformMockRecorder
}

// MockTerraformMockRecorder is the mock recorder for MockTerraform
type MockTerraformMockRecorder struct {
	mock *MockTerraform
}

// NewMockTerraform creates a new mock instance
func NewMockTerraform(ctrl *gomock.Controller) *MockTerraform {
	mock := &MockTerraform{ctrl: ctrl}
	mock.recorder = &MockTerraformMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockTerraform) EXPECT() *MockTerraformMockRecorder {
	return _m.recorder
}

// Output mocks base method
func (_m *MockTerraform) Output(stack Stack) (map[string]interface{}, error) {
	ret := _m.ctrl.Call(_m, "Output", stack)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Output indicates an expected call of Output
func (_mr *MockTerraformMockRecorder) Output(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Output", arg0)
}

// MockSSH is a mock of SSH interface
type MockSSH struct {
	ctrl     *gomock.Controller
	recorder *MockSSHMockRecorder
}

// MockSSHMockRecorder is the mock recorder for MockSSH
type MockSSHMockRecorder struct {
	mock *MockSSH
}

// NewMockSSH creates a new mock instance
func NewMockSSH(ctrl *gomock.Controller) *MockSSH {
	mock := &MockSSH{ctrl: ctrl}
	mock.recorder = &MockSSHMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockSSH) EXPECT() *MockSSHMockRecorder {
	return _m.recorder
}

// WriteConfig mocks base method
func (_m *MockSSH) WriteConfig() error {
	ret := _m.ctrl.Call(_m, "WriteConfig")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteConfig indicates an expected call of WriteConfig
func (_mr *MockSSHMockRecorder) WriteConfig() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WriteConfig")
}

// PassThrough mocks base method
func (_m *MockSSH) PassThrough(_param0 []string) {
	_m.ctrl.Call(_m, "PassThrough", _param0)
}

// PassThrough indicates an expected call of PassThrough
func (_mr *MockSSHMockRecorder) PassThrough(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PassThrough", arg0)
}

// Tunnel mocks base method
func (_m *MockSSH) Tunnel(hostname string, destination string, destinationPort int) Tunnel {
	ret := _m.ctrl.Call(_m, "Tunnel", hostname, destination, destinationPort)
	ret0, _ := ret[0].(Tunnel)
	return ret0
}

// Tunnel indicates an expected call of Tunnel
func (_mr *MockSSHMockRecorder) Tunnel(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Tunnel", arg0, arg1, arg2)
}

// Execute mocks base method
func (_m *MockSSH) Execute(host string, cmd string, args []string) (int, error) {
	ret := _m.ctrl.Call(_m, "Execute", host, cmd, args)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (_mr *MockSSHMockRecorder) Execute(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Execute", arg0, arg1, arg2)
}

// MockTunnel is a mock of Tunnel interface
type MockTunnel struct {
	ctrl     *gomock.Controller
	recorder *MockTunnelMockRecorder
}

// MockTunnelMockRecorder is the mock recorder for MockTunnel
type MockTunnelMockRecorder struct {
	mock *MockTunnel
}

// NewMockTunnel creates a new mock instance
func NewMockTunnel(ctrl *gomock.Controller) *MockTunnel {
	mock := &MockTunnel{ctrl: ctrl}
	mock.recorder = &MockTunnelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockTunnel) EXPECT() *MockTunnelMockRecorder {
	return _m.recorder
}

// Start mocks base method
func (_m *MockTunnel) Start() error {
	ret := _m.ctrl.Call(_m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (_mr *MockTunnelMockRecorder) Start() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start")
}

// Stop mocks base method
func (_m *MockTunnel) Stop() error {
	ret := _m.ctrl.Call(_m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (_mr *MockTunnelMockRecorder) Stop() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop")
}

// Port mocks base method
func (_m *MockTunnel) Port() int {
	ret := _m.ctrl.Call(_m, "Port")
	ret0, _ := ret[0].(int)
	return ret0
}

// Port indicates an expected call of Port
func (_mr *MockTunnelMockRecorder) Port() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Port")
}

// MockHost is a mock of Host interface
type MockHost struct {
	ctrl     *gomock.Controller
	recorder *MockHostMockRecorder
}

// MockHostMockRecorder is the mock recorder for MockHost
type MockHostMockRecorder struct {
	mock *MockHost
}

// NewMockHost creates a new mock instance
func NewMockHost(ctrl *gomock.Controller) *MockHost {
	mock := &MockHost{ctrl: ctrl}
	mock.recorder = &MockHostMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockHost) EXPECT() *MockHostMockRecorder {
	return _m.recorder
}

// ID mocks base method
func (_m *MockHost) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (_mr *MockHostMockRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ID")
}

// Hostname mocks base method
func (_m *MockHost) Hostname() string {
	ret := _m.ctrl.Call(_m, "Hostname")
	ret0, _ := ret[0].(string)
	return ret0
}

// Hostname indicates an expected call of Hostname
func (_mr *MockHostMockRecorder) Hostname() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Hostname")
}

// User mocks base method
func (_m *MockHost) User() string {
	ret := _m.ctrl.Call(_m, "User")
	ret0, _ := ret[0].(string)
	return ret0
}

// User indicates an expected call of User
func (_mr *MockHostMockRecorder) User() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "User")
}

// Roles mocks base method
func (_m *MockHost) Roles() []string {
	ret := _m.ctrl.Call(_m, "Roles")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Roles indicates an expected call of Roles
func (_mr *MockHostMockRecorder) Roles() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Roles")
}

// SSHConfig mocks base method
func (_m *MockHost) SSHConfig() string {
	ret := _m.ctrl.Call(_m, "SSHConfig")
	ret0, _ := ret[0].(string)
	return ret0
}

// SSHConfig indicates an expected call of SSHConfig
func (_mr *MockHostMockRecorder) SSHConfig() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SSHConfig")
}