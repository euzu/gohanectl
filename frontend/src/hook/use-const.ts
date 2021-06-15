import {useRef} from "react";
type ResultBox<T> = { v: T }

type ContantFunction<T> = () => T;

export default function useConst<T>(fnOrValue: T | ContantFunction<T>): T {
    const ref = useRef<ResultBox<T>>()
    if (!ref.current) {
        // @ts-ignore
        let value: T = typeof fnOrValue == 'function' ? fnOrValue() : fnOrValue as T;
        ref.current = { v: value }
    }
    return ref.current.v
}
