import React, { FC } from 'react';

interface Props {
  onLogin: () => void;
}

// TODO(Ryan): Should add a isLoggedIn prop, if they are then display user's name
const Header: FC<Props> = ({ onLogin }) => (
  <div className='header'>
    <button className='login-button' onClick={onLogin}>Login</button>
  </div>
);

export default Header;
