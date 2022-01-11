var deviceKey = '{{.DeviceKey}}';
var payloadData = _.getPayload();
var payload = _.parseJson(payloadData);
if (payload) {
    var state = payload["state"];
    if (state) {
        _.setSharedMem(deviceKey, 'active', state === 'on');
    }
    _.setSharedMemNonNull(deviceKey, 'wifi_signal', payload['linkquality']);
} else {
    _.log.error(deviceKey + ' payload not found');
}
