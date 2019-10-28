import React, { Component, createRef, RefObject, useState } from 'react';
import ColorPicker from '../../colorPicker/components/colorPicker';
import { Redirect } from 'react-router-dom';
import './Canvas.scss';
import { Icon } from 'antd';
import { ZOOM_CHANGE_FACTOR } from '../constants';
import Modal from 'react-modal'; 

interface Position {
  x: number;
  y: number;
}

interface Props {
  receivedError: boolean;
  registerContext: (context: CanvasRenderingContext2D) => void;
  updatePosition: (Position) => void;
  onMouseOut: () => void;
  zoomFactor: number;
  setZoomFactor: (newZoom: number) => void;
}

interface State {
  showColorPicker: boolean; 
}

class Canvas extends Component<Props, State> {
  canvasRef: RefObject<HTMLCanvasElement>;

  constructor(props) {
    super(props);
    this.canvasRef = createRef();
    this.state = {showColorPicker: false}
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

    // TODO: remove this code
    context!.fillStyle ='#000000';
    context!.fillRect(0, 0, 1000, 500);
    context!.fillStyle = '#ff0000';
    context!.fillRect(0, 500, 1000, 500);

    context!.scale(100, 100);

    this.canvasRef.current!.addEventListener('mousemove', (ev) => {
      const { x, y } = this.getMousePos(this.canvasRef.current, ev);
      this.props.updatePosition({x, y});
    }, false);

    this.canvasRef.current!.addEventListener('mouseout', () => {
      this.props.onMouseOut();
    })

    this.canvasRef.current!.addEventListener('click', (ev) => {
      console.log("setting to true");
      this.showColorPicker(); 
    });
  }

  getMousePos(canvas, evt) {
    var rect = canvas.getBoundingClientRect();
    return {
      x: evt.clientX - rect.left,
      y: evt.clientY - rect.top
    };
  }

  showColorPicker() {
    this.setState({
      showColorPicker: true
    });
  }

  render() {
    const { receivedError, zoomFactor, setZoomFactor } = this.props;
    return (
      // receivedError ? <Redirect to='/error'/> :
      <div className='canvas-container'>
        {this.state.showColorPicker ? (<div className='color-picker'>
          <ColorPicker
            onCancel={() => console.log('canceled')}
            onComplete={(c) => console.log('color selected: ', c)}
          />
        </div>) : null}
        <div className='zoom-canvas' style={{ transform: `scale(${zoomFactor}, ${zoomFactor})` }}>
          <canvas ref={this.canvasRef} />
        </div>
        <div className='zoom-controls'>
          <Icon
            type='plus-circle'
            onClick={() => setZoomFactor(zoomFactor + ZOOM_CHANGE_FACTOR)}
          />
          <Icon
            type='minus-circle'
            onClick={() => setZoomFactor(zoomFactor - ZOOM_CHANGE_FACTOR)}
          />
        </div>
      </div>
    );
  }
}

export default Canvas;
