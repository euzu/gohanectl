import React from "react";
import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@material-ui/core";

interface ConfirmDialogProps {
    title: string;
    message: string;
    onClose: () => void;
    onCancel?: () => void;
    onSubmit?: () => void;
    submitLabel?: string;
    cancelLabel?: string;
}

export default function ConfirmDialog(props: ConfirmDialogProps) {
    const {title, message, submitLabel, cancelLabel, onClose, onCancel, onSubmit} = props;

    const handleClose = () => {
        onCancel && onCancel();
        onClose();
    };
    const handleSubmit = () => {
        onSubmit && onSubmit();
        onClose();
    };

    return <Dialog
        open={true}
        onClose={handleClose}
        aria-labelledby="confirm-dialog-title"
        aria-describedby="confirm-dialog-description">
        <DialogTitle id="confirm-dialog-title">{title}</DialogTitle>
        <DialogContent>
            <DialogContentText id="confirm-dialog-description">
                {message}
            </DialogContentText>
        </DialogContent>
        <DialogActions>
            <Button onClick={handleClose} color="secondary">
                {cancelLabel || 'No'}
            </Button>
            <Button onClick={handleSubmit} color="primary" autoFocus>
                {submitLabel || 'Yes'}
            </Button>
        </DialogActions>
    </Dialog>
}
