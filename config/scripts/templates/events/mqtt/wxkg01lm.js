// click can be
// single, double, triple, quadruple, long, long_release, many.
var deviceKey = '{{.DeviceKey}}';
var payloadData = _.getPayload();
var payload = _.parseJson(payloadData);
if (payload) {
    _.setSharedMemNonNull(deviceKey, 'wifi_signal', payload['linkquality']);
    _.setSharedMemNonNull(deviceKey, 'battery', payload['battery']);
    _.setSharedMemNonNull(deviceKey, 'action', payload['action']);
    _.setSharedMemNonNull(deviceKey, 'voltage', payload['voltage']);
} else {
    _.log.error(deviceKey + ' payload not found');
}
