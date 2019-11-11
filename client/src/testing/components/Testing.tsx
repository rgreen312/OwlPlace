import React, { FC } from 'react';
import { sendUpdateMessage, sendLoginMessage } from '../../websocket/actions';

interface Props {
  sendUpdateMessage: (id, x, y, r, g, b) => void;
  sendLoginMessage: (id) => void;
}


const Testing: FC<Props> = ({ sendUpdateMessage, sendLoginMessage }) => (
  <div className='testing-page'>
    <h2>TESTING AREA</h2>
    <p>Click "Pixel 1" to send an update message to the server!</p>
    <button onClick = {() => sendUpdateMessage("user1", 10, 400, 255, 255, 255)} id="p1"> Pixel 1 </button>
    <button onClick = {() => sendLoginMessage("testemail@gmail.com")} id="login"> User Login </button>
  </div>
);

export default Testing;
