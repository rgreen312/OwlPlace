// import { FETCH_IMAGE_ENDPOINT } from './constants';
// import { HOSTNAME } from '../constants';
import { Color } from './types';
import * as ActionTypes from './actionTypes';
import { getZoomFactor } from './selectors';
import { getWebSocket } from '../websocket/selectors';
import { getUserEmail } from '../login/selectors';
import { connectError, makeUpdateMessage } from '../websocket/actions';
import { HOSTNAME } from '../constants';

const fetchImageDataStart = () => ({
  type: ActionTypes.FetchImageStart
});
export type FetchImageDataStart = ReturnType<typeof fetchImageDataStart>;

export const fetchImageData = () => dispatch => {
  dispatch(fetchImageDataStart());
};

export const setImage = (image: string) => ({
  type: ActionTypes.FetchImageSuccess,
  payload: {
    image
  }
});
export type SetInitialImage = ReturnType<typeof setImage>;

export const setInitialImage = (image: string) => dispatch => {
  console.log('setting image in state');
  dispatch(setImage(image));
}

const registerContext = (ctx: CanvasRenderingContext2D) => ({
  type: ActionTypes.RegisterContext,
  payload: {
    context: ctx
  }
});
export type RegisterContext = ReturnType<typeof registerContext>;

export const registerCanvasContext = (
  ctx: CanvasRenderingContext2D
) => dispatch => {
  dispatch(registerContext(ctx));
};

const updatePosition = (x: number, y: number) => ({
  type: ActionTypes.UpdatePosition,
  payload: {
    x,
    y
  }
});
export type UpdatePosition = ReturnType<typeof updatePosition>;

export const updateCursorPosition = (x: number, y: number) => (
  dispatch,
  getState
) => {
  const state = getState();
  const zoom = getZoomFactor(state);
  dispatch(updatePosition(Math.ceil(x / zoom), Math.ceil(y / zoom)));
};

const clearPosition = () => ({
  type: ActionTypes.ClearPosition
});
export type ClearPosition = ReturnType<typeof clearPosition>;

export const clearCursorPosition = () => dispatch => {
  dispatch(clearPosition());
};

const setZoom = (f: number) => ({
  type: ActionTypes.SetZoom,
  payload: {
    zoom: f
  }
});
export type SetZoom = ReturnType<typeof setZoom>;

export const setZoomFactor = (newFactor: number) => dispatch => {
  if (newFactor < 0) {
    return;
  }

  dispatch(setZoom(newFactor));
};

const updatePixelStart = (newColor: Color, x: number, y: number) => ({
  type: ActionTypes.UpdatePixelStart,
  payload: {
    color: newColor,
    position: { x, y }
  }
});
export type UpdatePixelStart = ReturnType<typeof updatePixelStart>;

export const updatePixel = (
  newColor: Color,
  x: number,
  y: number
) => (dispatch, getState) => {
  dispatch(updatePixelStart(newColor, x, y));
  const socket = getWebSocket(getState());
  const email = getUserEmail(getState());
  console.log("This is being called!", newColor.r, newColor.g, newColor.b, x, y)
  console.log(getState())
  if (socket && email) {
    console.log("Socket sending!")
    socket.send(makeUpdateMessage(email, x, y, newColor.r, newColor.g, newColor.b));
  } else {
    dispatch(connectError('Could not connect'))
  }
};

export const setTimeRemaining = (time: number) => ({
  type: ActionTypes.SetTimeRemaining,
  payload: {
    time,
  }
});
export type SetTimeRemaining = ReturnType<typeof setTimeRemaining>;

export const setTimeToNextMove = (time: number) => dispatch => {
  // Don't let time remaining be negative
  if (time < 0) {
    dispatch(setTimeRemaining(0));
  }
  dispatch(setTimeRemaining(time));
}

export const getImageData = () => async dispatch => {
  const response = await fetch(`https://cors-anywhere.herokuapp.com/${HOSTNAME}/json/image`);
  const json = await response.json();
  const data = json.formatString;

  if (data) {
    console.log('setting image to ', data);
    dispatch(setImage('data:image/png;base64,' + data));
  }
}
