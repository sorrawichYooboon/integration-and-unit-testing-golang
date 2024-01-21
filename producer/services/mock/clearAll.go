package mockService

import "github.com/stretchr/testify/mock"

func (m *IEventProducer) ClearAll() {
	m.Mock = mock.Mock{}
}

func (m *IAccountService) ClearAll() {
	m.Mock = mock.Mock{}
}
