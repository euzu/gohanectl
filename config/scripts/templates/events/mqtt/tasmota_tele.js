var deviceKey = '{{.DeviceKey}}';
var topic = _.getTopic();
var payload = _.getPayload();
var parts = _.getSplittedTopic();
if (parts && parts.length) {
    var field = parts[2];
    if (field === 'lwt') {
        _.setSharedMem(deviceKey, 'online', "online" === payload);
    } else if (field === 'state') {
        var state = _.parseJson(payload);
        var powerState = state['power1'] || state['power'];
        _.setSharedMem(deviceKey, 'active', powerState === 'on' || powerState === 1);

        var wifi = state['wifi'];
        if (wifi) {
            _.setSharedMemNonNull(deviceKey, 'wifi_signal', wifi['signal']);
            _.setSharedMemNonNull(deviceKey, 'wifi_rssi', wifi['rssi']);
        }
    } else if (field === 'sensor') {
        var sensor = _.parseJson(payload);
        if (sensor) {
            var energy = sensor['energy'];
            if (energy) {
                function setEnergy(key) {
                    var value = energy[key];
                    if (typeof value === 'number' && isFinite(value)) {
                        _.setSharedMem(deviceKey, key, value);
                    }
                }

                var power = energy["power"];
                if (power < 2) {
                    power = 0;
                }
                _.setSharedMem(deviceKey, 'power', power);

                setEnergy('totalstarttime');
                setEnergy('total');
                setEnergy('yesterday');
                setEnergy('today');
                setEnergy('apparentpower');
                setEnergy('reactivepower');
                setEnergy('factor');
                setEnergy('voltage');
                setEnergy('current');
            }
        }
    } else {
        _.log.debug(' tele-' + topic + ':   ' + field);
    }
} else {
    _.log.error(deviceKey + ' topic not found');
}
