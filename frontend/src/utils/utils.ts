
const Utils = {
    isNil: (v: any) : boolean => v == null,
    isBlank: (v : string): boolean => v == null || !v.trim().length
}
export default Utils;