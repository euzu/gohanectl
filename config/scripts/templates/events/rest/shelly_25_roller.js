(function () {
    var deviceKey = '{{.DeviceKey}}';
    var payload = _.getPayload();
    var data = JSON.parse(payload)
    if (data) {
        function handleRollers(rollerData) {
            var rollerFields = ['state', 'power', 'overtemperature', 'stop_reason', 'last_direction', 'current_pos','positioning', 'safety_switch'];
            for (var i=0, cnt = rollerFields.length; i < cnt; i++) {
                var rollerField = rollerFields[i];
                var rollerVal = rollerData[rollerField];
                if (rollerVal != null) {
                    if (rollerField === 'current_pos') {
                        _.setSharedMem(deviceKey, 'active', rollerVal < 100);
                        _.setSharedMem(deviceKey, 'position', rollerVal);
                    } else if (rollerField === 'overtemperature') {
                        _.setSharedMem(deviceKey, 'overtemperature', !!rollerVal);
                    } else {
                        _.setSharedMem(deviceKey, rollerField, rollerVal);
                    }
                }
            }
        }

        var val = data['wifi_sta'];
        if (val != null) {
            _.setSharedMemNonNull(deviceKey, 'wifi_rssi', val['rssi']);
        }
        val = data['rollers'];
        if (val != null) {
            handleRollers(val);
        } else {
            handleRollers(data);
        }

    } else {
        _.log.warn(deviceKey + ' :  unexpected payload: ' + payload);
    }
})();
