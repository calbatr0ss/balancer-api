package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"balancer-api/handlers"
)

var _ = Describe("CreateRecord handler", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		handler = handlers.Handler{fakeRecordService}
		handlerFunc := http.HandlerFunc(handler.CreateRecord)
		handlerFunc.ServeHTTP(res, req)
	})

	Context("when CreateRecord is called", func() {
		mockResponseBody := `{"recordId":1}` + "\n"

		BeforeEach(func() {
			mockID := uint64(1)
			fakeRecordService.CreateRecordReturns(&mockID, nil)
			req = httptest.NewRequest("POST", "/records", strings.NewReader(`{ "name": "test", "type": "Liability", "balance": 1.50 }`))
		})

		Context("when a db error occurs", func() {
			BeforeEach(func() {
				fakeRecordService.CreateRecordReturns(nil, errors.New("mock error"))
			})

			It("returns status Internal Server Error", func() {
				Expect(res.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("when an invalid type is provided", func() {
			BeforeEach(func() {
				req = httptest.NewRequest("POST", "/records", strings.NewReader(`{ "name": "test", "type": "bad type", "balance": 1.50 }`))
			})

			It("returns status Bad Request", func() {
				Expect(res.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when no name is provided", func() {
			BeforeEach(func() {
				req = httptest.NewRequest("POST", "/records", strings.NewReader(`{ "type": "bad type", "balance": 1.50 }`))
			})

			It("returns status Bad Request", func() {
				Expect(res.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when bad json is provided", func() {
			BeforeEach(func() {
				req = httptest.NewRequest("POST", "/records", strings.NewReader(`{ "type": "bad type", "b`))
			})

			It("returns status Bad Request", func() {
				Expect(res.Code).To(Equal(http.StatusBadRequest))
			})
		})

		It("returns status Created with new record ID", func() {
			Expect(res.Code).To(Equal(http.StatusCreated))
			Expect(res.Body.String()).To(Equal(mockResponseBody))
		})
	})
})
