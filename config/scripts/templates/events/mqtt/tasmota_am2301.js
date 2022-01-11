var deviceKey = '{{.DeviceKey}}';
var topic = _.getTopic();
var payload = _.getPayload();
var state = _.parseJson(payload);
if (state) {
    var value = state["am2301"];
    if (value) {
        _.setSharedMemNonNull(deviceKey, 'temperature', value['temperature']);
        _.setSharedMemNonNull(deviceKey, 'humidity', value['humidity']);
        _.setSharedMemNonNull(deviceKey, 'dewpoint', value['dewpoint']);
        _.setSharedMemNonNull(deviceKey, 'tempunit', state['tempunit']);
    }
} else {
    _.log.debug(deviceKey + ' : ' + topic + '  : ' + payload);
}
