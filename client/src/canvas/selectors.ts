import { State } from '../types';

export const getCanvasContext = (
  state: State
): CanvasRenderingContext2D | undefined =>
  state.canvas.canvasContext || undefined;

export const getCurrentPosition = (
  state: State
): { x: number; y: number } | undefined =>
  state.canvas.curPosition || undefined;

export const getZoomFactor = (state: State): number => state.canvas.zoomFactor;
