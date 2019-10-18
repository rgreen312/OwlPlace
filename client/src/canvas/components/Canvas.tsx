import React, { FC } from 'react';
import { Redirect } from 'react-router-dom'; 

interface Props {
  receivedError: boolean; 
}

const Canvas: FC<Props> = ({receivedError}) => {
  return (
    receivedError ? <Redirect to='/error'/> :
    <div>
      <div>
        <canvas></canvas>
      </div>
    </div>
  );
}

export default Canvas;
