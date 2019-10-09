import { createReducer } from '../createReducer';
import * as ActionTypes from './actionTypes';
import { combineReducers } from 'redux';
import { LoginSuccess } from './actions';


export interface State {
  isLoggedIn: boolean;
  name: string | null;
  email: string | null;
  userId: string | null;
}

const isLoggedIn = createReducer<State['isLoggedIn']>(false, {
  [ActionTypes.LoginStart]: () =>  false,
  [ActionTypes.LoginError]: () =>  false,
  [ActionTypes.SignOut]: () => false,
  [ActionTypes.LoginSuccess]: () => true,
});

const name = createReducer<State['name']>(null, {
  [ActionTypes.LoginStart]: () =>  null,
  [ActionTypes.LoginError]: () =>  null,
  [ActionTypes.SignOut]: () => null,
  [ActionTypes.LoginSuccess]: (state, action: LoginSuccess) => action.payload.name,
});

const email = createReducer<State['name']>(null, {
  [ActionTypes.LoginStart]: () =>  null,
  [ActionTypes.LoginError]: () =>  null,
  [ActionTypes.SignOut]: () => null,
  [ActionTypes.LoginSuccess]: (state, action: LoginSuccess) => action.payload.email,
});

const id = createReducer<State['userId']>(null, {
  [ActionTypes.LoginStart]: () =>  null,
  [ActionTypes.LoginError]: () =>  null,
  [ActionTypes.SignOut]: () => null,
  [ActionTypes.LoginSuccess]: (state, action: LoginSuccess) => action.payload.id,
});



export default combineReducers({
  isLoggedIn,
  name,
  email,
  id,
});
