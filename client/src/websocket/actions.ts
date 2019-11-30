import { HOSTNAME } from '../constants';
import * as ActionTypes from './actionTypes';
import * as CanvasActionTypes from '../canvas/actionTypes';
import { getWebSocket } from './selectors';
import { MsgType } from '../message';
import { setImage, setTimeToNextMove } from '../canvas/actions';
import { getCanvasContext, getLastMove } from '../../src/canvas/selectors';
import { getUserEmail } from '../login/selectors';
const startConnect = () => ({
  type: ActionTypes.StartConnect
});
export type StartConnect =  ReturnType<typeof startConnect>;

export const connectError = (error: string) => ({
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

export const openWebSocket = () => (dispatch, getState) => {
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

      if (getUserEmail(getState())) {
        socket.send(makeLoginMessage(getUserEmail(getState())!)); 
      }
      
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
    console.log("RECEIVED: " + json.type);
    switch (json.type) {
        case MsgType.IMAGE: {
            // let msg = new ImageMsg(data.formatString); //now what?
            let imageString = json.formatString
            console.log("Received an IMAGE message from the server!");
            dispatch(setImage('data:image/png;base64,' + imageString));
            break;
        }
        case MsgType.TESTING: {
          console.log("Received a TESTING message from the server!");
          console.log("Message: " + json.msg);
          break;
        }
        case MsgType.CHANGECLIENTPIXEL: {
          console.log("Received a CHANGECLIENTPIXEL message from the server!");
          let x = json.x
          let y = json.y
          // let color = { r: json.r, g: json.g, b: json.b }
          dispatch(setColor(x, y, json.r, json.g, json.b))
          break;
        }
        case MsgType.DRAWRESPONSE: {
          let status = json.status
          console.log("Received a DRAWRESPONSE message from the server!");
          console.log("The status was " + status);
          if (status === 503) {
            // If update fails, we reset the pixel to its previous color
            const prevMove = getLastMove(getState());
            if (prevMove) {
              const { x, y } = prevMove.position;
              const { r, g, b } = prevMove.color;
              dispatch(setColor(x, y, r, g, b));
              dispatch({ type: CanvasActionTypes.UpdatePixelError });
            }
          } else if (status == 429) {
            // If the user's cooldown hasn't expired yet
            if (json.remainingTime > 0) {
              dispatch(setTimeToNextMove(json.remainingTime));
            } else {
              // this might happen if someone manually makes a status message with code 429...
              console.log("Received an ill-formatted DRAWRESPONSE cooldown message from the server!");
            }
          } else {
            dispatch({ type: CanvasActionTypes.UpdatePixelSuccess });
          }
          break;
        }
        case MsgType.VERIFICATIONFAIL: {
          let status = json.status
          console.log("Received a VERIFICATIONFAIL message from the server!");
          console.log("The status was " + status);
          // If user verification fails, direct to error page
          dispatch({ type: ActionTypes.ConnectError });
          break;
        }
        case MsgType.USERLOGINRESPONSE: {
          let status = json.status
          let cooldown = json.cooldown
          console.log("The status was " + status);
          console.log("The remaining cooldown time for current user is: " + cooldown);
          dispatch(setTimeToNextMove(cooldown));
          break;
        }
        default: {
            console.log("Received a message from the server of an unknown type, message: " + data + " type: " + json.type) ;
            break;
        }
    }
    // TODO (Ryan): figure out the best way to handle this... probably need to write some middlewear

  };
}

export const makeUpdateMessage = (
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

const setColor = (x: number, y: number, r: number, g: number, b: number) => (dispatch, getState) => {
  console.log("inside set color in actions.ts")
  const state = getState();
  const ctx = getCanvasContext(state);
  if (ctx) {
    ctx!.fillStyle = 'rgb('+ r + ',' + g + ',' + b + ')';
    ctx!.fillRect(x , y , 1, 1);
  }
}

export const makeLoginMessage = (
  email: string
) => {
  return JSON.stringify({
    type: 2, 
    email: email
  })
}

export const sendUpdateMessage = (id, x, y, r, g, b) => (dispatch, getState) => {
  const socket = getWebSocket(getState());
  console.log("sending along websocket message")
  if (socket) {
    socket.send(makeUpdateMessage(id, x, y, r, g, b));

    // The follwing should be REMOVED when testing is done/ you want to only do single pixels
    // let lower = 495;
    // let upper = 505;
    // for (let i = lower; i < upper; i++) {
    //   for (let j = lower; j < upper; j++) {
    //     console.log("sending..")
    //     socket.send(makeUpdateMessage(id, i, j, r, g, b));
    //   }
    // }
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
