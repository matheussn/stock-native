export const MOVEMENT_CREATE_ACCESS_KEY = 'movement-create-access';

export function allowMovementCreateAccess() {
  try {
    window.sessionStorage.setItem(MOVEMENT_CREATE_ACCESS_KEY, '1');
  } catch {}
}

export function clearMovementCreateAccess() {
  try {
    window.sessionStorage.removeItem(MOVEMENT_CREATE_ACCESS_KEY);
  } catch {}
}

export function hasMovementCreateAccess() {
  try {
    return window.sessionStorage.getItem(MOVEMENT_CREATE_ACCESS_KEY) === '1';
  } catch {
    return false;
  }
}
