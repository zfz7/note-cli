// Code generated by mockery v2.46.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ConfigHelper is an autogenerated mock type for the ConfigHelper type
type ConfigHelper struct {
	mock.Mock
}

// ReadConfig provides a mock function with given fields:
func (_m *ConfigHelper) ReadConfig() (NoteConfig, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadConfig")
	}

	var r0 NoteConfig
	var r1 error
	if rf, ok := ret.Get(0).(func() (NoteConfig, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() NoteConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(NoteConfig)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Setup provides a mock function with given fields:
func (_m *ConfigHelper) Setup() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Config")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewConfigHelper creates a new instance of ConfigHelper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConfigHelper(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConfigHelper {
	mock := &ConfigHelper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
