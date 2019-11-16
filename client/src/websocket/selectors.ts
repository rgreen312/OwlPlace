import { State } from "../types";

export const getWebSocket = (state: State): WebSocket | undefined =>
  state.websocket.socket || undefined;
export const getIsConnected = (state: State): boolean =>
  state.websocket.isConnected;
export const receivedError = (state: State): boolean =>
  state.websocket.receivedError;
export const getIsLoadingState = (state: State): boolean =>
  state.websocket.isLoading;
export const getCoolDown = (state: State): number | null => 
  state.websocket.coolDown; 
