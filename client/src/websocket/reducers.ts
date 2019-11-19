import { createReducer } from '../createReducer';
import * as ActionTypes from './actionTypes';
import { ConnectSuccess } from './actions';
import { combineReducers } from 'redux';

export interface State {
  socket: WebSocket | null;
  isConnected: boolean;
  receivedError: boolean;
  isLoading: boolean;
}

const socket = createReducer<State['socket']>(null, {
  [ActionTypes.StartConnect]: () => null,
  [ActionTypes.ConnectError]: () => null,
  [ActionTypes.CloseConnection]: () => null,
  [ActionTypes.ConnectSuccess]: (state, action: ConnectSuccess) =>
    action.payload.socket
});

const isConnected = createReducer<State['isConnected']>(false, {
  [ActionTypes.StartConnect]: () => false,
  [ActionTypes.ConnectError]: () => false,
  [ActionTypes.CloseConnection]: () => false,
  [ActionTypes.ConnectSuccess]: () => true
});

const receivedError = createReducer<State['receivedError']>(false, {
  [ActionTypes.StartConnect]: () => false,
  [ActionTypes.ConnectError]: () => true,
  [ActionTypes.CloseConnection]: () => false,
  [ActionTypes.ConnectSuccess]: () => false
});

const isLoading = createReducer<State['isLoading']>(false, {
  [ActionTypes.StartConnect]: () => true,
  [ActionTypes.ConnectError]: () => false,
  [ActionTypes.CloseConnection]: () => false,
  [ActionTypes.ConnectSuccess]: () => false
});

export default combineReducers({
  socket,
  isConnected,
  receivedError,
  isLoading
});
