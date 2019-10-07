import React from 'react';
import { render } from 'react-dom'
import './index.scss';
import App from './App';
import { combineReducers, createStore, applyMiddleware } from 'redux';
import canvasReducers from './canvas/reducers';
import loginReducers from './login/reducers';
import { Provider } from 'react-redux';
import thunk from 'redux-thunk';

const rootReducer = combineReducers({
  canvas: canvasReducers,
  login: loginReducers
});

const store = createStore(rootReducer, applyMiddleware(thunk));

render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('root')
)
