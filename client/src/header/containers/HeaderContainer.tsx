import { connect } from 'react-redux';
import Header from '../components/Header';
import { login } from '../../login/actions';

interface DispatchProps {
  onLogin: () => void;
}

const mapDispatchToProps: DispatchProps = {
  onLogin: login
}

export default connect(null, mapDispatchToProps)(Header);
