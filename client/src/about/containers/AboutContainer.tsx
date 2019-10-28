import { connect } from 'react-redux';
import { sendUpdateMessage } from '../../websocket/actions';
import About from '../components/About';

interface DispatchProps {
  sendUpdateMessage: (id, x, y, r, g, b) => void;
}

const mapDispatchToProps: DispatchProps = {
  sendUpdateMessage: sendUpdateMessage
}

export default connect(
  null,
  mapDispatchToProps
)(About);