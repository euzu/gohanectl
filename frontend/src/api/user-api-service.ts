import {noop, Observable} from "rxjs";
import ApiService, {DefaultApiService} from "./api-service";
import IUserSettings from "../model/user";
import {first} from "rxjs/operators";

const USER_SETTINGS_API_PATH = 'user/setting'

export default interface UserApiService extends ApiService {
    getUserSettings(): Observable<IUserSettings>;

    saveUserSettings(settings: IUserSettings): void;
}

export class DefaultUserApiService extends DefaultApiService implements UserApiService {
    getUserSettings(): Observable<IUserSettings> {
        return this.get<IUserSettings>(USER_SETTINGS_API_PATH);
    }

    saveUserSettings(settings: IUserSettings): void {
        this.post<any>(USER_SETTINGS_API_PATH, settings).pipe(first()).subscribe({
            next: noop,
            error: noop,
        });
    }
}
