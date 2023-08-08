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
	"github.com/SawitProRecruitment/UserService/utils/authentication"
	"github.com/SawitProRecruitment/UserService/utils/context_helper"
	"github.com/SawitProRecruitment/UserService/utils/response"
	"github.com/SawitProRecruitment/UserService/utils/string_helper"
	"github.com/SawitProRecruitment/UserService/utils/test_helper"
	"github.com/c2fo/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/xorcare/pointer"
)

func TestUpdateUser(t *testing.T) {
	var (
		validReqBody = generated.UpdateUserJSONRequestBody{
			FullName:    pointer.String(test_helper.TestUserName),
			PhoneNumber: pointer.String(test_helper.TestUserPhone),
		}
		validReqBodyNameOnly = generated.UpdateUserJSONRequestBody{
			FullName: pointer.String("user name updated"),
		}
		validReqBodyPhoneOnly = generated.UpdateUserJSONRequestBody{
			PhoneNumber: pointer.String("+625638301212"),
		}

		validUser = repository.User{
			Id:          test_helper.TestUserId,
			PhoneNumber: "+6212345678900",
			FullName:    test_helper.TestUserName,
		}
	)

	testCases := []struct {
		title        string
		jwt          string
		request      *generated.UpdateUserJSONRequestBody
		aborted      bool
		invalidMime  bool
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
			title:            "empty request body",
			jwt:              test_helper.TestUserJWT,
			request:          nil,
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `request need to have either fullName or phoneNumber`,
		},
		{
			title:            "req body Content-Type is not application/json",
			jwt:              test_helper.TestUserJWT,
			request:          &validReqBody,
			invalidMime:      true,
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `Unsupported Media Type`,
		},
		{
			title: "fullName too short",
			jwt:   test_helper.TestUserJWT,
			request: &generated.UpdateUserJSONRequestBody{
				FullName: pointer.String("a"),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `fullName must be at least 3 characters in length`,
		},
		{
			title: "fullName too loog",
			jwt:   test_helper.TestUserJWT,
			request: &generated.UpdateUserJSONRequestBody{
				FullName: pointer.String(strings.Repeat("a", 70)),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `fullName must be a maximum of 60 characters in length`,
		},
		{
			title: "phoneNumber too short",
			jwt:   test_helper.TestUserJWT,
			request: &generated.UpdateUserJSONRequestBody{
				PhoneNumber: pointer.String("54321"),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber must be at least 10 characters in length`,
		},
		{
			title: "phoneNumber too long",
			jwt:   test_helper.TestUserJWT,
			request: &generated.UpdateUserJSONRequestBody{
				PhoneNumber: pointer.String(strings.Repeat("123", 10)),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber must be a maximum of 13 characters in length`,
		},
		{
			title: "phoneNumber do not start with +62",
			jwt:   test_helper.TestUserJWT,
			request: &generated.UpdateUserJSONRequestBody{
				PhoneNumber: pointer.String("12345678901"),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber should start with +62`,
		},
		{
			title: "phoneNumber not following E.164 phone format",
			jwt:   test_helper.TestUserJWT,
			request: &generated.UpdateUserJSONRequestBody{
				PhoneNumber: pointer.String("+6212AbcDef"),
			},
			expectations:     func(t *testing.T, s *serverMock) {},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrMsg:   `phoneNumber must be a valid E.164 formatted phone number`,
		},
		{
			title:   "error in Repository.GetUser by Id",
			jwt:     test_helper.TestUserJWT,
			request: &validReqBody,
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
			title:   "error in Repository.GetUser by PhoneNumber",
			jwt:     test_helper.TestUserJWT,
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					Id: test_helper.TestUserId,
				}).
					Return(validUser, nil)

				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: string_helper.GetAndTrimPointerStringValue(validReqBody.PhoneNumber),
				}).
					Return(repository.User{}, errors.New(response.InternalServerErrorMsg))
			},
			expectedHttpCode: http.StatusInternalServerError,
			expectedErrMsg:   response.InternalServerErrorMsg,
		},
		{
			title:   "phoneNumber in update request already registered for another user",
			jwt:     test_helper.TestUserJWT,
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					Id: test_helper.TestUserId,
				}).
					Return(validUser, nil)

				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: string_helper.GetAndTrimPointerStringValue(validReqBody.PhoneNumber),
				}).
					Return(validUser, nil)
			},
			expectedHttpCode: http.StatusConflict,
			expectedErrMsg:   response.PhoneAlreadyRegisteredErrorMsg,
		},
		{
			title:   "error in Repository.UpdateUser",
			jwt:     test_helper.TestUserJWT,
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					Id: test_helper.TestUserId,
				}).
					Return(validUser, nil)

				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: string_helper.GetAndTrimPointerStringValue(validReqBody.PhoneNumber),
				}).
					Return(repository.User{}, sql.ErrNoRows)

				s.repository.EXPECT().UpdateUser(gomock.Any(), nil, gomock.Any()).
					Return(errors.New(response.InternalServerErrorMsg))
			},
			expectedHttpCode: http.StatusInternalServerError,
			expectedErrMsg:   response.InternalServerErrorMsg,
		},
		{
			title:   "success - both fullName and phoneNumber updated",
			jwt:     test_helper.TestUserJWT,
			request: &validReqBody,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					Id: test_helper.TestUserId,
				}).
					Return(validUser, nil)

				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: string_helper.GetAndTrimPointerStringValue(validReqBody.PhoneNumber),
				}).
					Return(repository.User{}, sql.ErrNoRows)

				s.repository.EXPECT().UpdateUser(gomock.Any(), nil, gomock.Any()).
					Return(nil)
			},
			expectedHttpCode: http.StatusOK,
			expectedResp: generated.UserDataResponse{
				FullName:    string_helper.GetAndTrimPointerStringValue(validReqBody.FullName),
				PhoneNumber: string_helper.GetAndTrimPointerStringValue(validReqBody.PhoneNumber),
			},
		},
		{
			title:   "success - only phoneNumber updated",
			jwt:     test_helper.TestUserJWT,
			request: &validReqBodyPhoneOnly,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					Id: test_helper.TestUserId,
				}).
					Return(validUser, nil)

				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					PhoneNumber: string_helper.GetAndTrimPointerStringValue(validReqBodyPhoneOnly.PhoneNumber),
				}).
					Return(repository.User{}, sql.ErrNoRows)

				s.repository.EXPECT().UpdateUser(gomock.Any(), nil, gomock.Any()).
					Return(nil)
			},
			expectedHttpCode: http.StatusOK,
			expectedResp: generated.UserDataResponse{
				FullName:    validUser.FullName,
				PhoneNumber: string_helper.GetAndTrimPointerStringValue(validReqBodyPhoneOnly.PhoneNumber),
			},
		},
		{
			title:   "success - only fullName updated",
			jwt:     test_helper.TestUserJWT,
			request: &validReqBodyNameOnly,
			expectations: func(t *testing.T, s *serverMock) {
				s.repository.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					Id: test_helper.TestUserId,
				}).
					Return(validUser, nil)

				s.repository.EXPECT().UpdateUser(gomock.Any(), nil, gomock.Any()).
					Return(nil)
			},
			expectedHttpCode: http.StatusOK,
			expectedResp: generated.UserDataResponse{
				FullName:    string_helper.GetAndTrimPointerStringValue(validReqBodyNameOnly.FullName),
				PhoneNumber: validUser.PhoneNumber,
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
			req := httptest.NewRequest(echo.PUT, "/", reqBody)
			req.Header.Set(authentication.AuthHeaderKey, tc.jwt)
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
			ctx.SetPath("/v1/user")

			tc.expectations(t, s)

			err := s.server.UpdateUser(ctx)

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
