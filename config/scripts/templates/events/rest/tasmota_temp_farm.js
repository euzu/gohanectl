var deviceKey = '{{.DeviceKey}}';
var payload = _.getPayload();

var data = JSON.parse(payload)
if (data) {
    var sensor = data["statussns"];
    if (sensor) {
        _.setSharedMemNonNull(deviceKey, 'tempunit', sensor['tempunit']);
        var am2301 = sensor['am2301'];
        if (am2301) {
            _.setSharedMemNonNull(deviceKey, 'dht22', am2301['temperature']);
            _.setSharedMemNonNull(deviceKey, 'humidity', am2301['humidity']);
            _.setSharedMemNonNull(deviceKey, 'dewpoint', am2301['dewpoint']);
        }
        for (var i = 1; i < 9; i++) {
            var key = 'ds18b20-' + i;
            var ds18b20 = sensor[key];
            if (ds18b20) {
                _.setSharedMemNonNull(deviceKey, key, ds18b20['temperature']);
            }
        }
    }
} else {
    _.log.debug(deviceKey + ' :  ' + payload);
}
