import React from 'react';
import ReactDOM from 'react-dom';
import {act} from 'react-dom/test-utils';

import useTranslator from "../../../hook/use-translator";
import Sensor from "./sensor";
import IDevice, {DeviceType} from "../../../model/device";

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

describe('Sensor', () => {
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
            ReactDOM.render(<Sensor  device={device}/>, container);
        });

        expect(container.querySelector('div').classList.contains('device-panel')).toBe(true);
    });
});
