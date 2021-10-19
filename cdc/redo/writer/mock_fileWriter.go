// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package writer

import mock "github.com/stretchr/testify/mock"

// mockFileWriter is an autogenerated mock type for the fileWriter type
type mockFileWriter struct {
	mock.Mock
}

// AdvanceTs provides a mock function with given fields: commitTs
func (_m *mockFileWriter) AdvanceTs(commitTs uint64) {
	_m.Called(commitTs)
}

// Close provides a mock function with given fields:
func (_m *mockFileWriter) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Flush provides a mock function with given fields:
func (_m *mockFileWriter) Flush() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GC provides a mock function with given fields: checkPointTs
func (_m *mockFileWriter) GC(checkPointTs uint64) error {
	ret := _m.Called(checkPointTs)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(checkPointTs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsRunning provides a mock function with given fields:
func (_m *mockFileWriter) IsRunning() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Write provides a mock function with given fields: p
func (_m *mockFileWriter) Write(p []byte) (int, error) {
	ret := _m.Called(p)

	var r0 int
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
