import React from "react";
// @ts-ignore
import createSvgIcon from "@material-ui/icons/utils/createSvgIcon";

const createIcon = (name: string, path: string): any => createSvgIcon(React.createElement("path", {
    d: path
}), name);

export default createIcon;