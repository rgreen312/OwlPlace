import React, { Component, createRef, RefObject, useState } from 'react';
import ColorPicker from '../../colorPicker/components/colorPicker';
import { Redirect } from 'react-router-dom';
import './Canvas.scss';
import { Icon, Spin } from 'antd';
import { ZOOM_CHANGE_FACTOR } from '../constants';
import { Color, RGBColor } from 'react-color';
import classNames from 'classnames';

interface Props {
  receivedError: boolean;
  isLoading: boolean;
  registerContext: (context: CanvasRenderingContext2D) => void;
  updatePosition: (x: number, y: number) => void;
  position: { x: number; y: number };
  onMouseOut: () => void;
  onUpdatePixel: (newColor: Color, x: number, y: number) => void;
  zoomFactor: number;
  setZoomFactor: (newZoom: number) => void;
  initialImage?: string;
}

interface State {
  showColorPicker: boolean;
  previousColor: RGBColor | null;
  isDrag: boolean;
  translateX: number;
  translateY: number;
  dragStartX: number;
  dragStartY: number;
  colorPickerX: number;
  colorPickerY: number;
}

class Canvas extends Component<Props, State> {
  canvasRef: RefObject<HTMLCanvasElement>;

  constructor(props) {
    super(props);
    this.canvasRef = createRef();
    this.state = {
      showColorPicker: false,
      previousColor: null,
      isDrag: false,
      translateX: 0,
      translateY: 0,
      dragStartX: 0,
      dragStartY: 0,
      colorPickerX: 0,
      colorPickerY: 0
    };
    this.onCancel = this.onCancel.bind(this);
    this.onComplete = this.onComplete.bind(this);
    this.onColorChange = this.onColorChange.bind(this);
    this.updateTranslate = this.updateTranslate.bind(this);
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
    context!.fillStyle = '#000000';
    context!.fillRect(0, 0, 1000, 500);
    context!.fillStyle = '#ff0000';
    context!.fillRect(0, 500, 1000, 500);

    this.canvasRef.current!.addEventListener('mousemove', ev => {
      if (this.state.showColorPicker) return;
      const { x, y } = this.getMousePos(this.canvasRef.current, ev);
      this.props.updatePosition(x, y);
    });

    this.canvasRef.current!.addEventListener('mouseout', () => {
      if (this.state.showColorPicker) return;
      this.props.onMouseOut();
    });

    // On mousedown, get the current location to be used for dragging
    this.canvasRef.current!.addEventListener('mousedown', e => {
      const { zoomFactor } = this.props;
      const { translateX, translateY } = this.state;
      const startPositionX = e.clientX - translateX;
      const startPositionY = e.clientY - translateY;

      let pickerX = e.clientX + zoomFactor;
      let pickerY = e.clientY;

      // The color picker is 220 x 337 pixels
      if (pickerX + 220 > window.innerWidth) {
        pickerX -= 220;
      }
      if (pickerY + 337 > window.innerHeight - 70) {
        pickerY -= 337;
      }

      this.setState({
        dragStartX: startPositionX,
        dragStartY: startPositionY,
        colorPickerX: pickerX,
        colorPickerY: pickerY
      });

      // If the user moves after clicking, then they are dragging so we add listener
      this.canvasRef.current!.addEventListener(
        'mousemove',
        this.updateTranslate
      );
    });

    /**
     * On mouse up, we check if the user was dragging. If they were, we set isDrag to false and
     * remove the event listener for dragging.
     *
     * If they were not dragging, then we display the color picker so we can update the color of
     * the pixel.
     */
    this.canvasRef.current!.addEventListener('mouseup', ev => {
      this.canvasRef.current!.removeEventListener(
        'mousemove',
        this.updateTranslate
      );

      if (!this.state.isDrag) {
        const { x, y } = this.getMousePos(this.canvasRef.current, ev);
        this.props.updatePosition(x, y);
        this.showColorPicker();
      }

      const { x, y } = this.getMousePos(this.canvasRef.current, ev);
      const imageData = this.canvasRef
        .current!.getContext('2d')!
        .getImageData(x, y, 1, 1);
      this.setState({
        previousColor: {
          r: imageData.data[1],
          g: imageData.data[2],
          b: imageData.data[3]
        }
      });
      this.showColorPicker();

      this.setState({
        isDrag: false
      });
    });
  }

  componentDidUpdate(prevProps: Props) {
    console.log('component updated');
    if (this.props.initialImage && this.props.initialImage !== prevProps.initialImage) {
      console.log('image drawn');
      const context = this.canvasRef.current!.getContext('2d');
      if (context) {
        const image = new Image();
        image.onload = function() {
          context.drawImage(image, 0, 0);
        };
        image.src = this.props.initialImage;
      }
    }
  }

  updateTranslate(ev: MouseEvent) {
    const { dragStartX, dragStartY } = this.state;
    const x = ev.clientX - dragStartX;
    const y = ev.clientY - dragStartY;
    this.onCancel();
    this.setState({
      isDrag: true,
      translateX: x,
      translateY: y
    });
  }

  getMousePos(canvas, evt) {
    var rect = canvas.getBoundingClientRect();
    return {
      x: evt.clientX - rect.left,
      y: evt.clientY - rect.top
    };
  }

  onCancel() {
    this.hideColorPicker(true);
  }

  onComplete() {
    this.hideColorPicker(false);
  }

  onColorChange(c: RGBColor) {
    const context = this.canvasRef.current!.getContext('2d');
    const x = this.props.position.x - 1;
    const y = this.props.position.y - 1;
    context!.fillStyle = `rgb(${c.r}, ${c.g}, ${c.b})`;
    context!.fillRect(x, y, 1, 1);
    this.props.onUpdatePixel({ r: c.r, g: c.g, b: c.b }, x, y);
  }

  showColorPicker() {
    this.setState({
      showColorPicker: true
    });
  }

  hideColorPicker(didCancel: boolean) {
    this.setState({
      showColorPicker: false
    });

    // We should change the color back if cancel was pressed.
    if (didCancel) {
      const context = this.canvasRef.current!.getContext('2d');

      const x = this.props.position.x;
      const y = this.props.position.y;
      const c = this.state.previousColor!;

      context!.fillStyle = 'rgb(' + c.r + ',' + c.g + ',' + c.b + ')';
      context!.fillRect(x - 1, y - 1, 1, 1);

      this.props.onUpdatePixel({ r: c.r, g: c.g, b: c.b }, x, y);
    }
  }

  render() {
    const { receivedError, zoomFactor, setZoomFactor, isLoading } = this.props;
    const {
      translateX,
      translateY,
      isDrag,
      colorPickerX,
      colorPickerY
    } = this.state;
    return (
      // receivedError ? <Redirect to='/error'/> :
      <div>
        {isLoading && <Spin className='spinner-style' size='large' />}
        <div
          className={classNames({
            'hide-canvas': isLoading,
            'canvas-container': true
          })}
        >
          {this.state.showColorPicker && (
            <ColorPicker
              onColorChange={c => this.onColorChange(c)}
              onCancel={this.onCancel}
              onComplete={this.onComplete}
              className='color-picker'
              style={{ top: `${colorPickerY}px`, left: `${colorPickerX}px` }}
            />
          )}

          <div
            className={classNames({
              'pan-canvas': true,
              'drag-canvas': isDrag
            })}
            style={{ transform: `translate(${translateX}px, ${translateY}px)` }}
          >
            <div
              className='zoom-canvas'
              style={{ transform: `scale(${zoomFactor}, ${zoomFactor})` }}
            >
              <canvas ref={this.canvasRef} />
            </div>
          </div>
          <div className='zoom-controls'>
            <Icon
              type='plus-circle'
              onClick={() => {
                setZoomFactor(zoomFactor + ZOOM_CHANGE_FACTOR);
                this.onCancel();
              }}
              className='zoom-icon'
            />
            <Icon
              type='minus-circle'
              onClick={() => {
                setZoomFactor(zoomFactor - ZOOM_CHANGE_FACTOR);
                this.onCancel();
              }}
              className='zoom-icon'
            />
          </div>
        </div>
      </div>
    );
  }
}

export default Canvas;
