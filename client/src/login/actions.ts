import * as ActionTypes from './actionTypes';
import { googleAPILoaded } from '../App';

const loginStart = () => ({
  type: ActionTypes.LoginStart
});
export type LoginStart = ReturnType<typeof loginStart>;

export const loginSuccess = (name: string, id: string, email: string) => ({
  type: ActionTypes.LoginSuccess,
  payload: {
    name,
    id,
    email,
  }
});
export type LoginSuccess = ReturnType<typeof loginSuccess>;

const loginError = () => ({
  type: ActionTypes.LoginError
});
export type LoginError = ReturnType<typeof loginError>;

export const login = () => async dispatch => {
  dispatch(loginStart());

  await googleAPILoaded;

  await gapi.auth2.getAuthInstance().signIn();
  const googleUser = gapi.auth2.getAuthInstance().currentUser.get();
  const profile = googleUser.getBasicProfile();

  dispatch(loginSuccess(profile.getName(), profile.getId(), profile.getEmail()));
}
