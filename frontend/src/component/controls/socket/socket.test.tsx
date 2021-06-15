import React from 'react';
import ReactDOM from 'react-dom';
import {act} from 'react-dom/test-utils';

import useTranslator from "../../../hook/use-translator";
import IDevice, {DeviceType} from "../../../model/device";
import Socket from "./socket";

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

describe('Socket', () => {
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
            ReactDOM.render(<Socket  device={device}/>, container);
        });


        const devicePanel = container.querySelector('div');
        expect(devicePanel).toBeDefined();

        const expander = devicePanel.querySelector('.expander');
        expect(expander).toBeDefined();

        const deviceStatePanel = devicePanel.querySelector('.device-state-panel');
        expect(deviceStatePanel).toBeDefined();
    });
});
