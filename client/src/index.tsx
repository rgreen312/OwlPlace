import React from 'react';
import { render } from 'react-dom'
import './index.scss';
import App from './App';
import { combineReducers, createStore, applyMiddleware, compose } from 'redux';
import canvasReducers from './canvas/reducers';
import loginReducers from './login/reducers';
import webSocketReducers from './websocket/reducers';
import { Provider } from 'react-redux';
import thunk from 'redux-thunk';

const rootReducer = combineReducers({
  canvas: canvasReducers,
  login: loginReducers,
  websocket: webSocketReducers,
});

// @ts-ignore - redux devtools doesn't have type definitions
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
export const store = createStore(rootReducer, composeEnhancers(applyMiddleware(thunk)));

render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('root')
)
