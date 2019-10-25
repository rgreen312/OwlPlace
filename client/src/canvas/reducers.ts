import { Color } from './types';
import * as ActionTypes from './actionTypes';
import { combineReducers,  } from 'redux';
import { createReducer } from '../createReducer';
import { RegisterContext, UpdatePosition, SetZoom } from './actions';
import * as WebSocketActionTypes from '../websocket/actionTypes';
import { DEFAULT_ZOOM } from './constants'

export interface State {
  canvasContext: CanvasRenderingContext2D | null;
  initialImage: string | null;
  curPosition: { x: number, y: number} | null;
  zoomFactor: number;
}

const canvasContext = createReducer<State['canvasContext']>(null, {
  [WebSocketActionTypes.ConnectError]: () => null,
  [ActionTypes.RegisterContext]: (state, action: RegisterContext) => action.payload.context
});

const curPosition = createReducer<State['curPosition']>(null, {
  [ActionTypes.ClearPosition]: () => null,
  [ActionTypes.UpdatePosition]: (state, action: UpdatePosition) => action.payload
});

const zoomFactor = createReducer<State['zoomFactor']>(DEFAULT_ZOOM, {
  [ActionTypes.SetZoom]: (state, action: SetZoom) => action.payload.zoom
})

export default combineReducers({
  canvasContext,
  curPosition,
  zoomFactor,
});
