(function(){
    var deviceKey = '{{.DeviceKey}}';
    var topic = _.getTopic();
    var payload = _.getPayload();
    var topicParts = _.getSplittedTopic();
    if (topicParts.length) {
        var field = topicParts[2];
        if (field === 'overtemperature') {
            _.setSharedMem(deviceKey, field, !!parseFloat(payload));
        } else if (field === 'temperature') {
            _.setSharedMem(deviceKey, field, parseFloat(payload));
        } else if (field === 'roller') {
            if (topicParts.length > 4) {
                var rollerField = topicParts[4];
                if (rollerField === 'energy' || rollerField === 'power') {
                    _.setSharedMem(deviceKey, rollerField, parseFloat(payload));
                } else if (rollerField === 'pos') {
                    var position = parseInt(payload)
                    _.setSharedMem(deviceKey, 'position', position);
                    _.setSharedMem(deviceKey, 'active', position < 100);
                } else if (rollerField === 'stop_reason') {
                    _.setSharedMem(deviceKey, rollerField, payload);
                } else {
                    _.log.debug(deviceKey + ' : roller-' + topic + '  : ' + payload);
                }
            }
        } else if (field === 'input') {
            // var input = parseInt(topicParts[3]);
            // if (input === 0 || input === 1) {
            //     _.setSharedMem(deviceKey, 'input-' + input, parseFloat(payload));
            // }
        } else if (field === 'input_event') {
        } else if (field === 'temperature_f') {
        } else if (field === 'relay') {
        } else {
            _.log.debug(deviceKey + ' : ' + topic + '  : ' + payload);
        }
    } else {
        _.log.error(deviceKey + ' topic not found');
    }
})();
