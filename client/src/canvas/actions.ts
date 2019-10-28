// import { FETCH_IMAGE_ENDPOINT } from './constants';
// import { HOSTNAME } from '../constants';
import { Color } from './types';
import * as ActionTypes from './actionTypes';
import { getZoomFactor, getCanvasContext } from './selectors'; 
import { isNumber } from 'util';

const fetchImageDataStart = () => ({
  type: ActionTypes.FetchImageStart
});

export type FetchImageDataStart = ReturnType<typeof fetchImageDataStart>;

export const fetchImageData = () => (dispatch) => {
  dispatch(fetchImageDataStart());
  
}

const registerContext = (ctx: CanvasRenderingContext2D) => ({
  type: ActionTypes.RegisterContext,
  payload: {
    context: ctx
  }
});
export type RegisterContext = ReturnType<typeof registerContext>;

export const registerCanvasContext = (ctx: CanvasRenderingContext2D) => dispatch => {
  dispatch(registerContext(ctx));
}

const updatePosition = (x: number, y: number) => ({
  type: ActionTypes.UpdatePosition,
  payload: {
    x,
    y
  }
})
export type UpdatePosition = ReturnType<typeof updatePosition>;

export const updateCursorPosition = (x: number, y: number) => (dispatch, getState) => {
  const state = getState();
  const zoom = getZoomFactor(state);
  const testing = Math.ceil(x / zoom); 
  console.log("Here's x: " + isNumber(x));
  console.log("Here's zoom: " + zoom); 
  console.log("Here's testing: " + testing); 
  dispatch(updatePosition(Math.ceil(x / zoom), Math.ceil(y / zoom)));
}

const clearPosition = () => ({
  type: ActionTypes.ClearPosition
})
export type ClearPosition = ReturnType<typeof clearPosition>;

export const clearCursorPosition = () => dispatch => {
  dispatch(clearPosition());
}

const setZoom = (f: number) => ({
  type: ActionTypes.SetZoom,
  payload: {
    zoom: f
  }
});
export type SetZoom = ReturnType<typeof setZoom>;

export const setZoomFactor = (newFactor: number) => (dispatch, getState) => {
  
  if (newFactor < 0) {
    return;
  }

  const state = getState();
  const ctx = getCanvasContext(state);
  if (ctx) {
    ctx.scale(newFactor, newFactor);
  }
  
  dispatch(setZoom(newFactor));
}


export const updatePixel = (newColor: Color, x: number, y: number) => (dispatch) => {

}
