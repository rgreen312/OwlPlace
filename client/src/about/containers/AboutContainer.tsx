import { connect } from 'react-redux';
import { sendUpdateMessage, sendLoginMessage } from '../../websocket/actions';
import About from '../components/About';

interface DispatchProps {
  sendUpdateMessage: (id, x, y, r, g, b) => void;
  sendLoginMessage: (email) => void;
}

const mapDispatchToProps: DispatchProps = {
  sendUpdateMessage: sendUpdateMessage,
  sendLoginMessage: sendLoginMessage
}

export default connect(
  null,
  mapDispatchToProps
)(About);