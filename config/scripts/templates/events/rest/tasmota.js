var deviceKey = '{{.DeviceKey}}';
var payload = _.getPayload();

var data = JSON.parse(payload)
if (data) {
    var power = data["power"];
    if (power == null) {
        var status = data["status"];
        if (status) {
            power = status["power"];
        }
    }
    if (power != null) {
        _.setSharedMem(deviceKey, 'active', power === "on" || power === 1);
    }

    var sensor = data["statussns"];
    if (sensor && sensor['energy']) {
        var energy = sensor['energy'];

        function setEnergy(key) {
            var value = energy[key];
            if (typeof value === 'number' && isFinite(value)) {
                _.setSharedMem(deviceKey, key, value);
            }
        }

        var epower = energy["power"];
        if (epower < 2) {
            epower = 0;
        }
        _.setSharedMem(deviceKey, 'power', epower);

        setEnergy('total');
        setEnergy('yesterday');
        setEnergy('today');
        setEnergy('apparentpower');
        setEnergy('reactivepower');
        setEnergy('factor');
        setEnergy('voltage');
        setEnergy('current');
    }
    var statussts = data["statussts"];
    if (statussts && statussts['wifi']) {
        var wifi = statussts['wifi'];
        _.setSharedMemNonNull(deviceKey, 'wifi_signal', wifi['signal']);
        _.setSharedMemNonNull(deviceKey, 'wifi_rssi', wifi['rssi']);
    }
} else {
    _.log.debug(deviceKey + ' :  ' + payload);
}
