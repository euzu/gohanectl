(function () {
    var deviceKey = '{{.DeviceKey}}';
    var topic = _.getTopic();
    var payload = _.getPayload();
    var parts = _.getSplittedTopic();
    if (parts.length) {
        var field = parts[2];
        if (field === 'power' || field === 'power1') {
            _.setSharedMem(deviceKey, 'active', payload === 'on' || payload === 1);
        } else if (field === 'result') {
            var result = JSON.parse(payload)
            if (result) {
                var power = result["power1"] || result["power"]
                _.setSharedMem(deviceKey, 'active', power === 'on' || power === 1);
            }
        } else if (field === 'status') {
            var state = _.parseJson(payload);
            if (state) {
                var status = state['status'];
                if (status) {
                    var powerState = (status['power'] || status['power1'])
                    _.setSharedMem(deviceKey, 'active',  powerState === 1 || powerState === 'on');
                }
            }
        } else {
            _.log.debug(deviceKey + ': ' + topic + ' :   ' + payload);
        }
    } else {
        _.log.error(deviceKey + ' topic not found');
    }
})();
