import React, { FC, useState } from "react";
import { SketchPicker } from "react-color";
import { Button } from "antd";
import { Redirect } from "react-router-dom";

interface Props {
}
// http://casesandberg.github.io/react-color/

const ColorPicker: FC<Props> = ({}) => {
  let chosenColor = null; 
  const [hasChosenColor, setHasChosenColor] = useState(
    false
  ); 

  // Used for when a new color is selected on the color picker.
  const handleColorChange = (color) => {
    chosenColor = color;  
  };

  // Okay pressed.
  const colorChosen = () => {
    if (chosenColor != null) {
      console.log(chosenColor); 
      setHasChosenColor(true); 
    } else {
      // TODO: display message somewhere saying that you have to pick a color first
      // (which is not ideal. I'm not sure why this component works like that.)
    }
  }

  // Cancel pressed.
  const cancel = () => {
    setHasChosenColor(true); 
  }

  return ( 
    hasChosenColor ? <Redirect to='/'/> : 
    <div>
      <SketchPicker 
        onChangeComplete={handleColorChange}      
      />
      <Button onClick={colorChosen}>
        Okay
      </Button>
      <Button onClick={cancel}>
        Cancel
      </Button>
    </div>
  );
};



export default ColorPicker;
