import { connect } from 'react-redux';
import Header from '../components/Header';
import { login, signOut } from '../../login/actions';
import { getIsLoggedIn, getUserName } from '../../login/selectors';
import { getTimeToChange } from '../../canvas/selectors';
import { setTimeToNextMove } from '../../canvas/actions';

interface DispatchProps {
  onLogin: () => void;
  onLogout: () => void;
  setTimeRemaining: (time: number) => void;
}

interface StateProps {
  isLoggedIn: boolean;
  name?: string;
  timeToNextChange: number;
}

const mapDispatchToProps: DispatchProps = {
  onLogin: login,
  onLogout: signOut,
  setTimeRemaining: setTimeToNextMove,
}

const mapStateToProps = (state): StateProps => ({
  isLoggedIn: getIsLoggedIn(state),
  name: getUserName(state),
  timeToNextChange: getTimeToChange(state),
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Header);
