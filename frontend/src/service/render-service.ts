import Renderer from "../model/renderer";

type RendererDict = { [key: string]: Renderer };

export default class RenderService {
    private readonly renderer: RendererDict = {};

    register(name: string, renderer: Renderer) {
        this.renderer[name] = renderer;
    }

    render(renderer: string, value: unknown): unknown {
        const instance = this.renderer[renderer];
        if (instance) {
            return instance.render(value);
        }
        return value;
    }
}
