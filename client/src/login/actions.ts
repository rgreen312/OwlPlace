import * as ActionTypes from './actionTypes';

const loginStart = () => ({
  type: ActionTypes.LoginStart
});
export type LoginStart = ReturnType<typeof loginStart>;

const loginSuccess = (name: string, id: string, email: string) => ({
  type: ActionTypes.LoginSuccess
});
export type LoginSuccess = ReturnType<typeof loginSuccess>;

const loginError = () => ({
  type: ActionTypes.LoginError
});
export type LoginError = ReturnType<typeof loginError>;

export const login = () => dispatch => {
  dispatch(loginStart());

  /**
   * The Sign-In client object.
   */
  let auth2: any;

  gapi.load('auth2', function() {
    /**
     * Retrieve the singleton for the GoogleAuth library and set up the
     * client.
     */
    auth2 = gapi.auth2.init({
        client_id: '634069824484-ch6gklc2fevg9852aohe6sv2ctq7icbk.apps.googleusercontent.com'
    });
  });

    gapi.auth2.getAuthInstance().signIn().then( function() {
        const googleUser = gapi.auth2.getAuthInstance().currentUser.get();
        var profile = googleUser.getBasicProfile();
        console.log('ID: ' + profile.getId()); // Do not send to your backend! Use an ID token instead.
        console.log('Name: ' + profile.getName());
        console.log('Email: ' + profile.getEmail()); // This is null if the 'email' scope is
      }
    ); 
}
