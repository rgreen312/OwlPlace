import { connect } from 'react-redux';
import Canvas from '../components/Canvas';
import { receivedError, getIsLoadingState } from '../../websocket/selectors';
import {
  registerCanvasContext,
  updateCursorPosition,
  clearCursorPosition,
  setZoomFactor,
  updatePixel,
} from '../actions';
import {
  getZoomFactor,
  getCurrentPosition,
  getInitialImage,
  canUpdatePixel,
} from '../selectors';
import { Color } from '../types';
import { getIsLoggedIn } from '../../login/selectors';

interface DispatchProps {
  registerContext: (context: CanvasRenderingContext2D) => void;
  updatePosition: (x: number, y: number) => void;
  onMouseOut: () => void;
  setZoomFactor: (newZoom: number) => void;
  onUpdatePixel: (newColor: Color, x: number, y: number) => void;
}

interface StateProps {
  receivedError: boolean;
  zoomFactor: number;
  position: { x: number; y: number } | undefined;
  isLoading: boolean;
  initialImage?: string;
  canUpdatePixel: boolean;
  isLoggedIn: boolean; 
}

const mapDispatchToProps: DispatchProps = {
  registerContext: registerCanvasContext,
  updatePosition: updateCursorPosition,
  onMouseOut: clearCursorPosition,
  setZoomFactor: setZoomFactor,
  onUpdatePixel: updatePixel
};

const mapStateToProps = (state): StateProps => ({
  receivedError: receivedError(state),
  zoomFactor: getZoomFactor(state),
  position: getCurrentPosition(state),
  isLoading: getIsLoadingState(state),
  initialImage: getInitialImage(state),
  canUpdatePixel: canUpdatePixel(state),
  isLoggedIn: getIsLoggedIn(state)
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Canvas);
