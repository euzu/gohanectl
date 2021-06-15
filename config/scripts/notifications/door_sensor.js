(function () {
    if (_.getParamKey() === "active") {
        if (_.getParamNewValue() == false) {
            _.sendTelegram("eshickten gechildi! " + _.formatDate(new Date()));
        }
    }
})();