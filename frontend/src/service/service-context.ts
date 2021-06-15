import ConfigService from "./config-service";
import DeviceService from "./device-service";
import EventService from "./event-service";
import DialogService from "./dialog-service";
import AuthService from "./auth-service";
import UserService from "./user-service";
import RenderService from "./render-service";
import RendererFactory from "../renderer/renderer-factory";

export interface Services {
    config(): ConfigService;

    device(): DeviceService;

    event(): EventService;

    dialog(): DialogService;

    auth(): AuthService;

    user(): UserService;

    render(): RenderService;
}

class ServiceContextImpl implements Services {

    private readonly _configService: ConfigService = new ConfigService();
    private readonly _EventService: EventService = new EventService();
    private readonly _deviceService: DeviceService = new DeviceService();
    private readonly _dialogService: DialogService = new DialogService();
    private readonly _authService: AuthService = new AuthService();
    private readonly _userService: UserService = new UserService();
    private readonly _renderService: RenderService = new RenderService();

    auth() {
        return this._authService;
    }

    config() {
        return this._configService;
    }

    device() {
        return this._deviceService;
    }

    event() {
        return this._EventService;
    }

    dialog() {
        return this._dialogService;
    }

    user() {
        return this._userService;
    }

    render() {
        return this._renderService;
    }

}

const ServiceContext = new ServiceContextImpl();
RendererFactory(ServiceContext.render().register.bind(ServiceContext.render()));
export default ServiceContext;

