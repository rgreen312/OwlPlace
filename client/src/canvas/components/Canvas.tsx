import React, { Component, createRef, RefObject, useState } from 'react';
import ColorPicker from '../../colorPicker/components/colorPicker';
import { Redirect } from 'react-router-dom';
import './Canvas.scss';
import { Icon } from 'antd';
import { ZOOM_CHANGE_FACTOR } from '../constants';
import { Color, RGBColor } from 'react-color';

interface Props {
  receivedError: boolean;
  registerContext: (context: CanvasRenderingContext2D) => void;
  updatePosition: (x: number, y: number) => void;
  position: {x: number, y: number}; 
  onMouseOut: () => void;
  onUpdatePixel: (newColor: Color, x: number, y: number) => void;
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
    this.onCancel = this.onCancel.bind(this); 
    this.onComplete = this.onComplete.bind(this); 
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

    this.canvasRef.current!.addEventListener('mousemove', (ev) => {
      if (this.state.showColorPicker) return;
      const { x, y } = this.getMousePos(this.canvasRef.current, ev); 
      this.props.updatePosition(x, y);
    })

    this.canvasRef.current!.addEventListener('mouseout', () => {
      if (this.state.showColorPicker) return;
      this.props.onMouseOut();
    })

    this.canvasRef.current!.addEventListener('click', (ev) => {
      const { x, y } = this.getMousePos(this.canvasRef.current, ev);
      this.props.updatePosition(x, y);
      this.showColorPicker(); 
    }, false);
  }

  getMousePos(canvas, evt) {
    var rect = canvas.getBoundingClientRect();
    return {
      x: evt.clientX - rect.left,
      y: evt.clientY - rect.top
    };
  }

  onCancel() {
    this.hideColorPicker(); 
  }

  onComplete(c: RGBColor) {
    this.hideColorPicker(); 

    const context = this.canvasRef.current!.getContext('2d');

    const x = this.props.position.x; 
    const y = this.props.position.y;

    context!.fillStyle = 'rgb(' + c.r + ',' + c.g + ',' + c.b + ')'
    context!.fillRect(x, y, 1, 1);

    this.props.onUpdatePixel({ r: c.r, g: c.g, b: c.b}, x, y);
  }

  showColorPicker() {
    this.setState({
      showColorPicker: true
    });
  }

  hideColorPicker() {
    this.setState({
      showColorPicker: false
    })
  }

  render() {
    const { receivedError, zoomFactor, setZoomFactor } = this.props;
    return (
      // receivedError ? <Redirect to='/error'/> :
      <div className='canvas-container'>
        {this.state.showColorPicker && (<div className='color-picker'>
          <ColorPicker
            onCancel={this.onCancel}
            onComplete={(c) => this.onComplete(c)}
          />
        </div>)}
        <div className='zoom-canvas' style={{ transform: `scale(${zoomFactor}, ${zoomFactor})` }}>
          <canvas ref={this.canvasRef} />
        </div>
        <div className='zoom-controls'>
          <Icon
            type='plus-circle'
            onClick={() => setZoomFactor(zoomFactor + ZOOM_CHANGE_FACTOR)}
            className='zoom-icon'
          />
          <Icon
            type='minus-circle'
            onClick={() => setZoomFactor(zoomFactor - ZOOM_CHANGE_FACTOR)}
            className='zoom-icon'
          />
        </div>
      </div>
    );
  }
}

export default Canvas;
