import { State as LoginState } from './login/reducers';
import { State as CanvasState } from './canvas/reducers';

export interface State {
  login: LoginState;
  canvas: CanvasState;
}
