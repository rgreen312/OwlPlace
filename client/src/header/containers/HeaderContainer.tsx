import { connect } from 'react-redux';
import Header from '../components/Header';
import { login, signOut } from '../../login/actions';
import { getIsLoggedIn, getUserName } from '../../login/selectors';

interface DispatchProps {
  onLogin: () => void;
  onLogout: () => void;
  onError: (error) => void;
}

interface StateProps {
  isLoggedIn: boolean;
  name?: string;
}

const mapDispatchToProps: DispatchProps = {
  onLogin: login,
  onLogout: signOut,
  onWebSocketError: webSocketError
}

const mapStateToProps  = (state): StateProps => ({
  isLoggedIn: getIsLoggedIn(state),
  name: getUserName(state)
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Header);
