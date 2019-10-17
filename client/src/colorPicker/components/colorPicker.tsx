import React, { FC } from "react";
import { SketchPicker } from "react-color";
import { Button } from "antd";
import { Redirect } from "react-router-dom";

interface Props {
  colorPicked: (Color?) => void; 
}
// http://casesandberg.github.io/react-color/

const ColorPicker: FC<Props> = ({colorPicked}) => {
  let chosenColor = null; 

  const handleChangeComplete = (color) => {
    chosenColor = color;
    console.log(chosenColor); 
  };

  const colorChosen = () => {
    console.log(chosenColor); 
    console.log("okay Pressed")
    colorPicked(chosenColor)
  }

  const cancel = () => {
    console.log("cancelling"); 
    return <Redirect to='/'/>; 
  }

  return (
    <div>
      <SketchPicker 
        onChange={handleChangeComplete}      
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
