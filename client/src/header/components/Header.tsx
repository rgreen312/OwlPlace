import React, { FC } from 'react';
import { PageHeader, Button } from 'antd';

interface Props {
  onLogin: () => void;
  isLoggedIn: boolean;
  name?: string;
}

// TODO(Ryan): Should add a isLoggedIn prop, if they are then display user's name
const Header: FC<Props> = ({ onLogin, isLoggedIn, name }) => {
  console.log('logged in: ' + isLoggedIn)
  const loginButton = isLoggedIn
  ? (
    <>
      Hi, {name}
    </>
  ) : (
    <Button className='login-button' onClick={onLogin}>
      Login
    </Button>
  );

  return (
    <PageHeader
      title='OwlPlace'
      extra={loginButton}
    />
  )
}

export default Header;
