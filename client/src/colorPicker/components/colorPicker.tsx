import React, { FC } from "react";
import { SketchPicker } from "react-color";

interface Props {
  colorPicked: () => void; 
}
// http://casesandberg.github.io/react-color/

const ColorPicker: FC<Props> = ({colorPicked}) => {
  return (
    <SketchPicker 
      onChangeComplete={colorPicked}      
    />
  );
};

export default ColorPicker;
