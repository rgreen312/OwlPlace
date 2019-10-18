import { connect } from 'react-redux';
import Canvas from '../components/Canvas';
import { receivedError } from '../../websocket/selectors';

interface DispatchProps {

}

const mapDispatchToProps: DispatchProps = {
};

interface StateProps {
  receivedError: boolean; 
}
const mapStateToProps  = (state): StateProps => ({
  receivedError: receivedError(state)
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Canvas);
