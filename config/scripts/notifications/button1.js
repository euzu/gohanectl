(function () {
    if (_.getParamKey() === "click") {
        var value = _.getParamNewValue();
        if (value === 'single') {
            _.sendMqtt('cmnd/light_wz/Power', "toggle");
        } else if(value === 'double') {
        } else if(value === 'triple') {
        } else if(value === 'quadruple') {
        } else if(value === 'long') {
        } else if(value === 'long_release') {
        } else if(value === 'many') {
        }
    }
})();
