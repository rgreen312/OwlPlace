import { connect } from 'react-redux';
import Canvas from '../components/Canvas';
import { receivedError } from '../../websocket/selectors';
import { registerCanvasContext, updateCursorPosition, clearCursorPosition, setZoomFactor } from '../actions';
import { getZoomFactor, getCurrentPosition } from '../selectors';

interface DispatchProps {
  registerContext: (context: CanvasRenderingContext2D) => void;
  updatePosition: (x: number, y: number) => void;
  onMouseOut: () => void;
  setZoomFactor: (newZoom: number) => void;
}

interface StateProps {
  receivedError: boolean;
  zoomFactor: number;
  position: {x: number, y: number} | undefined; 
}

const mapDispatchToProps: DispatchProps = {
  registerContext: registerCanvasContext,
  updatePosition: updateCursorPosition,
  onMouseOut: clearCursorPosition,
  setZoomFactor: setZoomFactor
};

const mapStateToProps  = (state): StateProps => ({
  receivedError: receivedError(state),
  zoomFactor: getZoomFactor(state),
  position: getCurrentPosition(state), 
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Canvas);
