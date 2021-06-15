import React from "react";
import {render, unmountComponentAtNode} from "react-dom";
import ConfirmDialog from "../component/dialog/confirm-dialog";

interface ConfirmDialogProps {
    title: string;
    message: string;
    onCancel?: () => void;
    onSubmit?: () => void;
    successLabel?: string;
    cancelLabel?: string;
}

function createContainer() : HTMLDivElement {
    const containerElement = document.createElement('div');
    document.body.appendChild(containerElement);
    return containerElement;
}

export default class DialogService {

    confirm(props: ConfirmDialogProps) {
        const containerElement = createContainer();
        const {title, message, successLabel, cancelLabel, onCancel, onSubmit} = props;
        render(<ConfirmDialog title={title} message={message}
                              submitLabel={successLabel}
                              cancelLabel={cancelLabel}
                                     onClose={() => unmountComponentAtNode(containerElement)}
                                     onCancel={onCancel}
                                     onSubmit={onSubmit}/>, containerElement);
    }
}
