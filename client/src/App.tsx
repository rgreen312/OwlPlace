import React from "react";
import "./App.scss";
import RoutingContainer from './RoutingContainer';
import { connect } from 'react-redux';
import { checkLogin } from './login/actions';
import { openWebSocket } from './websocket/actions';

// // updateMessage is type 1
// const updateMessage = (
//   id: string,
//   x: number,
//   y: number,
//   r: number,
//   g: number,
//   b: number
// ) => {
//   return JSON.stringify({
//     type: 1,
//     userId: id,
//     x: x,
//     y: y,
//     r: r,
//     g: g,
//     b: b
//   });
// };
// const onClickP1 = (
//   id: string,
//   x: number,
//   y: number,
//   r: number,
//   g: number,
//   b: number
// ) => {
//   console.log("Sending update of Pixel 1");
//   socket.send(updateMessage(id, x, y, r, g, b));

//   return true;
// }

interface Props {
  checkLogin: () => void;
  openConnection: () => void;

}

class App extends React.Component<Props> {


  componentDidMount() {
    this.props.checkLogin();
    this.props.openConnection();
  }
  
  render() {
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
  }
};


const mapDispatchToProps: Props = {
  checkLogin: checkLogin,
  openConnection: openWebSocket

}

export default connect(
  null,
  mapDispatchToProps,
)(App);
