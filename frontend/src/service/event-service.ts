import {Subscription, Subject} from "rxjs";

export enum EventTopic {
    DeviceState = 'hanectl-state',
    Websocket = 'hanectl-websocket',
    Mqtt = 'hanectl-mqtt',
    PreventSwipe = 'hanectl-prevent-swipe'
}

export interface AppEvent {
    topic: EventTopic;
    payload: any;
}

type EventHandlerCallback = (e: AppEvent) => void;

export type EventListener = { topic: EventTopic, handler: EventHandlerCallback};

export default class EventService {
    private topics: any = {};

    handleEvents(e: AppEvent) {
        const topic = e.topic;
        if (topic) {
            const currentTopic: Subject<AppEvent> = this.topics[topic];
            if (currentTopic) {
                currentTopic.next(e);
                return true;
            }
        }
        return false;
    }

    addEventListeners(listeners: EventListener[]): Subscription | undefined {
        const subscription = new Subscription();
        for (let i = 0, cnt = listeners.length; i < cnt; i++) {
            this.addEventListener(listeners[i], subscription);
        }
        return subscription;
    }

    addEventListener(listener: EventListener, sub?: Subscription): Subscription | undefined {
        if (!listener) {
            return undefined;
        }
        let subject: Subject<AppEvent> = this.topics[listener.topic];
        if (!subject) {
            subject = new Subject<AppEvent>();
            this.topics[listener.topic] = subject;
        }
        const subscription = sub || new Subscription();
        subscription.add(subject.subscribe(e => listener.handler(e)));
        return subscription;
    }

    fireEvent<T>(topic: EventTopic, payload: T = null) {
        const message = {topic, payload};
        this.handleEvents(message);
    }

    close() {
    }

}
