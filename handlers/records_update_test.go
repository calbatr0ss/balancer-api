package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"balancer-api/handlers"
	"balancer-api/models"

	"github.com/go-chi/chi"
)

var _ = Describe("UpdateRecord handler", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		handler = handlers.Handler{fakeRecordService}
		handlerFunc := http.HandlerFunc(handler.UpdateRecord)
		handlerFunc.ServeHTTP(res, req)
	})

	Context("when UpdateRecord is called", func() {
		Context("when a url param ID is not provided", func() {
			It("returns status Bad Request", func() {
				Expect(res.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when a url param ID is provided", func() {
			BeforeEach(func() {
				fakeRecordService.GetRecordReturns(&models.Record{ID: 1, Name: "test", Type: "LIABILITY", Balance: 2.46}, nil)
				fakeRecordService.UpdateRecordReturns(nil)
				req = httptest.NewRequest("PUT", "/records", strings.NewReader(`{"name":"updated","type":"asset","balance":42}`))
			})

			JustBeforeEach(func() {
				// Add URL param
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			})

			Context("when bad json is provided", func() {
				BeforeEach(func() {
					req = httptest.NewRequest("PUT", "/records", strings.NewReader(`{"name":"test"}`))
				})

				It("returns status Bad Request", func() {
					Expect(res.Code).To(Equal(http.StatusBadRequest))
				})
			})

			Context("when a db error occurs on get", func() {
				BeforeEach(func() {
					fakeRecordService.GetRecordReturns(nil, errors.New("mock error"))
				})

				It("returns status Internal Server Error", func() {
					Expect(res.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			Context("when a db error occurs on update", func() {
				BeforeEach(func() {
					fakeRecordService.UpdateRecordReturns(errors.New("mock error"))
				})

				It("returns status Internal Server Error", func() {
					Expect(res.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			Context("when an invalid type is provided", func() {
				BeforeEach(func() {
					req = httptest.NewRequest("PUT", "/records", strings.NewReader(`{"name":"updated","type":"yeet","balance":42}`))
				})

				It("returns status Bad Request", func() {
					Expect(res.Code).To(Equal(http.StatusBadRequest))
				})
			})

			It("returns status OK", func() {
				Expect(res.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
