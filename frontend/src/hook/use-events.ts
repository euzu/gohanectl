import {useEffect, DependencyList} from "react";
import {EventListener} from "../service/event-service";
import {useServices} from "../provider/service-provider";

export default function useEvents(factory: () => EventListener[], deps?: DependencyList) {
    const services = useServices();
    useEffect(() => {
        const eventListener = factory();
        const subscription = services.event().addEventListeners(eventListener);
        return () => subscription?.unsubscribe();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [services, ...(deps || [])]);
}
