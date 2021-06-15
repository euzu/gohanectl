import {useTranslation} from "react-i18next";
import {useMemo} from "react";
import useConst from "./use-const";

export default function useTranslator(): (key: string, variables?: any) => string {
    const [t] = useConst(useTranslation('common'));
    return useMemo(() => (key: string, variables?: any) => t(key, variables), [t]);
}
