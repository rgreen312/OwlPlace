import React, { FC } from 'react';
import { combineReducers, createStore } from 'redux';
import canvasReducers from './canvas/reducers';
import { Provider } from 'react-redux';

const rootReducer = combineReducers({
  canvas: canvasReducers
});

const store = createStore(rootReducer);

const RootContainer: FC = ({ children }) => (
  <Provider store={store}>
    {children}
  </Provider>
)
