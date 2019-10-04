import React from 'react';
import "./App.scss";

/**
 * The Sign-In client object.
 */
var auth2: any;

/**
 * Initializes the Sign-In client.
 */
gapi.load('auth2', function() {
  /**
   * Retrieve the singleton for the GoogleAuth library and set up the
   * client.
   */
  auth2 = gapi.auth2.init({
      client_id: '634069824484-ch6gklc2fevg9852aohe6sv2ctq7icbk.apps.googleusercontent.com'
  });
});

function onSignIn() {
  gapi.auth2.getAuthInstance().signIn().then( function() {
      const googleUser = gapi.auth2.getAuthInstance().currentUser.get();
      var profile = googleUser.getBasicProfile();
      console.log('ID: ' + profile.getId()); // Do not send to your backend! Use an ID token instead.
      console.log('Name: ' + profile.getName());
      console.log('Email: ' + profile.getEmail()); // This is null if the 'email' scope is
    }
  ); 
}

const App: React.FC = () => {
  return (
    <div>
      <div className="top-nav-bar">
        <button className="login-btn" onClick={onSignIn}>
          <p className="login-text">login</p>
        </button>
      </div>
      <div className="main-wrapper">
        <h1>owlplaces</h1>
        <h2>change the canvas one pixel at a time</h2>
      </div>
    </div>
  );
};

export default App;
