package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"balancer-api/handlers"
	"balancer-api/models"
)

var _ = Describe("GetAllRecords handler", func() {
	BeforeEach(func() {
		setup()
	})

	JustBeforeEach(func() {
		handler = handlers.Handler{fakeRecordService}
		handlerFunc := http.HandlerFunc(handler.GetAllRecords)
		handlerFunc.ServeHTTP(res, req)
	})

	Context("when GetAllRecords is called", func() {
		var mockBody []models.Record
		var mockBodyResponse []byte

		BeforeEach(func() {
			mockBody = []models.Record{
				models.Record{ID: 0, Name: "test name", Type: "ASSET", Balance: 3.50},
				models.Record{ID: 1, Name: "test name", Type: "Liability", Balance: 2.50},
			}
			mockBodyResponse, _ = json.Marshal(mockBody)

			fakeRecordService.GetAllRecordsReturns(&[]models.Record{
				models.Record{ID: 0, Name: "test name", Type: "ASSET", Balance: 3.50},
				models.Record{ID: 1, Name: "test name", Type: "Liability", Balance: 2.50},
			}, nil)
			req = httptest.NewRequest("GET", "/records", nil)
		})

		Context("when a db error occurs", func() {
			BeforeEach(func() {
				fakeRecordService.GetAllRecordsReturns(nil, errors.New("mock error"))
			})

			It("returns status Internal Server Error", func() {
				Expect(res.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		It("returns status OK with records", func() {
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal(string(mockBodyResponse) + "\n"))
		})
	})
})
