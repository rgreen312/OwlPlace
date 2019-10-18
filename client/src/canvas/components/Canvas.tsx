import React, { FC } from 'react';
import ColorPicker from '../../colorPicker/components/colorPicker';
import { Redirect } from 'react-router-dom';

interface Props {
  receivedError: boolean; 
}

const Canvas: FC<Props> = ({receivedError}) => {
  return (
    receivedError ? <Redirect to='/error'/> :
    <div>
      <div>
        <ColorPicker
          onCancel={() => console.log('canceled')}
          onComplete={(c) => console.log('color selected: ', c)}
        />
      </div>
    </div>
  );
}

export default Canvas;
