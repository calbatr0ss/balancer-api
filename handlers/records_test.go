package handlers_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"balancer-api/handlers"
)

var _ = Describe("Record", func() {
	Describe("GetAllRecords", func() {
		BeforeEach(func() {
			res = httptest.NewRecorder()
		})

		JustBeforeEach(func() {
			handler := http.HandlerFunc(handlers.GetAllRecords)
			handler.ServeHTTP(res, req)
		})

		Context("when GetAllRecords is called", func() {
			BeforeEach(func() {
				req = httptest.NewRequest("GET", "/record", nil)
			})

			It("returns 200 with records", func() {
				Expect(res.Code).To(Equal(http.StatusOK))
			})

		})

	})

})
