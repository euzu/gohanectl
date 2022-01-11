var deviceKey = '{{.DeviceKey}}';
var topic = _.getTopic();
var payload = _.getPayload();
var state = _.parseJson(payload);
if (state) {
    _.setSharedMemNonNull(deviceKey, 'tempunit', state['tempunit']);
    var am2301 = state["am2301"];
    if (am2301) {
        _.setSharedMemNonNull(deviceKey, 'dht22', am2301['temperature']);
        _.setSharedMemNonNull(deviceKey, 'humidity', am2301['humidity']);
        _.setSharedMemNonNull(deviceKey, 'dewpoint', am2301['dewpoint']);
    }
    for (var i = 1; i < 9; i++) {
        var key = 'ds18b20-' + i;
        var ds18b20 = state[key];
        if (ds18b20) {
            _.setSharedMemNonNull(deviceKey, key, ds18b20['temperature']);
        }
    }
} else {
    _.log.debug(deviceKey + ' : ' + topic + '  : ' + payload);
}
