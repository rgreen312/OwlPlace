import { Color } from './types';
import * as ActionTypes from './actionTypes';
import { combineReducers,  } from 'redux';
import { createReducer } from '../createReducer';
import { RegisterContext } from './actions';
import * as WebSocketActionTypes from '../websocket/actionTypes';

export interface State {
  canvasContext: CanvasRenderingContext2D | null;
  initialImage: string | null;
}

const canvasContext = createReducer<State['canvasContext']>(null, {
  [WebSocketActionTypes.ConnectError]: () => null,
  [ActionTypes.RegisterContext]: (state, action: RegisterContext) => action.payload.context
})

export default combineReducers({
  canvasContext,
});
