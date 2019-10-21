import React, { Component, createRef, RefObject } from 'react';
import ColorPicker from '../../colorPicker/components/colorPicker';
import { Redirect } from 'react-router-dom';
import './Canvas.scss';
import { Icon } from 'antd';

interface Props {
  receivedError: boolean;
  registerContext: (context: CanvasRenderingContext2D) => void;
  updatePosition: (x: number, y: number) => void;
  onMouseOut: () => void;
  zoomFactor: number;
  setZoomFactor: (newZoom: number) => void;
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
    // const image = new Image();

    // image.onload = function() {
    //   if (context) {
    //     context.drawImage(image, 0, 0);
    //   }
    // };
    // image.src = this.props.initialImage;


    if (context) {
      this.props.registerContext(context);
    }

    context!.imageSmoothingEnabled = false;

    context!.fillStyle ='#000000';
    context!.fillRect(0, 0, 1000, 500);
    context!.fillStyle = '#ff0000';
    context!.fillRect(0, 500, 1000, 500);

    context!.scale(100, 100);

    this.canvasRef.current!.addEventListener('mousemove', (ev) => {
      const { x, y } = this.getMousePos(this.canvasRef.current, ev);
      this.props.updatePosition(x, y);
    }, false);

    this.canvasRef.current!.addEventListener('mouseout', () => {
      this.props.onMouseOut();
    })
  }

  getMousePos(canvas, evt) {
    var rect = canvas.getBoundingClientRect();
    return {
      x: evt.clientX - rect.left,
      y: evt.clientY - rect.top
    };
  }

  render() {
    const { receivedError, zoomFactor, setZoomFactor } = this.props;
    return (
      // receivedError ? <Redirect to='/error'/> :
      <div className='canvas-container'>
        {/* <ColorPicker
          onCancel={() => console.log('canceled')}
          onComplete={(c) => console.log('color selected: ', c)}
        /> */}
        <div className='zoom-canvas' style={{ transform: `scale(${zoomFactor}, ${zoomFactor})` }}>
          <canvas ref={this.canvasRef} />
        </div>
        <div className='zoom-controls'>
          <Icon
            type='plus-circle'
            onClick={() => setZoomFactor(zoomFactor + 10)}
          />
          <Icon
            type='minus-circle'
            onClick={() => setZoomFactor(zoomFactor - 10)}
          />
        </div>
      </div>
    );

    // return (
    //     <div className='zoom-canvas' style={{ transform: `scale(${zoomFactor}, ${zoomFactor})` }}>
    //       <canvas ref={this.canvasRef} />
    //     </div>
    // );
  }
}

export default Canvas;
