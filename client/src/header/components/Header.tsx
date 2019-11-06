import React, { FC, useState } from "react";
import { PageHeader, Button, Menu, Dropdown, Icon, Tag } from "antd";
import { Link, useLocation } from "react-router-dom";
import Timer from 'react-compound-timer';
import "./Header.scss";

interface Props {
  isLoggedIn: boolean;
  name?: string;
  onLogin: () => void;
  onLogout: () => void;
  timeToNextChange: number;
}

const Header: FC<Props> = ({ onLogin, isLoggedIn, name, onLogout, timeToNextChange }) => {
  //@ts-ignore
  window.onGoogleScriptLoad = () => {
    console.log("The google script has really loaded, cool!");
  };

  const loginButton = isLoggedIn ? (
    <div className='name-label'>Hi, {name}</div>
  ) : (
    <Button className='login-button' onClick={onLogin}>
      Login
    </Button>
  );

  let location = useLocation().pathname;

  const menu = (
    <Menu>
      <Menu.Item>
        {/*Check the last 5 characters in a string.*/}
        {location.substring(location.length - 5, location.length) != "about" ? (
          <Link to='/about'>About</Link>
        ) : (
          <Link to='/'>Home</Link>
        )}
      </Menu.Item>
      {isLoggedIn && <Menu.Item onClick={onLogout}>Sign Out</Menu.Item>}
    </Menu>
  );

  const dropdownMenu = (
    <Dropdown key='more' overlay={menu}>
      <Button
        style={{
          border: "none",
          padding: 0
        }}
      >
        <Icon
          type='ellipsis'
          style={{
            fontSize: 20,
            verticalAlign: "top"
          }}
        />
      </Button>
    </Dropdown>
  );

  const [showPrefix, setShowPrefix] = useState(false);

  const timer = (
    <div className='timer'>
      <Timer
        initialTime={timeToNextChange}
        direction='backward'
        checkpoints={[
          {
            time: 10000,
            callback: () => setShowPrefix(true)
          },
          {
            time: 0,
            callback: () => console.log('update the state here')
          }
        ]}
      >
        {() =>(
          <>
            <Timer.Minutes />:{showPrefix ? 0 : null}<Timer.Seconds />
          </>
        )}
      </Timer>
    </div>
  )

  return (
    <PageHeader
      title='OwlPlace'
      subTitle='change the canvas one pixel at a time'
      extra={[timer, loginButton, dropdownMenu]}
    />
  );
};

export default Header;
