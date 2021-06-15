(function () {
    var MIN_POWER = 13;
    var MAX_POWER = 150;
    var TIME_EXPIRE_FOR_STOP = 60000;
    var memKeyPower = '{{.DeviceKey}}-ender-power';
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
                _.sendMqtt('stat/ender/status', {print: 'done'});
                _.sendTelegram("Ender 3 fertig! " + _.formatDate(nowDate))
            }
            if (newValue < MIN_POWER) {
                if (startedAt) {
                    if (now - startedAt > TIME_EXPIRE_FOR_STOP) {
                        _.setSharedMem(memKeyPower, 'startedAt', 0);
                        _.sendMqtt('stat/ender/status', {print: 'done'});
                        _.sendTelegram("Ender 3 fertig! " + _.formatDate(nowDate))
                    }
                }
            }
        } else {
            if (newValue >= MAX_POWER) {
                _.setSharedMem(memKeyPower, 'startedAt', now);
                _.sendMqtt('stat/ender/status', {print: 'start'});
                _.sendTelegram("Ender 3 gestartet! " + _.formatDate(nowDate))
            }
        }
    }
})();
