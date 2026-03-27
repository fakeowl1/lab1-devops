package service

import "github.com/stretchr/testify/mock"

type MockHealthyRepo struct {
    mock.Mock
}

func (m *MockHealthyRepo) Ping() error {
    args := m.Called()
    return args.Error(0)
}
