package accountservice

import (
	"errors"
	"producer/commands"
	mockService "producer/services/mock"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_accountService_OpenAccount(t *testing.T) {
	mockEventProducer := mockService.NewIEventProducer(t)

	clearAllMock := func() {
		mockEventProducer.ClearAll()
	}
	tests := []struct {
		name               string
		mockServiceRequest commands.OpenAccountCommand

		wantServiceOrRepoCallWithAndResponse func()
		wantServiceOrRepoCallTimes           map[string]map[string]int
		wantMainServiceError                 error
		wantMainServiceResponse              string
	}{
		{
			name: "Test should return error when account holder request is empty",
			mockServiceRequest: commands.OpenAccountCommand{
				AccountHolder:  "",
				AccountType:    1,
				OpeningBalance: 1000,
			},
			wantMainServiceError: errors.New("bad request"),
		},
		{
			name: "Test should return error when account type request is zero",
			mockServiceRequest: commands.OpenAccountCommand{
				AccountHolder:  "John Doe",
				AccountType:    0,
				OpeningBalance: 1000,
			},
			wantMainServiceError: errors.New("bad request"),
		},
		{
			name: "Test should return error when opening balance request is zero",
			mockServiceRequest: commands.OpenAccountCommand{
				AccountHolder:  "John Doe",
				AccountType:    1,
				OpeningBalance: 0,
			},
			wantMainServiceError: errors.New("bad request"),
		},
		{
			name: "Test should return error when produce of event producer service return error",
			mockServiceRequest: commands.OpenAccountCommand{
				AccountHolder:  "John Doe",
				AccountType:    1,
				OpeningBalance: 1000,
			},
			wantServiceOrRepoCallWithAndResponse: func() {
				mockEventProducer.On("Produce", mock.Anything).Return(errors.New("error"))
			},
			wantServiceOrRepoCallTimes: map[string]map[string]int{
				"eventProducerService": {
					"Produce": 1,
				},
			},
		},
		{
			name: "Test should return account holder when event producer service return nil",
			mockServiceRequest: commands.OpenAccountCommand{
				AccountHolder:  "John Doe",
				AccountType:    1,
				OpeningBalance: 1000,
			},
			wantServiceOrRepoCallWithAndResponse: func() {
				mockEventProducer.On("Produce", mock.Anything).Return(nil)
			},
			wantServiceOrRepoCallTimes: map[string]map[string]int{
				"eventProducerService": {
					"Produce": 1,
				},
			},
			wantMainServiceResponse: "John Doe",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer clearAllMock()

			if test.wantServiceOrRepoCallWithAndResponse != nil {
				test.wantServiceOrRepoCallWithAndResponse()
			}

			accountService := NewAccountService(mockEventProducer)
			response, err := accountService.OpenAccount(test.mockServiceRequest)

			if test.wantMainServiceError != nil {
				assert.Equal(t, test.wantMainServiceError.Error(), err.Error())
			}

			if !reflect.DeepEqual(test.wantMainServiceResponse, "") {
				assert.Equal(t, test.wantMainServiceResponse, response)
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
