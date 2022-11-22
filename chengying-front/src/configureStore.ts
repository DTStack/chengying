import { createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import { Service } from './services';
import rootStore from './stores';

export default function configureStore(initialState: any) {
  thunkMiddleware.withExtraArgument(Service);
  return createStore(
    rootStore,
    initialState,
    compose(
      applyMiddleware(thunkMiddleware),
      (window as any).__REDUX_DEVTOOLS_EXTENSION__
        ? (window as any).__REDUX_DEVTOOLS_EXTENSION__ &&
            (window as any).__REDUX_DEVTOOLS_EXTENSION__()
        : (fn: any) => fn
    )
  );
}
