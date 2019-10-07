import { Color } from './types';
import * as ActionTypes from './actionTypes';
import { combineReducers,  } from 'redux';
import { createReducer } from '../createReducer';

export interface State {
  image: Color[][];
  location: {
    x: number;
    y: number;
  };
}



export default combineReducers({
  
});
