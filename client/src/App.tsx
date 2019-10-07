import React from "react";
import Footer from "./landingPage/Footer";
import Header from "./landingPage/Header";
import { Button } from "antd";
import "./App.scss";

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
};
/**
 * The Sign-In client object.
 */
var auth2: any;

/**
 * Initializes the Sign-In client.
 */
gapi.load("auth2", function() {
  /**
   * Retrieve the singleton for the GoogleAuth library and set up the
   * client.
   */
  auth2 = gapi.auth2.init({
    client_id:
      "634069824484-ch6gklc2fevg9852aohe6sv2ctq7icbk.apps.googleusercontent.com"
  });
});

function onSignIn() {
  gapi.auth2
    .getAuthInstance()
    .signIn()
    .then(function() {
      const googleUser = gapi.auth2.getAuthInstance().currentUser.get();
      var profile = googleUser.getBasicProfile();
      console.log("ID: " + profile.getId()); // Do not send to your backend! Use an ID token instead.
      console.log("Name: " + profile.getName());
      console.log("Email: " + profile.getEmail()); // This is null if the 'email' scope is
    });
}

const App: React.FC = () => {
  return (
    <div className="main-wrapper">
      <p>Click "Pixel 1" to send an update message to the server!</p>
      <button
        onClick={() => onClickP1("user1", 10, 400, 255, 255, 255)}
        id="p1"
      >
        {" "}
        Pixel 1{" "}
      </button>
      <Header />
      <Button
        className="login-btn"
        type="primary"
        icon="google"
        onClick={onSignIn}
      >
        Login with Google
      </Button>
      <Footer />
    </div>
  );
};

export default App;
