import { createReducer } from '../createReducer';
import * as ActionTypes from './actionTypes';
import { combineReducers } from 'redux';


interface State {
  isLoggedIn: boolean;
  name: string | null;
  email: string | null;
  id: string | null;
}

const isLoggedIn = createReducer<State['isLoggedIn']>(false, {
  [ActionTypes.LoginStart]: () =>  false,
  [ActionTypes.LoginStart]: () =>  false,
  [ActionTypes.LoginSuccess]: () => true,
});

export default combineReducers({
  isLoggedIn,
});
