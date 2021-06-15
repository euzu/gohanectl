import ApiService, {DefaultApiService} from "./api-service";
import {Observable} from "rxjs";
import ServerState from "../model/server-state";
import {IRoomConfig} from "../model/room";

const SERVER_STATE_API_PATH = 'status';
const ROOM_CONFIG_API_PATH = 'config/room'

export default interface ConfigApiService extends ApiService {
    getServerState(): Observable<ServerState>;

    getRoomConfig(): Observable<IRoomConfig>;
}

export class DefaultConfigApiService extends DefaultApiService implements ConfigApiService {
    getServerState(): Observable<ServerState> {
        return this.get<ServerState>(SERVER_STATE_API_PATH);
    }

    getRoomConfig(): Observable<IRoomConfig> {
        return this.get<IRoomConfig>(ROOM_CONFIG_API_PATH);
    }

}
