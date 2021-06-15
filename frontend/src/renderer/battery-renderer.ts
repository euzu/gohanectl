import Renderer from "../model/renderer";

export default class BatteryLowRenderer implements Renderer {
    render(value: unknown): unknown {
        return value ? 'change': 'ok';
    }
}