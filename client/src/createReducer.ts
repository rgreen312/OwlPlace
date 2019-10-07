export function createReducer<T>(
  initialState: T,
  handlers
) {
  return function(state: T = initialState, action): T {
    const handler = handlers[action && action.type];
    if (!handler) {
      return state;
    }
    return handler(state, action);
  };
}
