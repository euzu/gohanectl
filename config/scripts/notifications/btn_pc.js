(function () {
    var key = _.getParamKey();
    if (key === "action") {
        var value = _.getParamNewValue();
        if (value === 'single') {
            _.sendMqtt('z2m/socket_pc/set', '{"state": "on"}');
        } else if(value === 'double') {
        } else if(value === 'triple') {
        } else if(value === 'quadruple') {
        } else if(value === 'long') {
        } else if(value === 'long_release') {
        } else if(value === 'many') {
        }
    }
})();
