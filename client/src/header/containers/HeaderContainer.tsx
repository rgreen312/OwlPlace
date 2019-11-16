import { connect } from 'react-redux';
import Header from '../components/Header';
import { login, signOut } from '../../login/actions';
import { getIsLoggedIn, getUserName } from '../../login/selectors';
import { getTimeToChange } from '../../canvas/selectors';
import { setTimeToNextMove } from '../../canvas/actions';
import { getCoolDown } from '../../websocket/selectors';

interface DispatchProps {
  onLogin: () => void;
  onLogout: () => void;
  setTimeRemaining: (time: number) => void;
}

interface StateProps {
  isLoggedIn: boolean;
  name?: string;
  timeToNextChange: number;
  cooldown: number | null; 
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
  cooldown: getCoolDown(state)
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Header);
