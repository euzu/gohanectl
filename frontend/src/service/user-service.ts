import UserApiService, {DefaultUserApiService} from "../api/user-api-service";
import {Observable, ReplaySubject} from "rxjs";
import IUserSettings from "../model/user";
import Utils from "../utils/utils";
import {first} from "rxjs/operators";

export default class UserService {
    private loaded: boolean = false;
    private userSettings$: ReplaySubject<IUserSettings> = new ReplaySubject<IUserSettings>(null);

    constructor(private userApiService: UserApiService = new DefaultUserApiService()) {
    }

    private findSettingByKey(userSettings: IUserSettings, key: string): string {
        if (!Utils.isNil(userSettings)) {
            for (const setting of userSettings.settings) {
                if (setting.key === key) {
                    return setting.value;
                }
            }
        }
        return null;
    }

    private loadSettings() {
        if (!this.loaded) {
            this.loaded = true;
            this.userApiService.getUserSettings().pipe(first()).subscribe({
                next: data => this.userSettings$.next(data),
                error: (err) => this.userSettings$.error(err),
                complete: () => this.userSettings$.complete()
            })
        }
    }

    getUserSetting(key: string): Observable<string> {
        this.loadSettings();
        return new Observable((observer) => {
            this.userSettings$.pipe(first()).subscribe({
                next: (data) => {
                    const value = this.findSettingByKey(data, key);
                    if (value !== null) {
                        observer.next(value);
                    } else {
                        observer.next(null);
                    }
                },
                error: (err) => observer.error(err),
                complete: () => observer.complete()
            });
        });
    }

    getUserSettings(): Observable<IUserSettings> {
        return this.userSettings$;
    }

    saveUserSettings(userSettings: IUserSettings): void {
        this.userApiService.saveUserSettings(userSettings)
    }

}
