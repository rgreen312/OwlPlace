import React, { FC } from "react";
import "./App.scss";
import RoutingContainer from './RoutingContainer';
import * as ActionTypes from './login/actionTypes';
import { loginStart, loginSuccess } from "./login/actions";

let socket = new WebSocket("ws://127.0.0.1:3010/ws");
console.log("Attempting Connection...");

// open message is 0
socket.onopen = () => {
  console.log("Successfully Connected");
  socket.send(
    JSON.stringify({
      type: 0,
      message: "Hi From the Client! The websocket just opened"
    })
  );
};

// close message is 9
socket.onclose = event => {
  console.log("Socket Closed Connection: ", event);
  socket.send(
    JSON.stringify({
      type: 9,
      message: "Client Closed!"
    })
  );
};

socket.onerror = error => {
  console.log("Socket Error: ", error);
};

socket.onmessage = event => {
  var message = event.data;
  console.log("Recieved a message from the server, message: " + message);
};

// updateMessage is type 1
const updateMessage = (
  id: string,
  x: number,
  y: number,
  r: number,
  g: number,
  b: number
) => {
  return JSON.stringify({
    type: 1,
    userId: id,
    x: x,
    y: y,
    r: r,
    g: g,
    b: b
  });
};
const onClickP1 = (
  id: string,
  x: number,
  y: number,
  r: number,
  g: number,
  b: number
) => {
  console.log("Sending update of Pixel 1");
  socket.send(updateMessage(id, x, y, r, g, b));

  return true;
}

/**
 * The Sign-In client object.
 */
let auth2: any;
let googleUser: any;

export const googleAPILoaded: Promise<void> = new Promise(resolve => {
  gapi.load('auth2', () => {
    /**
     * Retrieve the singleton for the GoogleAuth library and set up the
     * client.
     */
    gapi.auth2.init({
        client_id: '634069824484-ch6gklc2fevg9852aohe6sv2ctq7icbk.apps.googleusercontent.com'
    }).then( function() {
        // Sign in the user if they are currently signed in.
        auth2 = gapi.auth2.getAuthInstance(); 
        if (auth2.isSignedIn.get() == true) {
          console.log("someone is already signed in"); 
          const profile = googleUser.getBasicProfile();
          
          const login = () => async dispatch => {
            dispatch(loginStart());
          
            dispatch(loginSuccess(profile.getName(), profile.getId(), profile.getEmail()));
          }
        } 
      }
    );

    resolve();
  });
});


const App: FC = () => {
  return (
    // <div>
    //   <div className="top-nav-bar">
    //     <button className="login-btn" onClick={onSignIn}>
    //       <p className="login-text">login</p>
    //     </button>
    //   </div>
    //   <div className="main-wrapper">
    //     <h1>owlplaces</h1>
    //     <h2>change the canvas one pixel at a time</h2>
    //     <p>Click "Pixel 1" to send an update message to the server!
    //       </p>
    //       <button onClick = {() => onClickP1("user1", 10, 400, 255, 255, 255)} id="p1"> Pixel 1 </button>
    //       {/* <button onClick= {onClickClose} id="close">Close</button> */}
    //   </div>
    // </div>
    <RoutingContainer />
  );
};

export default App;
