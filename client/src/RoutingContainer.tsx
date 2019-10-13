import React, { FC } from "react";
import { BrowserRouter as Router, Route } from "react-router-dom";
import ScrollToTop from "./components/ScrollToTop";
import CanvasPageContainer from "./canvas/containers/CanvasPageContainer";
import AboutPage from "./about/components/About";
import Header from "./header/containers/HeaderContainer";
import Footer from "./footer/components/Footer";
import ColorPicker from "./colorPicker/colorPicker";

const RoutingContainer: FC = () => (
  <Router>
    <ScrollToTop>
      <Header />
      <ColorPicker />
      <Route exact path="/" component={CanvasPageContainer} />
      <Route path="/about" component={AboutPage} />
      <Footer />
    </ScrollToTop>
  </Router>
);

export default RoutingContainer;
