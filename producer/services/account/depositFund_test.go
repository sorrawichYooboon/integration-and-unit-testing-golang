package accountservice

import (
	"errors"
	"producer/commands"
	mockService "producer/services/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_accountService_DepositFund(t *testing.T) {
	mockEventProducer := mockService.NewIEventProducer(t)

	clearAllMock := func() {
		mockEventProducer.ClearAll()
	}
	tests := []struct {
		name               string
		mockServiceRequest commands.DepositFundCommand

		wantServiceOrRepoCallWithAndResponse func()
		wantServiceOrRepoCallTimes           map[string]map[string]int
		wantServiceOrRepoError               error
		wantMainServiceError                 error
		wantMainServiceResponse              interface{}
	}{
		{
			name: "Test should return error when id request is empty",
			mockServiceRequest: commands.DepositFundCommand{
				ID:     "",
				Amount: 1000,
			},
			wantMainServiceError: errors.New("bad request"),
		},
		{
			name: "Test should return error when amount request is zero",
			mockServiceRequest: commands.DepositFundCommand{
				ID:     "123",
				Amount: 0,
			},
		},
		{
			name: "Test should return error when produce of event producer service return error",
			mockServiceRequest: commands.DepositFundCommand{
				ID:     "123",
				Amount: 1000,
			},
			wantServiceOrRepoCallWithAndResponse: func() {
				mockEventProducer.On("Produce", mock.Anything).Return(errors.New("error"))
			},
			wantMainServiceError: errors.New("error"),
		},
		{
			name: "Test should return nil when produce of event producer service return nil",
			mockServiceRequest: commands.DepositFundCommand{
				ID:     "123",
				Amount: 1000,
			},
			wantServiceOrRepoCallWithAndResponse: func() {
				mockEventProducer.On("Produce", mock.Anything).Return(nil)
			},
			wantMainServiceError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer clearAllMock()

			if test.wantServiceOrRepoCallWithAndResponse != nil {
				test.wantServiceOrRepoCallWithAndResponse()
			}

			accountService := NewAccountService(mockEventProducer)
			err := accountService.DepositFund(test.mockServiceRequest)

			if test.wantMainServiceError != nil {
				assert.Equal(t, test.wantMainServiceError.Error(), err.Error())
			}

			for serviceName, serviceCallTimes := range test.wantServiceOrRepoCallTimes {
				for methodName, times := range serviceCallTimes {
					switch serviceName {
					case "eventProducerService":
						mockEventProducer.AssertNumberOfCalls(t, methodName, times)
					default:
						t.Errorf("service %s or method %s not found", serviceName, methodName)
					}
				}
			}
		})
	}
}
