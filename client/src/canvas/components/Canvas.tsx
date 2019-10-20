import React, { Component, createRef, RefObject } from 'react';
import ColorPicker from '../../colorPicker/components/colorPicker';
import { Redirect } from 'react-router-dom';

interface Props {
  receivedError: boolean;
  registerContext: (context: CanvasRenderingContext2D) => void;
}

class Canvas extends Component<Props> {
  canvasRef: RefObject<HTMLCanvasElement>;

  constructor(props) {
    super(props);
    this.canvasRef = createRef();
  }

  componentDidMount() {
    this.canvasRef.current!.width = 1000;
    this.canvasRef.current!.height = 1000;

    const context = this.canvasRef.current!.getContext('2d');
    const image = new Image();

    image.onload = function() {
      if (context) {
        context.drawImage(image, 0, 0);
      }
    };
    // image.src = this.props.initialImage;

    if (context) {
      this.props.registerContext(context);
    }
  }

  render() {
    const { receivedError } = this.props;
    return (
      // receivedError ? <Redirect to='/error'/> :
      <div>
        <div>
          {/* <ColorPicker
            onCancel={() => console.log('canceled')}
            onComplete={(c) => console.log('color selected: ', c)}
          /> */}
          <canvas ref={this.canvasRef} />
        </div>
      </div>
    );
  }
}

export default Canvas;
