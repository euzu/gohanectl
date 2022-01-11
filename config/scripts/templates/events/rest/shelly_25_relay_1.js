var deviceKey = '{{.DeviceKey}}';
var payload = _.getPayload();
var data = JSON.parse(payload)
if (data) {
    function handleRelays(relaysData) {
        if (relaysData.length > 1) {
            var relayData = relaysData[1];
            var booleanFields = ["ison", "has_timer", "overpower", "overtemperature", "is_valid"];
            for (var i = 0, cnt = booleanFields.length; i < cnt; i++) {
                var booleanField = booleanFields[i];
                var booleanVal = relayData[booleanField];
                if (booleanVal != null) {
                    booleanVal = !!booleanVal;
                    if (booleanField === 'ison') {
                        _.setSharedMem(deviceKey, 'active', booleanVal);
                    } else {
                        _.setSharedMem(deviceKey, booleanField, booleanVal);
                    }
                }
            }

            var relayFields = ["timer_started", "timer_duration", "timer_remaining"];
            for (var j = 0, cnt2 = relayFields.length; j < cnt2; j++) {
                var relayField = relayFields[j];
                var relayVal = relayData[relayField];
                if (relayVal != null) {
                    _.setSharedMem(deviceKey, relayField, relayVal);
                }
            }
        }
    }

    var val = data['wifi_sta'];
    if (val != null) {
        _.setSharedMemNonNull(deviceKey, 'wifi_rssi', val['rssi']);
    }
    val = data['relays'];
    if (val != null) {
        handleRelays(val);
    } else {
        handleRelays(data);
    }

} else {
    _.log.warn(deviceKey + ' :  unexpected payload: ' + payload);
}
