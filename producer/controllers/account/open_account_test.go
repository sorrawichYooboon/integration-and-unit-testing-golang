package accountcontrollers

import (
	"producer/commands"
	internal "producer/internal"
	mockService "producer/services/mock"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func Test_Controller_Open_Account(t *testing.T) {
	mockAccountService := mockService.NewIAccountService(t)

	clearAllMock := func() {
		mockAccountService.ClearAll()
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
			name: "Test should return error if OpenAccount of account service return error",
			mockPayload: commands.OpenAccountCommand{
				AccountHolder:  "test",
				AccountType:    1,
				OpeningBalance: 1000,
			},
			wantServiceCallWithAndResponse: func() {
				mockAccountService.On("OpenAccount", commands.OpenAccountCommand{
					AccountHolder:  "test",
					AccountType:    1,
					OpeningBalance: 1000,
				}).Return("", fiber.ErrInternalServerError)
			},
			wantServiceCallTimes: map[string]map[string]int{
				"accountService": {
					"OpenAccount": 1,
				},
			},
			wantStatusCode: 500,
		},
		{
			name: "Test should return success if OpenAccount of account service return success",
			mockPayload: commands.OpenAccountCommand{
				AccountHolder:  "test",
				AccountType:    1,
				OpeningBalance: 1000,
			},
			wantServiceCallWithAndResponse: func() {
				mockAccountService.On("OpenAccount", commands.OpenAccountCommand{
					AccountHolder:  "test",
					AccountType:    1,
					OpeningBalance: 1000,
				}).Return("test", nil)
			},
			wantServiceCallTimes: map[string]map[string]int{
				"accountService": {
					"OpenAccount": 1,
				},
			},
			wantStatusCode: 201,
			wantControllerResponse: fiber.Map{
				"message": "open account success",
				"id":      "test",
			},
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

			accountContrller := NewAccountController(mockAccountService)
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
					case "accountService":
						mockAccountService.AssertNumberOfCalls(t, methodName, times)
					default:
						t.Errorf("service %s or method %s not found", serviceName, methodName)
					}
				}
			}
		})
	}
}
