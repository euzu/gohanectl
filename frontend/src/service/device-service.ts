import {Observable, of} from "rxjs";
import DeviceApiService, {DefaultDeviceApiService} from "../api/device-api-service";
import IDeviceConfig from "../model/device-config";
import IDevice from "../model/device";
import IDictionary from "../model/dictionary";
import Utils from "../utils/utils";
import {first} from "rxjs/operators";

export default class DeviceService {

    private deviceConfig: IDeviceConfig;

    constructor(private deviceApiService: DeviceApiService = new DefaultDeviceApiService()) {
    }

    getDeviceConfig(): Observable<IDeviceConfig> {
        if (Utils.isNil(this.deviceConfig)) {
            return new Observable((observer) => {
                this.deviceApiService.getDeviceConfig().pipe(first()).subscribe({
                    next: data => {
                        this.deviceConfig = data;
                        observer.next(data);
                    },
                    error: err => observer.error(err),
                    complete: () => observer.complete()
                })
            });
        }
        return of(this.deviceConfig);
    }

    getDeviceState(deviceKey: string): Observable<any> {
        return this.deviceApiService.getDeviceState(deviceKey);
    }

    getDeviceStates(): Observable<any> {
        return this.deviceApiService.getDeviceStates();
    }

    setDevicePower(deviceKey: string, on: boolean): Observable<boolean> {
        return this.deviceApiService.setDeviceCommand(deviceKey, {command: 'power', payload: on});
    }

    setDeviceTemperature(deviceKey: string, temperature: number): Observable<boolean> {
        return this.setDeviceCommand(deviceKey, {command: 'temperature', payload: temperature, temperature});
    }

    setDevicePosition(deviceKey: string, position: number): Observable<boolean> {
        return this.setDeviceCommand(deviceKey, {command: 'position', payload: position, position});
    }

    setDevicePause(deviceKey: string): Observable<boolean> {
        return this.setDeviceCommand(deviceKey, {command: 'pause', payload: true, pause: true});
    }

    setDeviceCommand(deviceKey: string, payload: IDictionary): Observable<boolean> {
        return this.deviceApiService.setDeviceCommand(deviceKey, payload);
    }

    isDeviceTimeout(device: IDevice, deviceState: any): boolean {
        // Timeout below 0 disables timeout check.
        // Devices with timeout below 0 are for example Contact Sensors.
        // They are only triggered mechanically to save battery.
        if (device.timeout < 0) {
            return false;
        }
        if (deviceState?.lastUpdated) {
            const elapsedTime = Date.now() - deviceState.lastUpdated;
            if (elapsedTime < (device.timeout * 1000)) {
                return false
            }
        }
        return true;
    }
}
