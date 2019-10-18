import React, { FC } from 'react';
import ColorPicker from '../../colorPicker/components/colorPicker';

interface Props {

}

const Canvas: FC<Props> = () => {
  return (
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
