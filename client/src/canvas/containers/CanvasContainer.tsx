import { connect } from 'react-redux';
import Canvas from '../components/Canvas';
import { receivedError } from '../../websocket/selectors';
import { registerCanvasContext } from '../actions';

interface DispatchProps {
  registerContext: (context: CanvasRenderingContext2D) => void;
}

interface StateProps {
  receivedError: boolean;
}

const mapDispatchToProps: DispatchProps = {
  registerContext: registerCanvasContext,
};

const mapStateToProps  = (state): StateProps => ({
  receivedError: receivedError(state),
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Canvas);
