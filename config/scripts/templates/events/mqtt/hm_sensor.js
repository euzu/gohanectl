var deviceKey = '{{.DeviceKey}}';
var topic = _.getTopic();
var payload = _.getPayload();
var state = _.parseJson(payload);
if (state) {
    var value = state["state"];
    if (value) {
        if (value === "open") {
            _.setSharedMem(deviceKey, "open", true);
            _.setSharedMem(deviceKey, "ts", new Date().toISOString());
        } else if (value === "closed") {
            _.setSharedMem(deviceKey, "open", false);
        } else if (_.startsWith(value, "battery: ")) {
            _.setSharedMem(deviceKey, "battery", value.slice("battery: ".length));
        } else if (_.startsWith(value, "trigger_cnt: ")) {
            _.setSharedMem(deviceKey, "trigger_cnt", value.slice("trigger_cnt: ".length));
        }
    }
} else {
    _.log.debug(deviceKey + ' : ' + topic + '  : ' + payload);
}

