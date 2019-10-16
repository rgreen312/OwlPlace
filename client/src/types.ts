import { State as LoginState } from './login/reducers';
import { State as CanvasState } from './canvas/reducers';
import { State as WebSocketState } from './websocket/reducers';

export interface State {
  login: LoginState;
  canvas: CanvasState;
  websocket: WebSocketState;
}
