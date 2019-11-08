import { Color } from './types';
import * as ActionTypes from './actionTypes';
import { combineReducers } from 'redux';
import { createReducer } from '../createReducer';
import { RegisterContext, UpdatePosition, SetZoom, SetTimeRemaining } from './actions';
import * as WebSocketActionTypes from '../websocket/actionTypes';
import { DEFAULT_ZOOM, TIME_BETWEEN_UPDATES_MS } from './constants';

export interface State {
  canvasContext: CanvasRenderingContext2D | null;
  initialImage: string | null;
  curPosition: { x: number; y: number } | null;
  zoomFactor: number;
  timeToNextChange: number;
}

const canvasContext = createReducer<State['canvasContext']>(null, {
  [WebSocketActionTypes.ConnectError]: () => null,
  [ActionTypes.RegisterContext]: (state, action: RegisterContext) =>
    action.payload.context
});

const curPosition = createReducer<State['curPosition']>(null, {
  [ActionTypes.ClearPosition]: () => null,
  [ActionTypes.UpdatePosition]: (state, action: UpdatePosition) =>
    action.payload
});

const zoomFactor = createReducer<State['zoomFactor']>(DEFAULT_ZOOM, {
  [ActionTypes.SetZoom]: (state, action: SetZoom) => action.payload.zoom
});

const timeToNextChange = createReducer<State['timeToNextChange']>(10000, { // TODO: default to 0
  [ActionTypes.SetTimeRemaining]: (state, action: SetTimeRemaining) => action.payload.time,
  [ActionTypes.UpdatePixelSuccess]: () => TIME_BETWEEN_UPDATES_MS,
});

export default combineReducers({
  canvasContext,
  curPosition,
  zoomFactor,
  timeToNextChange,
});
