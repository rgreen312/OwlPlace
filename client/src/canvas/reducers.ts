import { Color } from './types';
import * as ActionTypes from './actionTypes';
import { combineReducers,  } from 'redux';
import { createReducer } from '../createReducer';
import { RegisterContext, UpdatePosition } from './actions';
import * as WebSocketActionTypes from '../websocket/actionTypes';

export interface State {
  canvasContext: CanvasRenderingContext2D | null;
  initialImage: string | null;
  curPosition: { x: number, y: number} | null;
}

const canvasContext = createReducer<State['canvasContext']>(null, {
  [WebSocketActionTypes.ConnectError]: () => null,
  [ActionTypes.RegisterContext]: (state, action: RegisterContext) => action.payload.context
});

const curPosition = createReducer<State['curPosition']>(null, {
  [ActionTypes.ClearPosition]: () => null,
  [ActionTypes.UpdatePosition]: (state, action: UpdatePosition) => action.payload
});

export default combineReducers({
  canvasContext,
  curPosition,
});
