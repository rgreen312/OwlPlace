import { HOSTNAME } from '../constants';
import * as ActionTypes from './actionTypes';

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
    dispatch(connectSuccess(socket));
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
    dispatch(connectError(error.type));
  };

  socket.onmessage = event => {
    const { data } = event;
    // TODO (Ryan): figure out the best way to handle this... probably need to write some middlewear
    console.log("Recieved a message from the server, message: " + data);
  };
}

// TODO (ryan): create actions to disconnect from web socket and send messages
