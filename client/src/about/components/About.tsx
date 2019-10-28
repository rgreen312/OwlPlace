import React, { FC } from 'react';
import { sendUpdateMessage } from '../../websocket/actions'


interface Props {
  sendUpdateMessage: (id, x, y, r, g, b) => void;
}

const About: FC<Props> = ({ sendUpdateMessage }) => (
  <div className='about-page'>
    Update the canvas one pixel at a time...

    <h2>change the canvas one pixel at a time</h2>
      <p>Click "Pixel 1" to send an update message to the server!
      </p>
        <button onClick = {() => sendUpdateMessage("user1", 10, 400, 255, 255, 255)} id="p1"> Pixel 1 </button>
        {/* <button onClick= {onClickClose} id="close">Close</button> */}
  </div>
);

export default About;
