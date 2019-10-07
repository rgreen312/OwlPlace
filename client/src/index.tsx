import React from 'react';
import { render } from 'react-dom'
import './index.scss';
import App from './App';
import { combineReducers, createStore, applyMiddleware } from 'redux';
import canvasReducers from './canvas/reducers';
import { Provider } from 'react-redux';
import thunk from 'redux-thunk';
import { testAction } from './canvas/actions';

const rootReducer = combineReducers({
  canvas: canvasReducers
});

const store = createStore(canvasReducers, applyMiddleware(thunk));
// store.dispatch(testAction());

render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('root')
)
