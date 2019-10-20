import { State } from '../types';

export const getCanvasContext = (state: State): CanvasRenderingContext2D | undefined => 
  state.canvas.canvasContext || undefined;
