import { State } from '../types';

export const getIsLoggedIn = (state: State): boolean => state.login.isLoggedIn;
export const getUserName = (state: State): string | undefined => state.login.name || undefined;
export const getUserId = (state: State): string | undefined => state.login.userId || undefined;
export const getUserEmail = (state: State): string | undefined => state.login.email || undefined;
