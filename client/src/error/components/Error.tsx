import React, { Component } from "react";
import { Redirect } from "react-router-dom";
import "./Error.scss";

interface Props {}

interface State {}

class Error extends Component<Props, State> {
  shouldRedirect = false;

  constructor(props) {
    super(props);
  }

  componentDidMount() {
    window.addEventListener("beforeunload", () => {
      this.shouldRedirect = true;
      this.forceUpdate();
    });
  }

  render() {
    return this.shouldRedirect ? (
      <Redirect to="/" />
    ) : (
      <div className="error-page">
        <h1>Whoops, something went wrong.</h1>
        <p>Hang tight, we're working hard to fix the issue!</p>
      </div>
    );
  }
}

export default Error;
