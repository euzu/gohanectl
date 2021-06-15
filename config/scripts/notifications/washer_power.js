(function () {
    var memKeyPower = '{{.DeviceKey}}-washer-power';
    var key = _.getParamKey();
    if (key === "power") {
        var newValue = _.getParamNewValue();
        //var oldValue = _.getParamOldValue();

        var nowDate = new Date();
        var now = nowDate.getTime();

        var startedAt = _.getSharedMem(memKeyPower, 'startedAt');
        if (startedAt) {
            if (newValue === 0) {
                _.setSharedMem(memKeyPower, 'startedAt', 0);
                _.sendTelegram("Waschmaschine fertig! " + _.formatDate(nowDate))
            }
            if (newValue < 5) {
                if (startedAt) {
                    if (now - startedAt > 60000) {
                        _.setSharedMem(memKeyPower, 'startedAt', 0);
                        _.sendTelegram("Waschmaschine fertig! " + _.formatDate(nowDate))
                    }
                }
            }
        } else {
            if (newValue >= 50) {
                _.setSharedMem(memKeyPower, 'startedAt', now);
            }
        }
    }
})();
