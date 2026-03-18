import {createContext} from 'react';

export const AppNavigationGuardContext = createContext({
  registerBeforeLeaveHandler: () => {},
  requestNavigation: () => {}
});
