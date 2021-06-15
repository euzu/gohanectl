import ConfigApiService, {DefaultConfigApiService} from "../api/config-api-service";
import {Observable, of} from "rxjs";
import ServerState from "../model/server-state";
import {IRoomConfig} from "../model/room";
import Utils from "../utils/utils";
import {first} from "rxjs/operators";

export default class ConfigService {
    private roomConfig : IRoomConfig;
    constructor(private configApiService: ConfigApiService = new DefaultConfigApiService()) {
    }

    getServerState(): Observable<ServerState> {
        return this.configApiService.getServerState();
    }

    getRoomConfig(): Observable<IRoomConfig> {
        if (Utils.isNil(this.roomConfig)) {
            return new Observable((observer) => {
                this.configApiService.getRoomConfig().pipe(first()).subscribe({
                    next: data => {
                        this.roomConfig = data;
                        observer.next(data);
                    },
                    error: err => observer.error(err),
                    complete: () => observer.complete()
                })
            });
        }
        return of(this.roomConfig);
    }
}
