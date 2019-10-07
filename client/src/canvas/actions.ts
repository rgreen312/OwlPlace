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


export const updatePixel = (newColor: Color, x: number, y: number) => {

}


const testActionType = (data: string) => ({
  type: ActionTypes.FetchImageStart,
  payload:{ data }
});
export type TestActionType = ReturnType<typeof testActionType>;

export const testAction = () => (dispatch, getState) => {
  console.log('dispatching action');
  dispatch(testActionType('testing'));
}
