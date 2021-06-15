import React from "react";
import './spinner.scss';

export default function Spinner(props: any) {
    if (props.loading) {
        return (
            <div className='loader'>
                <span/>
                <span/>
                <span/>
            </div>
        )
    }
    return null;
}
