import { FC, useEffect } from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';

// @ts-ignore - TODO: figure out why typescript doesn't like this
const ScrollToTop: FC<RouteComponentProps> = ({
  children,
  location: { pathname }
}) => {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname]);

  return children || null;
};

export default withRouter(ScrollToTop);
