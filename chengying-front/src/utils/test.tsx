import * as React from 'react';
import { render } from '@testing-library/react';
import { createStore, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';
import thunkMiddleware from 'redux-thunk';

// this is a handy function that I normally make available for all my tests
// that deal with connected components.
// you can provide initialState for the entire store that the ui is rendered with
export function renderWithRedux(
  Component,
  rootReducer,
  initialState: any = {}
) {
  const store: any = createStore(
    rootReducer,
    initialState,
    applyMiddleware(thunkMiddleware)
  );
  const layout = <Provider store={store}>{Component}</Provider>;
  return {
    ...render(layout),
    // adding `store` to the returned utilities to allow us
    // to reference it in our tests (just try to avoid using
    // this to test implementation details).
    store,
  };
}
