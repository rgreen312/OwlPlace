import React, { FC } from 'react';
import { HashRouter as Router, Route } from 'react-router-dom';
import ScrollToTop from './components/ScrollToTop';
import CanvasPageContainer from './canvas/containers/CanvasPageContainer';
import AboutPage from './about/containers/AboutContainer';
import Header from './header/containers/HeaderContainer';
import Footer from './footer/components/Footer';
import ErrorPage from './error/components/Error';

const RoutingContainer: FC = () => (
  <Router basename='/OwlPlace'>
    <ScrollToTop>
      <Header />
      <Route exact path='/' component={CanvasPageContainer} />
      <Route path='/about' component={AboutPage} />
      <Route path='/error' component={ErrorPage} />
      <Footer />
    </ScrollToTop>
  </Router>
);

export default RoutingContainer;
