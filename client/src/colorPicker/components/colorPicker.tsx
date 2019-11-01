import React, { FC, useState } from "react";
import { SketchPicker } from "react-color";
import { Button } from "antd";
import { Redirect } from "react-router-dom";
import './colorPicker.scss';

interface RGBColor {
  r: number;
  g: number;
  b: number;
}

interface Props {
  onColorChange: (color: RGBColor) => void; 
  onComplete: () => void;
  onCancel: () => void;
}
// http://casesandberg.github.io/react-color/

const ColorPicker: FC<Props> = ({ onColorChange, onComplete, onCancel }) => {
  const [color, setColor] = useState({ r: 0, g: 0, b: 0 })
  const complete = () => onComplete();

  const onColorSelection = (color) => {
    onColorChange(color); 
  }

  return ( 
    <div>
      <SketchPicker 
        color={color}
        onChange={(c) => onColorSelection({ r: c.rgb.r, b: c.rgb.b, g: c.rgb.g })}
      />
      <div className='button-bar'>
        <Button onClick={complete} className='okay-button'>
          Okay
        </Button>
        <Button onClick={onCancel} className='cancel-button'>
          Cancel
        </Button>
      </div>
    </div>
  );
};



export default ColorPicker;
