import React, { FC } from 'react';
import { PageHeader, Button, Menu, Dropdown, Icon } from 'antd';
import { Link } from 'react-router-dom';

interface Props {
  isLoggedIn: boolean;
  name?: string;
  onLogin: () => void;
  onLogout: () => void;
}

// TODO(Ryan): Should add a isLoggedIn prop, if they are then display user's name
const Header: FC<Props> = ({ onLogin, isLoggedIn, name, onLogout }) => {

  //@ts-ignore
  window.onGoogleScriptLoad = () => {
    console.log('The google script has really loaded, cool!');
  }

  const loginButton = isLoggedIn
  ? (
    <div key='name' style={{ display: 'inline-block' }}>
      Hi, {name}
    </div>
  ) : (
    <Button className='login-button' onClick={onLogin} key='signin'>
      Login
    </Button>
  );

  // TODO(ryan): add sign out functionality
  const menu = (
    <Menu>
      <Menu.Item>
        <Link to='/about'>About</Link>
      </Menu.Item>
      {isLoggedIn && (
        <Menu.Item onClick={onLogout}>
          Sign Out
        </Menu.Item>
      )}
    </Menu>
  );

  const dropdownMenu = (
    <Dropdown key="more" overlay={menu}>
      <Button
        style={{
          border: 'none',
          padding: 0,
        }}
      >
        <Icon
          type="ellipsis"
          style={{
            fontSize: 20,
            verticalAlign: 'top',
          }}
        />
      </Button>
    </Dropdown>
  );

  return (
    <PageHeader
      title='OwlPlace'
      extra={[loginButton, dropdownMenu]}
    />
  )
}

export default Header;
