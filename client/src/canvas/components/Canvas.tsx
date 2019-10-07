import React, { FC } from 'react';

interface Props {
  onClick: () => void;
}

const Canvas: FC<Props> = ({ onClick }) => {
  return (
    <div>
      <div>
        This is a canvas
        <button onClick={onClick}>Test</button>
      </div>
    </div>
  );
}

export default Canvas;
