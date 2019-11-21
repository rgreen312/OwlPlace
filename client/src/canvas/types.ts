export interface Color {
  r: number;
  g: number;
  b: number;
}

export interface Position {
  x: number;
  y: number;
}

export interface Move {
  position: Position;
  color: Color;
}
