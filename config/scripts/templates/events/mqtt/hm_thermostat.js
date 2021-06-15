(function(){
    var deviceKey = '{{.DeviceKey}}';
    var topic = _.getTopic();
    var payload = _.getPayload();
    var state = _.parseJson(payload);
    if (state && state.temp) {
        _.setSharedMem(deviceKey, "temp", parseFloat(state["temp"]));
        _.setSharedMem(deviceKey, "desired", parseFloat(state["desired"]));
        _.setSharedMem(deviceKey, "valve", parseFloat(state["valve"]));
    } else {
        _.log.debug(deviceKey + ' : ' + topic + '  : ' + payload);
    }
})();
