import React from "react";
import "./timeout-indicator.scss";

interface TimeoutIndicatorParams {
    shiftX?: number;
    shiftY?: number;
}

export default function TimeoutIndicator(params: TimeoutIndicatorParams) {
    const {shiftX, shiftY} = params;
    const myStyle = {
        transform: 'translate3d(' + (shiftX || 0) + 'px,' + (shiftY || 0) + 'px,0)'
    };
    return <div className={"timeout-indicator-container"}>
        <div className={"fade-in timeout-indicator"} style={myStyle}/>
    </div>;
}