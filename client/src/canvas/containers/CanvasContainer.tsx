import { connect } from 'react-redux';
import Canvas from '../components/Canvas';
import { testAction } from '../actions';

interface DispatchProps {
  onClick: () => void;
}

const mapDispatchToProps: DispatchProps = {
  onClick: testAction
};

export default connect(
  null,
  mapDispatchToProps
)(Canvas);
