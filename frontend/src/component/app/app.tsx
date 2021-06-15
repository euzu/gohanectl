import React from "react";
import './app.scss';
import {ServiceProvider} from "../../provider/service-provider";
import Devices from "../devices/devices";
import useWebsocket from "../../hook/use-websocket";
import Authenticator from "../auth/auth";
import {I18nextProvider} from "react-i18next";
import i18next from "i18next";
import i18n_init from "./i18n";
import Header from "../header/header";

i18n_init();

const factory = () => <Devices/>;

function App() {
    useWebsocket();
    return (
        <I18nextProvider i18n={i18next}>
            <ServiceProvider>
                <div className="app">
                    <Header/>
                    <div className="app-content">
                        <Authenticator factory={factory}/>
                    </div>
                </div>
            </ServiceProvider>
        </I18nextProvider>
    );
}

export default App;
