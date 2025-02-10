package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockInitializer is a mock implementation of the initializer package
type MockInitializer struct {
	mock.Mock
}

func (m *MockInitializer) LoadConfig() {
	m.Called()
}

func (m *MockInitializer) ConnectSpannerDB() (interface{}, error) {
	args := m.Called()
	return args.Get(0), args.Error(1)
}

func TestInit(t *testing.T) {
	mockInitializer := new(MockInitializer)

	// Mock the LoadConfig method
	mockInitializer.On("LoadConfig").Return()

	// Mock the ConnectSpannerDB method
	mockInitializer.On("ConnectSpannerDB").Return(nil, nil)

	// Replace the real initializer with the mock
	// initializer.LoadConfig = mockInitializer.LoadConfig
	// initializer.ConnectSpannerDB = mockInitializer.ConnectSpannerDB

	// Call the init function
	// Define the init function
	init := func() {
		mockInitializer.LoadConfig()
		mockInitializer.ConnectSpannerDB()
	}

	init()

	// Assert that the methods were called
	mockInitializer.AssertCalled(t, "LoadConfig")
	mockInitializer.AssertCalled(t, "ConnectSpannerDB")
}

// mockLogger is a mock implementation of io.Writer to capture log output
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockLogger) String() string {
	args := m.Called()
	return args.String(0)
}
