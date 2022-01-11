var deviceKey = '{{.DeviceKey}}';
var payloadData = _.getPayload();
var payload = _.parseJson(payloadData);
if (payload) {
    _.setSharedMemNonNull(deviceKey, 'wifi_signal', payload['linkquality']);
    _.setSharedMemNonNull(deviceKey, 'battery', payload['battery']);
    _.setSharedMemNonNull(deviceKey, 'humidity', payload['humidity']);
    _.setSharedMemNonNull(deviceKey, 'temperature', payload['temperature']);
    _.setSharedMemNonNull(deviceKey, 'voltage', payload['voltage']);
} else {
    _.log.error(deviceKey + ' payload not found');
}
