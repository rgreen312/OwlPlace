import { State } from '../types';

export const getWebSocket = (state: State): WebSocket | undefined => state.websocket.socket || undefined;
export const getIsConnected = (state: State): boolean => state.websocket.isConnected;
