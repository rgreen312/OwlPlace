import React, { FC } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import ScrollToTop from './components/ScrollToTop';
import CanvasPageContainer from './canvas/containers/CanvasPageContainer';

const RoutingContainer: FC = () => (
  <Router>
    <ScrollToTop>
      {/* <Header /> */}

      <Route exact path="/" component={CanvasPageContainer} />
    </ScrollToTop>
  </Router>
);


export default RoutingContainer;
