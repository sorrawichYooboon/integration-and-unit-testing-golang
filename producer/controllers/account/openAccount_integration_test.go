//go:build integration

package accountcontrollers

import (
	"producer/commands"
	internal "producer/internal"
	accountservice "producer/services/account"
	mockService "producer/services/mock"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func Test_Integration_Controller_Open_Account(t *testing.T) {
	mockEventProducer := mockService.NewIEventProducer(t)
	accountservice := accountservice.NewAccountService(mockEventProducer)

	clearAllMock := func() {
		mockEventProducer.ClearAll()
	}

	tests := []struct {
		name        string
		mockPayload interface{}

		wantServiceCallWithAndResponse func()
		wantServiceCallTimes           map[string]map[string]int
		wantStatusCode                 int
		wantErrorCode                  string
		wantControllerResponse         interface{}
	}{
		{
			name: "Test should return success if produce event success",
			mockPayload: commands.OpenAccountCommand{
				AccountHolder:  "test",
				AccountType:    1,
				OpeningBalance: 1000,
			},
			wantServiceCallWithAndResponse: func() {
				mockEventProducer.On("Produce", mock.Anything).Return(nil)
			},
			wantServiceCallTimes: map[string]map[string]int{
				"eventProducer": {
					"Produce": 1,
				},
			},
			wantStatusCode: 201,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer clearAllMock()

			jsonByte := internal.MarshalJSONData(test.mockPayload)

			// ------------------ gin ------------------
			// request := internal.CreateHTTPRequest(http.MethodPost, "/mock-endpoint", jsonByte)

			// response := httptest.NewRecorder()
			// ctx, _ := gin.CreateTestContext(response)
			// ctx.Request = request

			// ------------------ fiber ------------------
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Request().SetBody(jsonByte)

			if test.wantServiceCallWithAndResponse != nil {
				test.wantServiceCallWithAndResponse()
			}

			accountContrller := NewAccountController(accountservice)
			accountContrller.OpenAccount(ctx)

			// ------------------ gin ------------------
			// if test.wantStatusCode != 0 {
			// 	assert.Equal(t, test.wantStatusCode, response.Code)
			// }

			// ------------------ fiber ------------------
			if test.wantStatusCode != 0 {
				assert.Equal(t, test.wantStatusCode, ctx.Response().StatusCode())
			}

			if test.wantControllerResponse != nil {
				var wantControllerResponseByte []byte
				switch test.wantControllerResponse.(type) {
				case fiber.Map:
					wantControllerResponseByte = internal.MarshalJSONData(test.wantControllerResponse.(fiber.Map))
				case string:
					wantControllerResponseByte = []byte(test.wantControllerResponse.(string))
				}

				assert.Equal(t, wantControllerResponseByte, ctx.Response().Body())
			}

			for serviceName, serviceCallTimes := range test.wantServiceCallTimes {
				for methodName, times := range serviceCallTimes {
					switch serviceName {
					case "eventProducer":
						mockEventProducer.AssertNumberOfCalls(t, methodName, times)
					default:
						t.Errorf("service %s or method %s not found", serviceName, methodName)
					}
				}
			}
		})
	}
}
