import {Observable} from "rxjs";
import ApiService, {DefaultApiService} from "./api-service";
import IDeviceList from "../model/device-config";
import IDictionary from "../model/dictionary";
import {first} from "rxjs/operators";

const DEVICE_CONFIG_API_PATH = 'devices'
const DEVICE_STATES_API_PATH = 'devices/status'
const DEVICE_API_PATH = 'device'

export default interface DeviceApiService extends ApiService {
    getDeviceConfig(): Observable<IDeviceList>;

    getDeviceStates(): Observable<any>;

    getDeviceState(deviceKey: string): Observable<any>;

    setDeviceCommand(deviceKey: string, value: IDictionary): Observable<boolean>;
}

export class DefaultDeviceApiService extends DefaultApiService implements DeviceApiService {
    getDeviceConfig(): Observable<IDeviceList> {
        return this.get<IDeviceList>(DEVICE_CONFIG_API_PATH);
    }

    getDeviceStates(): Observable<any> {
        return this.get<any>(DEVICE_STATES_API_PATH);
    }

    getDeviceState(deviceKey: string): Observable<any> {
        return this.get<any>(DEVICE_API_PATH + '/' + deviceKey);
    }

    setDeviceCommand(deviceKey: string, params: IDictionary): Observable<boolean> {
        const payload = Object.assign({}, params, {device: deviceKey})
        return new Observable((observer) => {
            this.post<any>(DEVICE_API_PATH, payload).pipe(first()).subscribe({
                next: val => observer.next(val),
                error: err => observer.error(err),
                complete: () => observer.complete()
            })
        });
    }

}
