package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"balancer-api/balancer/balancerfakes"
	"balancer-api/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	req               *http.Request
	res               *httptest.ResponseRecorder
	fakeRecordService *balancerfakes.FakeRecordService
	handler           handlers.Handler
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

func setup() {
	fakeRecordService = new(balancerfakes.FakeRecordService)
	res = httptest.NewRecorder()
}
