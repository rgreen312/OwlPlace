import { connect } from 'react-redux';
import Header from '../components/Header';
import { login } from '../../login/actions';
import { getIsLoggedIn, getUserName } from '../../login/selectors';

interface DispatchProps {
  onLogin: () => void;
}

interface StateProps {
  isLoggedIn: boolean;
  name?: string;
}

const mapDispatchToProps: DispatchProps = {
  onLogin: login
}

const mapStateToProps  = (state): StateProps => ({
  isLoggedIn: getIsLoggedIn(state),
  name: getUserName(state)
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Header);
