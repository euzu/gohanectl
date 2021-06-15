(function () {
    var deviceKey = '{{.DeviceKey}}';
    var payload = _.getPayload();

    var data = JSON.parse(payload)
    if (data) {
        var sensor = data["statussns"];
        if (sensor && sensor['am2301']) {
            _.setSharedMemNonNull(deviceKey, 'tempunit', sensor['tempunit']);
            var am2301 = sensor['am2301'];
            _.setSharedMemNonNull(deviceKey, 'temperature', am2301['temperature']);
            _.setSharedMemNonNull(deviceKey, 'humidity', am2301['humidity']);
            _.setSharedMemNonNull(deviceKey, 'dewpoint', am2301['dewpoint']);
        }
    } else {
        _.log.debug(deviceKey + ' :  ' + payload);
    }
})();
