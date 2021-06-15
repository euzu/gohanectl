(function () {
    var deviceKey = '{{.DeviceKey}}';
    var topic = _.getTopic();
    var payload = _.getPayload();
    var state = _.parseJson(payload);
    if (state && state["status"]) {
        _.setSharedMem(deviceKey, "active", state["status"] === "on");
    } else {
        _.log.debug(deviceKey + ' : ' + topic + '  : ' + payload);
    }
})();
