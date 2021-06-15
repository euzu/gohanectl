import React from 'react';
import ReactDOM from 'react-dom';
import {act} from 'react-dom/test-utils';

import {useServices} from "../../provider/service-provider";
import Authenticator from './auth'
import {Observable, of} from "rxjs";

import useTranslator from "../../hook/use-translator";

jest.mock("../../provider/service-provider");
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

describe('Auth', () => {
    it("unauthenticated", () => {
        const servicesMock = {
            auth: () => {
                return {
                    authenticated: () => of(false),
                    isAuthenticated: () => of(false)
                }
            },
        }
        // @ts-ignore
        useServices.mockReturnValue(servicesMock);
        // @ts-ignore
        useTranslator.mockReturnValue((value: string) => value)
        act(() => {
            ReactDOM.render(<Authenticator factory={() => <React.Fragment/>}/>, container);
        });

        expect(container.querySelector('div').classList.contains('login')).toBe(true);
    });

    it("authenticated", () => {
        const servicesMock = {
            auth: () => {
                return {
                    authenticated: () => of(true),
                    isAuthenticated: () => of(true)
                }
            },
        }
        // @ts-ignore
        useServices.mockReturnValue(servicesMock);
        // @ts-ignore
        useTranslator.mockReturnValue((value: string) => value)
        act(() => {
            ReactDOM.render(<Authenticator factory={() => <div id={'test-content'}/>}/>, container);
        });

        expect(container.querySelector('#test-content')).toBeDefined();
    });

    it("waiting", () => {
        const servicesMock = {
            auth: () => {
                return {
                    authenticated: () => new Observable(),
                    isAuthenticated: () => of(false)
                }
            },
        }
        // @ts-ignore
        useServices.mockReturnValue(servicesMock);
        // @ts-ignore
        useTranslator.mockReturnValue((value: string) => value)
        act(() => {
            ReactDOM.render(<Authenticator factory={() => <div id={'test-content'}/>}/>, container);
        });
        expect(container.querySelector('div').classList.contains('loader')).toBe(true);
    });
});