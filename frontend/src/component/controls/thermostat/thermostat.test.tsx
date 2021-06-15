import React from 'react';
import ReactDOM from 'react-dom';
import {act} from 'react-dom/test-utils';

import useTranslator from "../../../hook/use-translator";
import IDevice, {DeviceType} from "../../../model/device";
import Thermostat from "./thermostat";

jest.mock("../../../hook/use-translator");

let container: HTMLDivElement;

beforeEach(() => {
    container = document.createElement('div');
    document.body.appendChild(container);
});

afterEach(() => {
    document.body.removeChild(container);
    container = null;
});

describe('Thermostat', () => {
    it("render", () => {
        const device: IDevice = {
            type: DeviceType.sensor,
            caption: "Sensor",
            confirm: false,
            deviceKey: "test-device",
            expanded: false,
            icon: "",
            invert: {position: false, state_color: false},
            optimistic: false,
            room: "",
            supplemental: [],
            timeout: 0,
            url: "http://localhost/"
        };
        // @ts-ignore
        useTranslator.mockReturnValue((value: string) => value)
        act(() => {
            ReactDOM.render(<Thermostat  device={device}/>, container);
        });


        const devicePanel = container.querySelector('div');
        expect(devicePanel).toBeDefined();

        const thermostatPanel = devicePanel.querySelector('.thermostat-panel');
        expect(thermostatPanel).toBeDefined();
        const thermostat = thermostatPanel.querySelector('.thermostat');
        expect(thermostat).toBeDefined();

        expect(thermostat.querySelector('round-slider')).toBeDefined();
        expect(thermostat.querySelector('.thermostat-temp')).toBeDefined();
        expect(thermostat.querySelector('.thermostat-value')).toBeDefined();
    });
});
