import React, {useEffect, useState} from "react";
import {useServices} from "../../provider/service-provider";
import Login from "../login/login";
import {noop} from "rxjs";
import Spinner from "../progress/spinner";

interface AuthenticatorParams {
    factory() : React.ReactNode;
}

export default function Authenticator(params: AuthenticatorParams) {
    const {factory} = params;
    const services = useServices();
    const [progress, setProgress] = useState(true);
    const [authenticated, setAuthenticated] = useState(false);
    const [content, setContent] = useState(null);

    useEffect(() => {
            const sub = services.auth().authenticated().subscribe({
                next: noop,
                error: noop,
                complete: () => {
                    setProgress(val => false);
                }
            });
        return () => sub.unsubscribe();
    }, [services]);

    useEffect(() => {
        const subs = services.auth().isAuthenticated().subscribe((val: boolean) => {
            if(val && content === null) {
                setContent(() => factory());
            }
            setProgress(!val);
            setAuthenticated(() => val);
        });
        return () => subs.unsubscribe();
    }, [services, factory, content]);

    if (progress) {
         return <Spinner loading={progress}/>;
     }
    if (authenticated) {
        return  content;
    }
    return <Login/>
}
