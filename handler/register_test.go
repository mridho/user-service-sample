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
	"github.com/xorcare/pointer"
)

func TestRegister(t *testing.T) {
	var (
		validReqBody = generated.RegisterJSONRequestBody{
			FullName:    test_helper.TestUserName,
			Password:    test_helper.TestUserPassword,
			PhoneNumber: test_helper.TestUserPhone,
		}

		validUser = repository.User{
			Id:          test_helper.TestUserId,
			PhoneNumber: test_helper.TestUserPhone,
			FullName:    test_helper.TestUserName,
		}
	)

	testCases := []struct {
		title        string
		request      *generated.RegisterJSONRequestBody
		aborted      bool
		invalidMime  bool
		expectations func(t *testing.T, s *serverMock)

		expectedHttpCode int
		expectedErrMsg   string
		expectedResp     generated.RegisterResponse
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
			expectedErrMsg:   `{"messages":["fullName is a required field","password is a required field","phoneNumber is a required field"]}`,
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
			title: "fullName too short",
			request: &generated.RegisterJSONRequestBody{
				FullName: "a",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `fullName must be at least 3 characters in length`,
		},
		{
			title: "fullName too loog",
			request: &generated.RegisterJSONRequestBody{
				FullName: strings.Repeat("a", 70),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `fullName must be a maximum of 60 characters in length`,
		},
		{
			title: "phoneNumber too short",
			request: &generated.RegisterJSONRequestBody{
				FullName:    test_helper.TestUserName,
				PhoneNumber: "54321",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber must be at least 10 characters in length`,
		},
		{
			title: "phoneNumber too long",
			request: &generated.RegisterJSONRequestBody{
				FullName:    test_helper.TestUserName,
				PhoneNumber: strings.Repeat("123", 10),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber must be a maximum of 13 characters in length`,
		},
		{
			title: "phoneNumber do not start with +62",
			request: &generated.RegisterJSONRequestBody{
				FullName:    test_helper.TestUserName,
				PhoneNumber: "12345678901",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber should start with +62`,
		},
		{
			title: "phoneNumber not following E.164 phone format",
			request: &generated.RegisterJSONRequestBody{
				FullName:    test_helper.TestUserName,
				PhoneNumber: "+6212AbcDef",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber must be a valid E.164 formatted phone number`,
		},
		{
			title: "password too short",
			request: &generated.RegisterJSONRequestBody{
				FullName:    test_helper.TestUserName,
				PhoneNumber: test_helper.TestUserPhone,
				Password:    "abc",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `password must be at least 6 characters in length`,
		},
		{
			title: "password too long",
			request: &generated.RegisterJSONRequestBody{
				FullName:    test_helper.TestUserName,
				PhoneNumber: test_helper.TestUserPhone,
				Password:    strings.Repeat("abc", 30),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `password must be a maximum of 64 characters in length`,
		},
		{
			title: "password does not contain capital letter",
			request: &generated.RegisterJSONRequestBody{
				FullName:    test_helper.TestUserName,
				PhoneNumber: test_helper.TestUserPhone,
				Password:    "abcdefg",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `password must contain at least 1 capital characters, 1 number, and 1 special (non alpha-numeric) character`,
		},
		{
			title: "password does not contain number",
			request: &generated.RegisterJSONRequestBody{
				FullName:    test_helper.TestUserName,
				PhoneNumber: test_helper.TestUserPhone,
				Password:    "aBcdefg",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `password must contain at least 1 capital characters, 1 number, and 1 special (non alpha-numeric) character`,
		},
		{
			title: "password does not contain special character",
			request: &generated.RegisterJSONRequestBody{
				FullName:    test_helper.TestUserName,
				PhoneNumber: test_helper.TestUserPhone,
				Password:    "aB3defg",
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `password must contain at least 1 capital characters, 1 number, and 1 special (non alpha-numeric) character`,
		},
		{
			title:   "error in Repository.GetUser",
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: validReqBody.PhoneNumber,
				}).Return(repository.User{}, errors.New(response.InternalServerErrorMsg))
			},
			expectedHttpCode: http.StatusInternalServerError,
			expectedErrMsg:   response.InternalServerErrorMsg,
		},
		{
			title:   "phone number already registered",
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: validReqBody.PhoneNumber,
				}).Return(validUser, nil)
			},
			expectedHttpCode: http.StatusConflict,
			expectedErrMsg:   response.PhoneAlreadyRegisteredErrorMsg,
		},
		{
			title:   "error in Repository.InsertUser",
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: validReqBody.PhoneNumber,
				}).
					Return(repository.User{}, sql.ErrNoRows)

				s.repository.EXPECT().InsertUser(gomock.Any(), nil, gomock.AssignableToTypeOf(repository.InsertUserInput{})).
					Return(repository.InsertUserOutput{}, errors.New(response.InternalServerErrorMsg))
			},
			expectedHttpCode: http.StatusInternalServerError,
			expectedErrMsg:   response.InternalServerErrorMsg,
		},
		{
			title:   "success",
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: validReqBody.PhoneNumber,
				}).
					Return(repository.User{}, sql.ErrNoRows)

				s.repository.EXPECT().InsertUser(gomock.Any(), nil, gomock.AssignableToTypeOf(repository.InsertUserInput{})).
					Return(repository.InsertUserOutput{
						Id: test_helper.TestUserId,
					}, nil)
			},
			expectedHttpCode: http.StatusCreated,
			expectedResp: generated.RegisterResponse{
				Id:      test_helper.TestUserId,
				Message: pointer.String("user registration success"),
			},
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
			ctx.SetPath("/v1/register")

			tc.expectations(t, s)

			err := s.server.Register(ctx)

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
