import config from "../config";
import {useServices} from "../provider/service-provider";
import {Services} from "../service/service-context";
import {EventTopic} from "../service/event-service";

const PING_INTERVAL = 30 * 1000;
const JSON_RECORD_SEPARATOR = String.fromCharCode(0x1E);

export class Websocket {
    private connected: boolean;
    private ws: WebSocket;
    private services: Services;
    private tsLastPing: number;

    constructor() {
        this.connect();
        const checkPing = () => {
            if ((performance.now() - this.tsLastPing) > PING_INTERVAL) {
                this.services.event().fireEvent(EventTopic.Websocket, {error: 'ping timeout', type: 'timeout'});
            }
        }
        setInterval(checkPing, PING_INTERVAL * 2);
    }

    public isConnected(): boolean {
        return this.connected;
    }

    private connect() {
        const ws = new WebSocket(config.api.wsUrl);
        ws.onopen = () => {
            this.connected = true;
            this.services.event().fireEvent(EventTopic.Websocket, {connected: true});
        }
        ws.onclose = () => {
            this.connected = false;
            this.services.event().fireEvent(EventTopic.Websocket, {connected: false});
            const reconnect = () => this.connect();
            setTimeout(reconnect, 1000);
        }
        ws.onerror = (error: Event) => {
            ws.close();
            this.handleError(error);
        }
        ws.onmessage = (msg: MessageEvent) => this.handleMessage(msg.data);
    }

    public send(data: string) {
        this.ws.send(data);
    }

    private handleMessage(data: any) {
        let payload = data;
        if (typeof data === 'string') {
            const msgs = data.split(JSON_RECORD_SEPARATOR)
            for (let i=0, cnt = msgs.length; i < cnt; i++) {
                try {
                    const message = msgs[i].trim();
                    if (message.length > 0) {
                        payload = JSON.parse(message);
                    }
                } catch (err) {
                    console.log("failed to parse json:", data);
                }
                if (payload?.ping) {
                    this.tsLastPing = performance.now();
                    this.services.event().fireEvent(EventTopic.Websocket, {ping: this.tsLastPing});
                } else if (payload?.mqtt) {
                    this.services.event().fireEvent(EventTopic.Mqtt, {mqtt: payload.mqtt === 1});
                } else if (payload?.deviceKey && this.services) {
                    this.services.event().fireEvent(EventTopic.DeviceState, payload);
                }
            }
        }
    }

    private handleError(error: Event) {
        this.services.event().fireEvent(EventTopic.Websocket, {error: error, type: 'general'});
    }

    public setServices(services: Services) {
        if (this.services !== services) {
            this.services = services;
        }
    }
}

const websocket = new Websocket();

export default function useWebsocket(): Websocket {
    const services = useServices();
    websocket.setServices(services);
    return websocket;
}
