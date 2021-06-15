(function(){
    var deviceKey = '{{.DeviceKey}}';
    var topic = _.getTopic();
    var payload = _.getPayload();
    var parts = _.getSplittedTopic();
    if (parts.length) {
        var field = parts[2];
        if (field === "relay") {
            _.setSharedMem(deviceKey, 'active', payload === 'on');
        } else if (field === 'input' || field === 'input_event') {
        } else {
            _.log.debug(deviceKey + ' : ' + topic + '  : ' + payload);
        }
    } else {
        _.log.error(deviceKey + ' topic not found');
    }
})();
