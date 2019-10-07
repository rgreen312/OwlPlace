import { Color } from './types';
import * as ActionTypes from './actionTypes';
import { combineReducers,  } from 'redux';
import { createReducer } from '../createReducer';
import { TestActionType } from './actions'

interface Update {
  x: number;
  y: number;
  color: Color;
}

interface State {
  image: Color[][];
  location: {
    x: number;
    y: number;
  };
  updates: Update[];
  test: string;
}

const imageReducer = createReducer([][], {
  [ActionTypes.FetchImageError]: () =>  [][],
  [ActionTypes.FetchImageSuccess]: (state: State, action) => [][],
});

const testReducer = createReducer(null, {
  [ActionTypes.FetchImageError]: () =>  null,
  [ActionTypes.FetchImageSuccess]: (state: State, action: TestActionType) => action.payload.data,
});

export default combineReducers({
  testReducer,
});
