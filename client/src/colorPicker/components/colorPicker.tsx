import React, { FC, useState } from "react";
import { SketchPicker } from "react-color";
import { Button } from "antd";
import { Redirect } from "react-router-dom";

interface Props {
}
// http://casesandberg.github.io/react-color/

const ColorPicker: FC<Props> = ({}) => {
  const [chosenColor, setChosenColor] = useState(
    ""
  ); 
  const [hasChosenColor, setHasChosenColor] = useState(
    false
  ); 

  // Used for when a new color is selected on the color picker.
  const handleColorChange = (color) => {
    setChosenColor(color); 
    console.log(chosenColor); 
  };

  // Okay pressed.
  const colorChosen = () => {
    if (chosenColor != "") {
      setHasChosenColor(true); 
    } else {
      // TODO: display message somewhere saying that you have to pick a color first
      // (which is not ideal. I'm not sure why this component works like that.)
    }
  }

  // Cancel pressed.
  const cancel = () => {
    console.log("cancelling"); 
    setHasChosenColor(true); 
  }

  return ( 
    hasChosenColor ? <Redirect to='/'/> : 
    <div>
      <SketchPicker 
        onChange={handleColorChange}      
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
