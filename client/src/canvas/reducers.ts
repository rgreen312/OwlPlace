import { Move } from './types';
import * as ActionTypes from './actionTypes';
import { combineReducers } from 'redux';
import { createReducer } from '../createReducer';
import { RegisterContext, UpdatePosition, SetZoom, SetInitialImage, SetTimeRemaining, UpdatePixelStart} from './actions';
import * as WebSocketActionTypes from '../websocket/actionTypes';
import { DEFAULT_ZOOM, TIME_BETWEEN_UPDATES_MS } from './constants';

export interface State {
  canvasContext: CanvasRenderingContext2D | null;
  initialImage: string | null;
  curPosition: { x: number; y: number } | null;
  zoomFactor: number;
  timeToNextChange: number;
  lastMove: Move | null;
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

const initialImage = createReducer<State['initialImage']>(null, {
  [ActionTypes.FetchImageError]: () => null,
  [ActionTypes.FetchImageStart]: () => null,
  [ActionTypes.FetchImageSuccess]: (state, action: SetInitialImage) => action.payload.image,
})
const timeToNextChange = createReducer<State['timeToNextChange']>(0, {
  [ActionTypes.SetTimeRemaining]: (state, action: SetTimeRemaining) => action.payload.time,
  [ActionTypes.UpdatePixelSuccess]: () => TIME_BETWEEN_UPDATES_MS,
});

const lastMove = createReducer<State['lastMove']>(null, {
  [ActionTypes.UpdatePixelStart]: (state, action: UpdatePixelStart) => ({ color: action.payload.color, position: action.payload.position }),
});

export default combineReducers({
  canvasContext,
  curPosition,
  zoomFactor,
  initialImage,
  timeToNextChange,
  lastMove
});
