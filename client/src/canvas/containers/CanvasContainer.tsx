import { connect } from 'react-redux';
import Canvas from '../components/Canvas';
import { receivedError } from '../../websocket/selectors';

interface DispatchProps {
  registerContext: (context: CanvasRenderingContext2D) => void;
}

interface StateProps {
  receivedError: boolean;
  initialImage: string;
}

const mapDispatchToProps: DispatchProps = {
  registerContext: 
};

const mapStateToProps  = (state): StateProps => ({
  receivedError: receivedError(state),
  initialImage: 
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Canvas);
