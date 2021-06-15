(function () {
    var deviceKey = '{{.DeviceKey}}';
    var payloadData = _.getPayload();
    var payload = _.parseJson(payloadData);
    if (payload) {
        var contact = payload["contact"];
        if (contact != null) {
            _.setSharedMem(deviceKey, 'active', contact === false);
        }
        _.setSharedMemNonNull(deviceKey, 'wifi_signal', payload['linkquality']);
        _.setSharedMemNonNull(deviceKey, 'battery', payload['battery']);
        _.setSharedMemNonNull(deviceKey, 'battery_low', payload['battery_low']);
        _.setSharedMemNonNull(deviceKey, 'voltage', payload['voltage']);
    } else {
        _.log.error(deviceKey + ' payload not found');
    }
})();
