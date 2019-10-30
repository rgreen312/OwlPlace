import React, { Component, createRef, RefObject } from "react";
import ColorPicker from "../../colorPicker/components/colorPicker";
import { Redirect } from "react-router-dom";
import "./Canvas.scss";
import { Icon } from "antd";
import { ZOOM_CHANGE_FACTOR } from "../constants";

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
  mouseDown: Boolean;
  widthCanvas: any;
  heightCanvas: any;
  xleftView: any;
  ytopView: any;
  widthViewOriginal: any;
  heightViewOriginal: any;
  widthView: any;
  heightView: any;
  lastX: any;
  lastY: any;

  constructor(props) {
    super(props);
    this.canvasRef = createRef();
    this.mouseDown = false;
    this.widthCanvas = 0;
    this.heightCanvas = 0;
    this.xleftView = 0;
    this.ytopView = 0;
    this.widthViewOriginal = 1.0; //actual width and height of zoomed and panned display
    this.heightViewOriginal = 1.0;
    this.widthView = this.widthViewOriginal; //actual width and height of zoomed and panned display
    this.heightView = this.heightViewOriginal;
    this.lastX = 0;
    this.lastY = 0;
    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseWheel = this.handleMouseWheel.bind(this);
    this.handleDblClick = this.handleDblClick.bind(this);
    this.handleMouseUp = this.handleMouseUp.bind(this);
    this.handleMouseDown = this.handleMouseDown.bind(this);
  }

  componentDidMount() {
    this.setup();
  }

  setup() {
    this.canvasRef.current!.width = 1000;
    this.canvasRef.current!.height = 1000;
    this.widthCanvas = this.canvasRef.current!.width;
    this.heightCanvas = this.canvasRef.current!.height;
    const context = this.canvasRef.current!.getContext("2d");

    if (context) {
      this.props.registerContext(context);
    }

    this.canvasRef.current!.addEventListener(
      "dblclick",
      this.handleDblClick,
      false
    ); // dblclick to zoom in at point, shift dblclick to zoom out.
    this.canvasRef.current!.addEventListener(
      "mousedown",
      this.handleMouseDown,
      false
    ); // click and hold to pan
    this.canvasRef.current!.addEventListener(
      "mousemove",
      this.handleMouseMove,
      false
    );
    this.canvasRef.current!.addEventListener(
      "mouseup",
      this.handleMouseUp,
      false
    );
    this.canvasRef.current!.addEventListener(
      "mousewheel",
      this.handleMouseWheel,
      false
    );

    this.draw(context);
  }

  draw(ctx) {
    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.scale(
      this.widthCanvas / this.widthView,
      this.heightCanvas / this.heightView
    );
    ctx.translate(-this.xleftView, -this.ytopView);

    ctx.fillStyle = "yellow";
    ctx.fillRect(
      this.xleftView,
      this.ytopView,
      this.widthView,
      this.heightView
    );
    ctx.fillStyle = "blue";
    ctx.fillRect(0.1, 0.5, 0.1, 0.1);
    ctx.fillStyle = "red";
    ctx.fillRect(0.3, 0.2, 0.4, 0.2);
  }

  handleDblClick(event) {
    var X =
      event.clientX -
      this.canvasRef.current!.offsetLeft -
      this.canvasRef.current!.clientLeft +
      this.canvasRef.current!.scrollLeft; //Canvas coordinates
    var Y =
      event.clientY -
      this.canvasRef.current!.offsetTop -
      this.canvasRef.current!.clientTop +
      this.canvasRef.current!.scrollTop;
    var x = (X / this.widthCanvas) * this.widthView + this.xleftView; // View coordinates
    var y = (Y / this.heightCanvas) * this.heightView + this.ytopView;

    var scale = event.shiftKey == 1 ? 1.5 : 0.5; // shrink (1.5) if shift key pressed
    this.widthView *= scale;
    this.heightView *= scale;

    if (
      this.widthView > this.widthViewOriginal ||
      this.heightView > this.heightViewOriginal
    ) {
      this.widthView = this.widthViewOriginal;
      this.heightView = this.heightViewOriginal;
      x = this.widthView / 2;
      y = this.heightView / 2;
    }

    this.xleftView = x - this.widthView / 2;
    this.ytopView = y - this.heightView / 2;

    const context = this.canvasRef.current!.getContext("2d");

    this.draw(context);
  }

  handleMouseDown(event) {
    this.mouseDown = true;
  }

  handleMouseUp(event) {
    this.mouseDown = false;
  }

  handleMouseMove(event) {
    var X =
      event.clientX -
      this.canvasRef.current!.offsetLeft -
      this.canvasRef.current!.clientLeft +
      this.canvasRef.current!.scrollLeft;
    var Y =
      event.clientY -
      this.canvasRef.current!.offsetTop -
      this.canvasRef.current!.clientTop +
      this.canvasRef.current!.scrollTop;

    if (this.mouseDown) {
      var dx = ((X - this.lastX) / this.widthCanvas) * this.widthView;
      var dy = ((Y - this.lastY) / this.heightCanvas) * this.heightView;
      this.xleftView -= dx;
      this.ytopView -= dy;
    }
    this.lastX = X;
    this.lastY = Y;

    const context = this.canvasRef.current!.getContext("2d");
    this.draw(context);
  }

  handleMouseWheel(event) {
    var x = this.widthView / 2 + this.xleftView; // View coordinates
    var y = this.heightView / 2 + this.ytopView;

    var scale = event.wheelDelta < 0 || event.detail > 0 ? 1.1 : 0.9;
    this.widthView *= scale;
    this.heightView *= scale;

    if (
      this.widthView > this.widthViewOriginal ||
      this.heightView > this.heightViewOriginal
    ) {
      this.widthView = this.widthViewOriginal;
      this.heightView = this.heightViewOriginal;
      x = this.widthView / 2;
      y = this.heightView / 2;
    }

    // scale about center of view, rather than mouse position. This is different than dblclick behavior.
    this.xleftView = x - this.widthView / 2;
    this.ytopView = y - this.heightView / 2;
    const context = this.canvasRef.current!.getContext("2d");

    this.draw(context);
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
      <div className="canvas-container">
        {/* <ColorPicker
          onCancel={() => console.log('canceled')}
          onComplete={(c) => console.log('color selected: ', c)}
        /> */}
        <div
          className="zoom-canvas"
          //   style={{ transform: `scale(${zoomFactor}, ${zoomFactor})` }}
        >
          <canvas ref={this.canvasRef} />
        </div>
        <div className="zoom-controls">
          <Icon
            type="plus-circle"
            onClick={() => setZoomFactor(zoomFactor + ZOOM_CHANGE_FACTOR)}
          />
          <Icon
            type="minus-circle"
            onClick={() => setZoomFactor(zoomFactor - ZOOM_CHANGE_FACTOR)}
          />
        </div>
      </div>
    );
  }
}

export default Canvas;
