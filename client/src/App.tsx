import React from "react";
import "./App.scss";

const App: React.FC = () => {
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
      </div>
    </div>
  );
};

export default App;
