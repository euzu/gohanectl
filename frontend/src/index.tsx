import React from 'react';
import ReactDOM from 'react-dom';
import 'react-app-polyfill/ie9';
import 'react-app-polyfill/ie11';
import 'typeface-roboto';
import './index.scss';
//import {noop} from 'rxjs';
// @ts-ignore
import {createBrowserHistory} from 'history';
//import reportWebVitals from "./reportWebVitals";
import App from "./component/app/app";

// @ts-ignore
export const history = createBrowserHistory({basename: process.env.PUBLIC_URL});


// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
// reportWebVitals(noop);

ReactDOM.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>,
    document.getElementById('root')
);
