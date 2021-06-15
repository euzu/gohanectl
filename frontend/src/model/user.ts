export interface IUserSetting {
    key: string;
    value: string;
}

export default interface IUserSettings {
    settings: IUserSetting[];
}
