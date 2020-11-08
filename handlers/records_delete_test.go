package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"balancer-api/handlers"

	"github.com/go-chi/chi"
)

var _ = Describe("DeleteRecord handler", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		handler = handlers.Handler{fakeRecordService}
		handlerFunc := http.HandlerFunc(handler.DeleteRecord)
		handlerFunc.ServeHTTP(res, req)
	})

	Context("when DeleteRecord is called", func() {
		var mockBody handlers.IDPayload
		var mockBodyResponse []byte

		BeforeEach(func() {
			mockBody = handlers.IDPayload{ID: uint64(1)}
			mockBodyResponse, _ = json.Marshal(mockBody)
			fakeRecordService.DeleteRecordReturns(nil)
			req = httptest.NewRequest("DELETE", "/records", nil)
		})

		Context("when a url param ID is not provided", func() {
			It("returns status Bad Request", func() {
				Expect(res.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when a url param ID is provided", func() {
			BeforeEach(func() {
				// Add URL param
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			})

			Context("when a db error occurs", func() {
				BeforeEach(func() {
					fakeRecordService.DeleteRecordReturns(errors.New("mock error"))
				})

				It("returns status Internal Server Error", func() {
					Expect(res.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			It("returns status OK with record ID", func() {
				Expect(res.Code).To(Equal(http.StatusOK))
				Expect(res.Body.String()).To(Equal(string(mockBodyResponse) + "\n"))
			})
		})
	})
})
