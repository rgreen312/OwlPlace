import { HOSTNAME } from '../constants';
import * as ActionTypes from './actionTypes';
import { getWebSocket } from './selectors';
import { ERROR, IMAGE, Msg, ErrorMsg, ImageMsg, TESTING, DRAWRESPONSE } from '../message';
import { setImage } from '../canvas/actions';

const startConnect = () => ({
  type: ActionTypes.StartConnect
});
export type StartConnect =  ReturnType<typeof startConnect>;

const connectError = (error: string) => ({
  type: ActionTypes.ConnectError,
  payload: {
    error
  }
});
export type ConnectError =  ReturnType<typeof connectError>;

const closeConnection = () => ({
  type: ActionTypes.CloseConnection,
});
export type CloseConnection =  ReturnType<typeof closeConnection>;

const connectSuccess = (socket: WebSocket) => ({
  type: ActionTypes.ConnectSuccess,
  payload: {
    socket
  }
});
export type ConnectSuccess = ReturnType<typeof connectSuccess>;

export const openWebSocket = () => dispatch => {
  dispatch(startConnect());

  const socket = new WebSocket(`ws://${HOSTNAME}/ws`);
    // open message is 0
    socket.onopen = () => {
      socket.send(
        JSON.stringify({
          type: 0,
          message: "Hi From the Client! The websocket just opened"
        })
      );

      // socket.send(makeUpdateMessage("AAAAAA", 6, 9, 4, 2, 0));
      
    dispatch(connectSuccess(socket));
  };

  // close message is 9
  socket.onclose = event => {
    dispatch(closeConnection());
  };

  socket.onerror = error => {
    dispatch(connectError(error.type));
  };

  socket.onmessage = event => {
    const { data } = event;
    let json = JSON.parse(data);
    switch (json.type) {
        case IMAGE: {
            // let msg = new ImageMsg(data.formatString); //now what?
            let imageString = json.formatString
            console.log("Received an IMAGE message from the server!");
            console.log("Format string: " + imageString);
            dispatch(setImage('data:image/png;base64,' + imageString));
            break;
        }
        case TESTING: {
          console.log("Received a TESTING message from the server!");
          console.log("Message: " + json.msg);
          break;
        }
        case DRAWRESPONSE: {
          let status = json.status
          console.log("Received a DRAWRESPONSE message from the server!");
          console.log("The status was " + status)
          break;
        }
        default: {
            console.log("Received a message from the server of an unknown type, message: " + data);
            break;
        }
    }
    // TODO (Ryan): figure out the best way to handle this... probably need to write some middlewear

  };
}

const makeUpdateMessage = (
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

const makeLoginMessage = (
  email: string
) => {
  return JSON.stringify({
    type: 2, 
    email: email
  })
}

export const sendUpdateMessage = (id, x, y, r, g, b) => (dispatch, getState) => {
  const socket = getWebSocket(getState());
  if (socket) {
    socket.send(makeUpdateMessage(id, x, y, r, g, b));

    // The follwing should be REMOVED when testing is done/ you want to only do single pixels
    let lower = 495;
    let upper = 505;
    for (let i = lower; i < upper; i++) {
      for (let j = lower; j < upper; j++) {
        console.log("sending..")
        socket.send(makeUpdateMessage(id, i, j, r, g, b));
      }
    }
  }
}

export const sendLoginMessage = (email) => (dispatch, getState) => {
  console.log("sending login message for email: " + email)
  const socket = getWebSocket(getState());
  if (socket) {
    socket.send(makeLoginMessage(email));
  }
}

export const closeWebSocket = () => (dispatch, getState) => {
  const socket = getWebSocket(getState());
  if (socket) {
    socket.send(
      JSON.stringify({
        type: 9,
        message: "Client Closed!"
      })
    );
  }
  dispatch(closeConnection());
}
