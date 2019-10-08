import React, { FC } from "react";
import "./App.scss";
import RoutingContainer from './RoutingContainer';

const App: FC = () => {
  return (
    // <div>
    //   <div className="top-nav-bar">
    //     <button className="login-btn" onClick={onSignIn}>
    //       <p className="login-text">login</p>
    //     </button>
    //   </div>
    //   <div className="main-wrapper">
    //     <h1>owlplaces</h1>
    //     <h2>change the canvas one pixel at a time</h2>
    //     <p>Click "Pixel 1" to send an update message to the server!
    //       </p>
    //       <button onClick = {() => onClickP1("user1", 10, 400, 255, 255, 255)} id="p1"> Pixel 1 </button>
    //       {/* <button onClick= {onClickClose} id="close">Close</button> */}
    //   </div>
    // </div>
    <RoutingContainer />
  );
};

export default App;
