// import { FETCH_IMAGE_ENDPOINT } from './constants';
// import { HOSTNAME } from '../constants';
import { Color } from './types';
import * as ActionTypes from './actionTypes';

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

export const updateCursorPosition = (x: number, y: number) => dispatch => {
  // TODO: verify that numbers are actually within canvas
  dispatch(updatePosition(x, y));
}

const clearPosition = () => ({
  type: ActionTypes.ClearPosition
})
export type ClearPosition = ReturnType<typeof clearPosition>;

export const clearCursorPosition = () => dispatch => {
  dispatch(clearPosition());
}


export const updatePixel = (newColor: Color, x: number, y: number) => (dispatch) => {

}
