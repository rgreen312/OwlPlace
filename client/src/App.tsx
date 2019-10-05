import React from "react";
import "./App.scss";

const App: React.FC = () => {

  var ws: WebSocket | null;
  
  const onClickOpen = () => {
    console.log("opening")
        if (ws) {
            return false;
        }
        ws = new WebSocket("ws://localhost:3000/");
        ws.onopen = function(evt) {
            console.log("OPEN");
        }
        ws.onclose = function(evt) {
            console.log("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            console.log("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            console.log("ERROR! ");
        }
        // return false;
  }
    const onClickP1 = () => {

    }

    const onClickClose = () => {
      if (!ws) {
        return false;
      }
      ws.close();
      return false;
    }

  return (
    <div>
      <div className="top-nav-bar">
        <button className="login-btn">
          <p className="login-text">login</p>
        </button>
      </div>
      <div className="main-wrapper">
        <h1>owlplaces</h1>
        <h2>change the canvas one pixel at a time</h2>
        <p>Click "Open" to create a connection to the server, 
          "Send" to send a message to the server and "Close" to close the connection. 
          You can change the message and send multiple times.
          </p>
          <button onClick= {onClickOpen} id="open"> Open </button>
          <button onClick= {onClickP1} id="p1"> Pixel 1 </button>
          <button onClick= {onClickClose} id="close">Close</button>
        
      </div>

      
    </div>
    
  );
};

// <script>  
// window.addEventListener("load", function(evt) {
//     var output = document.getElementById("output");
//     var input = document.getElementById("input");
//     var ws;
//     var print = function(message) {
//         var d = document.createElement("div");
//         d.innerHTML = message;
//         output.appendChild(d);
//     };
//     document.getElementById("open").onclick = function(evt) {
//         if (ws) {
//             return false;
//         }
//         ws = new WebSocket("{{.}}");
//         ws.onopen = function(evt) {
//             print("OPEN");
//         }
//         ws.onclose = function(evt) {
//             print("CLOSE");
//             ws = null;
//         }
//         ws.onmessage = function(evt) {
//             print("RESPONSE: " + evt.data);
//         }
//         ws.onerror = function(evt) {
//             print("ERROR: " + evt.data);
//         }
//         return false;
//     };
//     document.getElementById("send").onclick = function(evt) {
//         if (!ws) {
//             return false;
//         }
//         print("SEND: " + input.value);
//         ws.send(input.value);
//         return false;
//     };
//     document.getElementById("close").onclick = function(evt) {
//         if (!ws) {
//             return false;
//         }
//         ws.close();
//         return false;
//     };
// });
// </script>



export default App;
