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
  ctx.fillStyle = '#000000';
  ctx.fillRect(0, 0, 100, 100);
}


export const updatePixel = (newColor: Color, x: number, y: number) => (dispatch) => {

}
