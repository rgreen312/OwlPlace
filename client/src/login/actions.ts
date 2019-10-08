import * as ActionTypes from './actionTypes';

const loginStart = () => ({
  type: ActionTypes.LoginStart
});
export type LoginStart = ReturnType<typeof loginStart>;

const loginSuccess = (name: string, id: string, email: string) => ({
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

  /**
   * The Sign-In client object.
   */
  let auth2: any;

  const googleAPILoaded: Promise<void> = new Promise(resolve => {
    gapi.load('auth2', () => {
      /**
       * Retrieve the singleton for the GoogleAuth library and set up the
       * client.
       */
      auth2 = gapi.auth2.init({
          client_id: '634069824484-ch6gklc2fevg9852aohe6sv2ctq7icbk.apps.googleusercontent.com'
      });
      resolve();
    });
  });

  await googleAPILoaded;

  await gapi.auth2.getAuthInstance().signIn();
  const googleUser = gapi.auth2.getAuthInstance().currentUser.get();
  const profile = googleUser.getBasicProfile();

  dispatch(loginSuccess(profile.getName(), profile.getId(), profile.getEmail()));
}
