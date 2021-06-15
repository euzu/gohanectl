package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/test/mock_test"
	"gohanectl/hanectl/test/rest_test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func contextUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ClaimsCtxKey, model.Dictionary{
			UsernameClaimKey: "user",
			AuthoritiesClaimKey: []interface{}{"USER"},
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func TestDeviceConfig(t *testing.T) {
	devicesDto := new(model.DevicesDto)
	serviceFactoryMock := new(mock_test.ServiceFactoryMock)
	deviceServiceMock := new(mock_test.DeviceServiceMock)
	serviceFactoryMock.On("GetDeviceService").Return(deviceServiceMock)
	deviceServiceMock.On("GetDevicesDto", mock.AnythingOfTypeArgument("model.DeviceFilter")).Return(devicesDto, nil)

	r := chi.NewRouter()
	r.Use(contextUser)
	deviceApi(r, serviceFactoryMock)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if resp, body := rest_test.DoTestRequest(t, ts, "GET", "/devices", nil); body != "root" && resp.StatusCode != 200 {
		t.Fatalf(body)
	}

	serviceFactoryMock.AssertExpectations(t)
	deviceServiceMock.AssertExpectations(t)
}

func TestDeviceStates(t *testing.T) {
	states := make(model.Dictionary)
	serviceFactoryMock := new(mock_test.ServiceFactoryMock)
	deviceServiceMock := new(mock_test.DeviceServiceMock)
	serviceFactoryMock.On("GetDeviceService").Return(deviceServiceMock)
	deviceServiceMock.On("DeviceStates").Return(states)

	r := chi.NewRouter()
	deviceApi(r, serviceFactoryMock)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if resp, body := rest_test.DoTestRequest(t, ts, "GET", "/devices/status", nil); body != "root" && resp.StatusCode != 200 {
		t.Fatalf(body)
	}

	serviceFactoryMock.AssertExpectations(t)
	deviceServiceMock.AssertExpectations(t)
}

func TestDeviceState(t *testing.T) {
	deviceKey := "deviceKey"
	states := make(model.Dictionary)
	serviceFactoryMock := new(mock_test.ServiceFactoryMock)
	deviceServiceMock := new(mock_test.DeviceServiceMock)
	serviceFactoryMock.On("GetDeviceService").Return(deviceServiceMock)
	deviceServiceMock.On("DeviceState", deviceKey).Return(states)

	r := chi.NewRouter()
	deviceApi(r, serviceFactoryMock)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if resp, body := rest_test.DoTestRequest(t, ts, "GET", "/device/" + deviceKey, nil); body != "root" && resp.StatusCode != 200 {
		t.Fatalf(body)
	}

	serviceFactoryMock.AssertExpectations(t)
	deviceServiceMock.AssertExpectations(t)
}

func TestDeviceCommand(t *testing.T) {
	deviceKey := "deviceKey"
	deviceState := model.Dictionary{
		"device": deviceKey,
	}
	payloadContent, _ := json.Marshal(deviceState)
	payload := bytes.NewReader(payloadContent)
	device := &model.Device{
		DeviceKey: deviceKey,
		Authorities: []string{"USER"},
	}
	serviceFactoryMock := new(mock_test.ServiceFactoryMock)
	deviceServiceMock := new(mock_test.DeviceServiceMock)
	serviceFactoryMock.On("GetDeviceService").Return(deviceServiceMock)
	deviceServiceMock.On("GetDevice", deviceKey).Return(device, nil)
	deviceServiceMock.On("DeviceCommand", deviceKey, deviceState).Return(true)

	r := chi.NewRouter()
	r.Use(contextUser)
	deviceApi(r, serviceFactoryMock)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if resp, body := rest_test.DoTestRequest(t, ts, "POST", "/device", payload); body != "root" && resp.StatusCode != 200 {
		t.Fatalf(body)
	}

	serviceFactoryMock.AssertExpectations(t)
	deviceServiceMock.AssertExpectations(t)
}
