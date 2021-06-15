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
        } else if (field === 'relay') {
            var relayField = topicParts[3];
            var relayIndex = parseInt(relayField);
            if (typeof relayIndex === 'number') {
                if (relayIndex === 0) {
                    if (topicParts.length > 4) {
                        relayField = topicParts[4];
                        if (relayField === "energy" || relayField === "power") {
                            _.setSharedMem(deviceKey, 'relay-' + relayIndex + '-' + relayField, parseFloat(payload));
                        } else {
                            _.log.debug(deviceKey + ' : ' + topic + '  : unknown relayField : ' + payload);
                        }
                    } else {
                        _.setSharedMem(deviceKey, 'active', payload === 'on');
                    }
                }
            } else {
                _.log.debug(deviceKey + ' : ' + topic + '  : ' + payload);
            }
        } else if (field === 'temperature_f') {
        } else if (field === 'roller') {
        } else if (field === 'input') {
        } else if (field === 'input_event') {
        } else {
            _.log.debug(deviceKey + ' : ' + topic + '  : ' + payload);
        }
    } else {
        _.log.error(deviceKey + ' topic not found');
    }
})();
