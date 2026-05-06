import { act } from '@testing-library/react';
/**
 * Wraps a state update function in act() and returns a promise that resolves when the update is complete
 * @param {Function} updateFn - The function that triggers the state update
 * @returns {Promise} A promise that resolves when the update is complete
 */
export const wrapStateUpdate = async (updateFn) => {
  await act(async () => {
    await updateFn();
  });
};
/**
 * Wraps multiple state updates in act() and returns a promise that resolves when all updates are complete
 * @param {Array<Function>} updateFns - Array of functions that trigger state updates
 * @returns {Promise} A promise that resolves when all updates are complete
 */
export const wrapMultipleStateUpdates = async (updateFns) => {
  await act(async () => {
    await Promise.all(updateFns.map(fn => fn()));
  });
};
/**
 * Creates a mock function that wraps state updates in act()
 * @param {Function} fn - The function to wrap
 * @returns {Function} A wrapped function that automatically uses act()
 */
export const createActWrappedMock = (fn) => {
  return async (...args) => {
    await act(async () => {
      await fn(...args);
    });
  };
}; 