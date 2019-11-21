import React, { FC, useState, useEffect } from 'react';
import { PageHeader, Button, Menu, Dropdown, Icon, notification } from 'antd';
import { Link, useLocation } from 'react-router-dom';
import './Header.scss';

interface Props {
  isLoggedIn: boolean;
  name?: string;
  onLogin: () => void;
  onLogout: () => void;
  timeToNextChange: number;
  setTimeRemaining: (time: number) => void;
}

const Header: FC<Props> = ({
  onLogin,
  isLoggedIn,
  name,
  onLogout,
  timeToNextChange,
  setTimeRemaining
}) => {
  //@ts-ignore
  window.onGoogleScriptLoad = () => {
    console.log('The google script has really loaded, cool!');
  };

  const loginButton = isLoggedIn ? (
    <div className='name-label' key='name'>Hi, {name}</div>
  ) : (
    <Button className='login-button' onClick={onLogin} key='login'>
      Login
    </Button>
  );

  let location = useLocation().pathname;
  
  const menu = (
    <Menu>
      <Menu.Item>
        {/*Check the last 5 characters in a string.*/}
        {location.substring(location.length - 5, location.length) !== 'about' ? (
          <Link to='/about'>About</Link>
        ) : (
          <Link to='/'>Home</Link>
        )}
      </Menu.Item>

      {isLoggedIn && <Menu.Item onClick={onLogout}>Sign Out</Menu.Item>}

      {window.location.hostname === 'localhost' && <Menu.Item><Link to='/testing'>Testing</Link></Menu.Item>}
    </Menu>
  );

  const dropdownMenu = (
    <Dropdown key='more' overlay={menu}>
      <Button
        style={{
          border: 'none',
          padding: 0
        }}
      >
        <Icon
          type='ellipsis'
          style={{
            fontSize: 20,
            verticalAlign: 'top'
          }}
        />
      </Button>
    </Dropdown>
  );

  const [time, setTime] = useState(timeToNextChange);

  useEffect(() => {
    setTime(timeToNextChange);
  }, [timeToNextChange])

  useEffect(() => {
    // @ts-ignore
    let interval: null | Timeout = null;
    if (time <= 0) {
      clearInterval(interval);
      setTimeRemaining(0);
      notification.success({
        message: 'Place another pixel',
        description: 'You can now update another pixel on the canvas.'
      });
      return;
    }

    interval = setInterval(() => {
      setTime(time => time - 1000);
    }, 1000);

    return () => clearInterval(interval);
  }, [time, setTimeRemaining]);

  const timerComponent = (
    <div className='timer'>
      {Math.floor(time / 60000)}:
      {/* Get the number of seconds remaining, and force a trailing 0 if necessary */}
      {('0' + Math.floor((time % 60000) / 1000)).slice(-2)}
    </div>
  );

  return (
    <PageHeader
      title='OwlPlace'
      subTitle='change the canvas one pixel at a time'
      extra={[isLoggedIn ? timerComponent : null, loginButton, dropdownMenu]}
    />
  );
};

export default Header;
