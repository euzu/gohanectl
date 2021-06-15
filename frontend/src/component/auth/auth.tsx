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
    const [children, setChildren] = useState(null);

    useEffect(() => {
        const sub = services.auth().authenticated().subscribe({next: noop,error: noop, complete:() => setProgress(val => false)});
        return () => sub.unsubscribe();
    }, [services]);

    useEffect(() => {
        const subs = services.auth().isAuthenticated().subscribe((val: boolean) => {
            setAuthenticated(() => val);
            if(val && children === null) {
                setChildren(() => factory());
            }
        });
        return () => subs.unsubscribe();
    }, [services, factory, children]);

    if (progress) {
         return <Spinner loading={progress}/>;
     }
    if (authenticated) {
        return  children;
    }
    return <Login/>
}
