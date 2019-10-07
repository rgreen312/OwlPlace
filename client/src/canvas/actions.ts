// import { FETCH_IMAGE_ENDPOINT } from './constants';
// import { HOSTNAME } from '../constants';
import { Color } from './types';
import * as ActionTypes from './actionTypes';

const fetchImageDataStart = () => ({
  type: ActionTypes.FetchImageStart
});

export type FetchImageDataStart = ReturnType<typeof fetchImageDataStart>;

//@ts-ignore
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

export const testAction = () => {
  console.log('In test action');
  //@ts-ignore
  return (dispatch) => {
    console.log('dispatching action');
    dispatch(testActionType('testing'));
  }
}
