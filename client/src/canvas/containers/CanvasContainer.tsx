import { connect } from 'react-redux';
import Canvas from '../components/Canvas';
import { receivedError } from '../../websocket/selectors';
import { registerCanvasContext, updateCursorPosition, clearCursorPosition } from '../actions';

interface DispatchProps {
  registerContext: (context: CanvasRenderingContext2D) => void;
  updatePosition: (x: number, y: number) => void;
  onMouseOut: () => void;
}

interface StateProps {
  receivedError: boolean;
}

const mapDispatchToProps: DispatchProps = {
  registerContext: registerCanvasContext,
  updatePosition: updateCursorPosition,
  onMouseOut: clearCursorPosition,
};

const mapStateToProps  = (state): StateProps => ({
  receivedError: receivedError(state),
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Canvas);
