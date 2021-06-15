(function () {
    var deviceKey = '{{.DeviceKey}}';
    var payloadData = _.getPayload();
    var payload = _.parseJson(payloadData);
    if (payload) {
        var state = payload["state"];
        if (state) {
            _.setSharedMem(deviceKey, 'active', state === 'on');
        }
        _.setSharedMemNonNull(deviceKey, 'wifi_signal', payload['linkquality']);
        _.setSharedMemNonNull(deviceKey, 'current', payload['current']);
        _.setSharedMemNonNull(deviceKey, 'voltage', payload['voltage']);
        _.setSharedMemNonNull(deviceKey, 'power', payload['power']);
        _.setSharedMemNonNull(deviceKey, 'total', payload['energy']);
    } else {
        _.log.error(deviceKey + ' payload not found');
    }
})();