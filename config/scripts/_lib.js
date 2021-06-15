function stdlib() {
    var log = {
        debug: function (message) {
            _otto_log('debug', message);
        },
        info: function (message) {
            _otto_log('info', message);
        },
        warn: function (message) {
            _otto_log('warn', message);
        },
        error: function (message) {
            _otto_log('error', message);
        },
    }

    function getParam(key) {
        return _otto_getParam(key, _otto_params);
    }

    function getSharedMem(deviceKey, key) {
        return _otto_getMem(deviceKey, key);
    }

    function setSharedMem(deviceKey, key, value) {
        _otto_setMem(deviceKey, key, value);
    }
    function setSharedMemNonNull(deviceKey, key, value) {
        if (value != null) {
            setSharedMem(deviceKey, key, value);
        }
    }

    function sendTelegram(message) {
        _otto_telegram(message);
    }

    function sendMqtt(topic, payload) {
        if (typeof payload !== "string") {
            _otto_mqtt(topic, JSON.stringify(payload));
        } else {
            _otto_mqtt(topic, payload);
        }
    }

    function startsWith(str, prefix) {
        return str != null && str.slice(0, prefix.length) === prefix;
    }

    function endsWith(str, suffix) {
        return str != null && str.slice(-suffix.length) === suffix;
    }

    function parseJson(json) {
        try {
            return JSON.parse(json);
        } catch (e) {
            _otto_log('error', "Failed to parse json:" + json);
            return json;
        }
    }

    function getTopic() {
        return getParam("topic");
    }

    function getSplittedTopic() {
        return getParam("topic") ? getParam("topic").split('/') : [];
    }

    function getPayload() {
        return getParam("payload");
    }

    function getParamKey() {
        return getParam("key");
    }

    function getParamNewValue() {
        return getParam("newValue");
    }

    function getParamOldValue() {
        return getParam("oldValue");
    }

    function formatDate(date) {
        if (date) {
            return date.toLocaleString().replace(',', '');
        }
        return '';
    }

    return {
        log: log,
        getParam: getParam,
        getSharedMem: getSharedMem,
        setSharedMem: setSharedMem,
        setSharedMemNonNull: setSharedMemNonNull,
        sendTelegram: sendTelegram,
        sendMqtt: sendMqtt,
        startsWith: startsWith,
        endsWith: endsWith,
        parseJson: parseJson,
        getTopic: getTopic,
        getSplittedTopic: getSplittedTopic,
        getPayload: getPayload,
        getParamKey: getParamKey,
        getParamNewValue: getParamNewValue,
        getParamOldValue: getParamOldValue,
        formatDate: formatDate
    }
}

var _ = stdlib();