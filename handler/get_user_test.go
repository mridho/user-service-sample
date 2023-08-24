package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"user-service-sample/generated"
	"user-service-sample/repository"
	"user-service-sample/utils/authentication"
	"user-service-sample/utils/context_helper"
	"user-service-sample/utils/response"
	"user-service-sample/utils/test_helper"

	"github.com/c2fo/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestGetUser(t *testing.T) {

	var (
		validUser = repository.User{
			Id:          test_helper.TestUserId,
			PhoneNumber: test_helper.TestUserPhone,
			FullName:    test_helper.TestUserName,
		}
	)

	testCases := []struct {
		title        string
		jwt          string
		aborted      bool
		expectations func(t *testing.T, s *serverMock)

		expectedHttpCode int
		expectedErrMsg   string
		expectedResp     generated.UserDataResponse
	}{
		{
			title:            "request aborted",
			aborted:          true,
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusGone,
			expectedErrMsg:   context_helper.ErrRequestCanceled.Error(),
		},
		{
			title:            "auth token invalid",
			jwt:              "Bearer invalid-token",
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusForbidden,
			expectedErrMsg:   response.AccessForbiddenErrorMsg,
		},
		{
			title: "error in Repository.GetUser",
			jwt:   test_helper.TestUserJWT,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					Id: test_helper.TestUserId,
				}).
					Return(repository.User{}, errors.New(response.InternalServerErrorMsg))
			},
			expectedHttpCode: http.StatusInternalServerError,
			expectedErrMsg:   response.InternalServerErrorMsg,
		},
		{
			title: "success",
			jwt:   test_helper.TestUserJWT,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					Id: test_helper.TestUserId,
				}).
					Return(validUser, nil)
			},
			expectedHttpCode: http.StatusOK,
			expectedResp: generated.UserDataResponse{
				FullName:    validUser.FullName,
				PhoneNumber: validUser.PhoneNumber,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// Setup
			s := setupServerMock(t)
			defer s.cleanUp()

			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			req.Header.Set(authentication.AuthHeaderKey, tc.jwt)
			rec := httptest.NewRecorder()

			if tc.aborted {
				// simulate request aborted
				c, cancel := context.WithCancel(req.Context())
				cancel()
				req = req.WithContext(c)
			}

			ctx := e.NewContext(req, rec)
			ctx.SetPath("/v1/user")

			tc.expectations(t, s)

			err := s.server.GetUser(ctx)

			// Assertions
			if tc.expectedHttpCode >= http.StatusOK && // code 2XX
				tc.expectedHttpCode <= http.StatusIMUsed {

				expectedRespJson, _ := json.Marshal(tc.expectedResp)

				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHttpCode, rec.Code)
				assert.Equal(t, string(expectedRespJson), strings.TrimSpace(rec.Body.String()))
			} else {
				assert.Equal(t, tc.expectedHttpCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.expectedErrMsg)
			}
		})
	}
}
