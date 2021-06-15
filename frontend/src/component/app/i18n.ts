import i18next from "i18next";
import {initReactI18next} from "react-i18next";
import common_en from "../../i18n/en_common.json";
import common_de from "../../i18n/de_common.json";

/* to change the language use
        i18n.changeLanguage(language);
 */

export default function i18n_init() {

    i18next.use(initReactI18next)
    i18next.init({
        interpolation: {escapeValue: false},  // React already does escaping
        // @ts-ignore
        lng: navigator?.language || navigator?.userLanguage || 'en',                              // language to use
        resources: {
            en: {
                common: common_en               // 'common' is our custom namespace
            },
            de: {
                common: common_de
            },
        },
    }).catch((err: any) => console.log('failed to load i18n'));
}
