var deviceKey = '{{.DeviceKey}}';
var payload = _.getPayload();

var data = JSON.parse(payload)
if (data) {
    var val = data['ison'];
    if (val != null) {
        _.setSharedMem(deviceKey, 'active', val === true);
    }
    val = data['wifi_sta'];
    if (val != null) {
        _.setSharedMemNonNull(deviceKey, 'wifi_rssi', val['rssi']);
    }
} else {
    _.log.warn(deviceKey + ' :  unexpected payload: ' + payload);
}
