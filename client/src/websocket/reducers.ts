import { createReducer } from '../createReducer';
import * as ActionTypes from './actionTypes';
import { ConnectSuccess } from './actions';
import { combineReducers } from 'redux';
 
export interface State {
  socket: WebSocket | null;
  isConnected: boolean;
}

const socket = createReducer<State['socket']>(null, {
  [ActionTypes.StartConnect]: () => null,
  [ActionTypes.ConnectError]: () => null,
  [ActionTypes.ConnectSuccess]: (state, action: ConnectSuccess) => action.payload.socket,
});

const isConnected = createReducer<State['isConnected']>(false, {
  [ActionTypes.StartConnect]: () => false,
  [ActionTypes.ConnectError]: () => false,
  [ActionTypes.ConnectSuccess]: () => true,
});

export default combineReducers({
  socket,
  isConnected
})
