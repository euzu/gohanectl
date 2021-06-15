import BatteryLowRenderer from "./battery-renderer";
import Renderer from "../model/renderer";

export default function RendererFactory(registrator: (name: string, renderer: Renderer) => void) {
    registrator('battery-low', new BatteryLowRenderer());
}