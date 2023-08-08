package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils/context_helper"
	"github.com/SawitProRecruitment/UserService/utils/response"
	"github.com/SawitProRecruitment/UserService/utils/test_helper"
	"github.com/c2fo/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestLogin(t *testing.T) {

	var (
		validReqBody = generated.LoginJSONRequestBody{
			PhoneNumber: test_helper.TestUserPhone,
			Password:    test_helper.TestUserPassword,
		}

		validUser = repository.User{
			Id:           test_helper.TestUserId,
			PhoneNumber:  test_helper.TestUserPhone,
			FullName:     test_helper.TestUserName,
			PasswordHash: test_helper.TestUserPasswordHash,
			Salt:         test_helper.TestUserSalt,
		}
	)

	testCases := []struct {
		title        string
		request      *generated.LoginJSONRequestBody
		aborted      bool
		invalidMime  bool
		expectations func(t *testing.T, s *serverMock)

		expectedHttpCode int
		expectedErrMsg   string
		expectedResp     string
	}{
		{
			title:            "request aborted",
			aborted:          true,
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusGone,
			expectedErrMsg:   context_helper.ErrRequestCanceled.Error(),
		},
		{
			title:            "empty request body",
			request:          nil,
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `{"messages":["password is a required field","phoneNumber is a required field"]}`,
		},
		{
			title:            "req body Content-Type is not application/json",
			request:          &validReqBody,
			invalidMime:      true,
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `Unsupported Media Type`,
		},
		{
			title: "phoneNumber too short",
			request: &generated.LoginJSONRequestBody{
				PhoneNumber: "54321",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber must be at least 10 characters in length`,
		},
		{
			title: "phoneNumber too long",
			request: &generated.LoginJSONRequestBody{
				PhoneNumber: strings.Repeat("123", 10),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber must be a maximum of 13 characters in length`,
		},
		{
			title: "phoneNumber do not start with +62",
			request: &generated.LoginJSONRequestBody{
				PhoneNumber: "12345678901",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber should start with +62`,
		},
		{
			title: "phoneNumber not following E.164 phone format",
			request: &generated.LoginJSONRequestBody{
				PhoneNumber: "+6212AbcDef",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber must be a valid E.164 formatted phone number`,
		},
		{
			title:   "error in Repository.GetUser",
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: validReqBody.PhoneNumber,
				}).
					Return(repository.User{}, errors.New(response.InternalServerErrorMsg))
			},
			expectedHttpCode: http.StatusInternalServerError,
			expectedErrMsg:   response.InternalServerErrorMsg,
		},
		{
			title:   "incorrect login phoneNumber not registered",
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: validReqBody.PhoneNumber,
				}).
					Return(repository.User{}, sql.ErrNoRows)
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   response.IncorrectLoginErrorMsg,
		},
		{
			title: "incorrect login wrong password",
			request: &generated.LoginJSONRequestBody{
				PhoneNumber: test_helper.TestUserPhone,
				Password:    "wrongP4$sWrd",
			},
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: validReqBody.PhoneNumber,
				}).
					Return(validUser, nil)
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   response.IncorrectLoginErrorMsg,
		},
		{
			title:   "error in Repository.IncrementUserLoginCount only log error - login success",
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: validReqBody.PhoneNumber,
				}).
					Return(validUser, nil)

				s.repository.EXPECT().IncrementUserLoginCount(gomock.Any(), nil, validUser).
					Return(errors.New(response.InternalServerErrorMsg))
			},
			expectedHttpCode: http.StatusOK,
			expectedResp:     validUser.Id,
		},
		{
			title:   "success",
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: validReqBody.PhoneNumber,
				}).
					Return(validUser, nil)

				s.repository.EXPECT().IncrementUserLoginCount(gomock.Any(), nil, validUser).
					Return(nil)
			},
			expectedHttpCode: http.StatusOK,
			expectedResp:     validUser.Id,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// Setup
			s := setupServerMock(t)
			defer s.cleanUp()

			var reqBody io.Reader
			if tc.request != nil {
				reqBodyJson, _ := json.Marshal(*tc.request)
				reqBody = bytes.NewReader(reqBodyJson)
			}

			e := echo.New()
			req := httptest.NewRequest(echo.POST, "/", reqBody)
			if !tc.invalidMime {
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			rec := httptest.NewRecorder()

			if tc.aborted {
				// simulate request aborted
				c, cancel := context.WithCancel(req.Context())
				cancel()
				req = req.WithContext(c)
			}

			ctx := e.NewContext(req, rec)
			ctx.SetPath("/v1/login")

			tc.expectations(t, s)

			err := s.server.Login(ctx)

			// Assertions
			if tc.expectedHttpCode >= http.StatusOK && // code 2XX
				tc.expectedHttpCode <= http.StatusIMUsed {

				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHttpCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.expectedResp)
			} else {
				assert.Equal(t, tc.expectedHttpCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.expectedErrMsg)
			}
		})
	}
}
