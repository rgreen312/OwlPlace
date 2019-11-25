import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import './Error.scss';

interface Props {}

interface State {}

class Error extends Component<Props, State> {
  shouldRedirect = false;

  componentDidMount() {
    window.addEventListener('beforeunload', () => {
      this.shouldRedirect = true;
      this.forceUpdate();
    });
  }

  render() {
    return this.shouldRedirect ? (
      <Redirect to='/' />
    ) : (
      <div className='error-page'>
        <div className='error-title'>Whoops, something went wrong.</div>
        <p>Hang tight, we're working hard to fix the issue!</p>
      </div>
    );
  }
}

export default Error;
