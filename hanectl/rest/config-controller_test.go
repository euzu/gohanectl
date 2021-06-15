package rest

import (
	"github.com/go-chi/chi"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/test/mock_test"
	"gohanectl/hanectl/test/rest_test"
	"net/http/httptest"
	"testing"
)

func TestServerStatus(t *testing.T) {
	serviceFactoryMock := new(mock_test.ServiceFactoryMock)
	configServiceMock := new(mock_test.ConfigServiceMock)
	serviceFactoryMock.On("GetConfigService").Return(configServiceMock)
	configServiceMock.On("GetServerStatus").Return(new(model.ServerStatus), nil)

	r := chi.NewRouter()
	configApi(r, serviceFactoryMock)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if resp, body := rest_test.DoTestRequest(t, ts, "GET", "/status", nil); body != "root" && resp.StatusCode != 200 {
		t.Fatalf(body)
	}

	serviceFactoryMock.AssertExpectations(t)
	configServiceMock.AssertExpectations(t)
}

func TestConfigRooms(t *testing.T) {
	serviceFactoryMock := new(mock_test.ServiceFactoryMock)
	configServiceMock := new(mock_test.ConfigServiceMock)
	serviceFactoryMock.On("GetConfigService").Return(configServiceMock)
	configServiceMock.On("GetRoomsConfig").Return(make(model.Dictionary), nil)

	r := chi.NewRouter()
	configApi(r, serviceFactoryMock)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if resp, body := rest_test.DoTestRequest(t, ts, "GET", "/config/room", nil); body != "root" && resp.StatusCode != 200 {
		t.Fatalf(body)
	}

	serviceFactoryMock.AssertExpectations(t)
	configServiceMock.AssertExpectations(t)

}