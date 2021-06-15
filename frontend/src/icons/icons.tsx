import React from "react";
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpandLessIcon from '@material-ui/icons/ExpandLess';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import CallMadeIcon from '@material-ui/icons/CallMade';
import PauseCircleFilledIcon from "@material-ui/icons/PauseCircleFilled";
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import FormatQuoteIcon from '@material-ui/icons/FormatQuote';
import AllInclusiveIcon from '@material-ui/icons/AllInclusive';
import PowerSettingsNewIcon from '@material-ui/icons/PowerSettingsNew';
import LogoIcon from './logo';
import SocketIcon from './socket';
import ThermometerIcon from './thermometer';
import SensorIcon from './sensor';
import MqttIcon from './mqtt';
import BulbIcon from "./bulb";
import BulbUnknownIcon from "./bulb-unknown";
import SocketUnknownIcon from "./socket-unknown";
import ThermostatIcon from "./thermostat";
import RollerShutterClosedIcon from "./roller-shutter-closed";
import RollerShutterOpenIcon from "./roller-shutter-open";
import RollerShutter75Icon from "./roller-shutter-75";
import RollerShutter50Icon from "./roller-shutter-50";
import RollerShutter25Icon from "./roller-shutter-25";
import RollerShutterUnknownIcon from "./roller-shutter-unknown";
import UpIcon from "./up";
import DownIcon from "./down";
import ButtonIcon from "./button";

const ICON = {
    All: <AllInclusiveIcon/>,
    Quote: <FormatQuoteIcon/>,
    Logo: <LogoIcon/>,
    ArrowDown: <ExpandMoreIcon/>,
    ArrowUp: <ExpandLessIcon/>,
    Pause: <PauseCircleFilledIcon/>,
    Bulb: <BulbIcon/>,
    BulbUnknown: <BulbUnknownIcon/>,
    Socket: <SocketIcon/>,
    SocketUnknown: <SocketUnknownIcon/>,
    Thermometer: <ThermometerIcon/>,
    Sensor: <SensorIcon/>,
    Ping: <SettingsInputAntennaIcon/>,
    Mqtt: <MqttIcon/>,
    Link: <CallMadeIcon/>,
    Thermostat: <ThermostatIcon/>,
    RollerShutterClosed: <RollerShutterClosedIcon/>,
    RollerShutterOpen: <RollerShutterOpenIcon/>,
    RollerShutter75: <RollerShutter75Icon/>,
    RollerShutter50: <RollerShutter50Icon/>,
    RollerShutter25: <RollerShutter25Icon/>,
    RollerShutterUnknown: <RollerShutterUnknownIcon/>,
    Up: <UpIcon/>,
    Down: <DownIcon/>,
    Error: <ErrorOutlineIcon/>,
    Power: <PowerSettingsNewIcon/>,
    Button: <ButtonIcon/>
};

export function getIconByName(name: string): React.ReactElement {
    if (name) {
        // @ts-ignore
        return ICON[name];
    }
    return null;
}

export default ICON;
