import React, { FC } from 'react';
import { PageHeader, Button, Menu, Dropdown, Icon } from 'antd';
import { Link } from 'react-router-dom';

interface Props {
  onLogin: () => void;
  isLoggedIn: boolean;
  name?: string;
}

// TODO(Ryan): Should add a isLoggedIn prop, if they are then display user's name
const Header: FC<Props> = ({ onLogin, isLoggedIn, name }) => {

  //@ts-ignore
  window.onGoogleScriptLoad = () => {
    console.log('The google script has really loaded, cool!');
  }

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

  // TODO(ryan): add sign out functionality
  const menu = (
    <Menu>
      <Menu.Item>
        <Link to='/about'>About</Link>
      </Menu.Item>
      {isLoggedIn && (
        <Menu.Item>
          <button onClick={() => {}}>Sign out</button>
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
