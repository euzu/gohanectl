import React, {useRef, useState} from "react";
import "./login.scss"
import {Button, Container, Grid, TextField} from "@material-ui/core";
import {useServices} from "../../provider/service-provider";
import {first} from "rxjs/operators";
import useTranslator from "../../hook/use-translator";
import ICON from "../../icons/icons";

interface LoginProperties {
}

export default function Login(props: LoginProperties) {
    const services = useServices();
    const translate = useTranslator();
    const usernameInputRef = useRef<HTMLInputElement>();
    const passwordInputRef = useRef<HTMLInputElement>();
    const [success, setSuccess] = useState<boolean>(true)

    const login = () => {
        const username = usernameInputRef.current.value;
        const password = passwordInputRef.current.value;
        if (username.trim().length > 0 && password.trim().length > 0) {
            services.auth().login(username, password).pipe(first()).subscribe({
                next: (value: boolean) => setSuccess(value),
                error: () => setSuccess(false)
            });
        }
    };

    return <div className={"login"}>
        <Container className={"login-container"} maxWidth="xs">
            <form>
                <Grid container spacing={3}>
                    <Grid item xs={12}>
                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <TextField fullWidth label="Username" name="username" size="small" variant="outlined"
                                           inputRef={ref => {
                                               if (ref) {
                                                   usernameInputRef.current = ref;
                                               }
                                           }}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    label="Password"
                                    name="password"
                                    size="small"
                                    type="password"
                                    variant="outlined"
                                    inputRef={ref => {
                                        if (ref) {
                                            passwordInputRef.current = ref;
                                        }
                                    }}
                                />
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item xs={12}>
                        <Button color="primary" fullWidth type="submit" variant="contained" onClick={(evt: any) => {
                            evt.preventDefault();
                            login();
                        }}>
                            Log in
                        </Button>
                    </Grid>
                    <Grid item xs={12}>
                        <div className={"alert" + (success ? ' alert-hidden' : '')}><div className={"alert-icon"}>{ICON.Error}</div>{translate('text.login-failed')}</div>
                    </Grid>
                </Grid>
            </form>
        </Container>
    </div>
}
