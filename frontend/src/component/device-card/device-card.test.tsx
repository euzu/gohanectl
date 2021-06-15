import React from 'react';
import ReactDOM from 'react-dom';
import {act} from 'react-dom/test-utils';
import DeviceCard from "./device-card";
import IDevice, {DeviceType} from "../../model/device";
import useTranslator from "../../hook/use-translator";

jest.mock("../../hook/use-translator");

let container: HTMLDivElement;

beforeEach(() => {
    container = document.createElement('div');
    document.body.appendChild(container);
});

afterEach(() => {
    document.body.removeChild(container);
    container = null;
});

describe('DeviceCard', () => {
   it('Render Device', () => {
       // @ts-ignore
       useTranslator.mockReturnValue((value: string) => value)

       const device: IDevice = {
           type: DeviceType.light,
           caption: "Light",
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
       act(() => {
           ReactDOM.render(<DeviceCard device={device} hidden={false}/>, container);
       });
       const cardContainer = container.querySelector('div');
       const cardCaption = cardContainer.firstChild;
       const cardContent = cardContainer.childNodes.length > 1 && cardContainer.childNodes[1];
       expect(cardContainer?.classList.contains('device-card')).toBe(true);
       // @ts-ignore
       expect(cardCaption?.classList.contains('device-caption')).toBe(true);
       // @ts-ignore
       expect(cardContent?.classList.contains('device-card-content')).toBe(true);

       expect(cardCaption.firstChild?.textContent).toBe(device.caption);
       // @ts-ignore
       expect(cardCaption.firstChild?.href).toBe(device.url);

       const devicePanel = cardContent.firstChild;
       // @ts-ignore
       expect(devicePanel?.classList.contains('device-panel')).toBe(true);

       // @ts-ignore
       expect(devicePanel?.firstChild?.classList.contains('device')).toBe(true);

   });
    it('Render Device hidden', () => {
        // @ts-ignore
        useTranslator.mockReturnValue((value: string) => value)

        const device: IDevice = {
            type: DeviceType.light,
            caption: "Light",
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
        act(() => {
            ReactDOM.render(<DeviceCard device={device} hidden={true}/>, container);
        });
        const cardContainer = container.querySelector('div');
        expect(cardContainer?.classList.contains('device-card')).toBe(true);
        expect(cardContainer?.classList.contains('hidden')).toBe(true);

    })
});