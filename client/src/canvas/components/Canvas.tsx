import React, { FC } from 'react';

interface Props {
  onClick: () => void;
}

const Canvas: FC<Props> = ({ onClick }) => {
  return (
    <div>
      <div>
        <canvas></canvas>
      </div>
    </div>
  );
}

export default Canvas;
