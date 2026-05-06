'use client';
import { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { hydrateFromLocalStorage } from '../redux/slices/uiPreferencesSlice';

/**
 * Hydrates UI preferences from localStorage on client mount
 * This prevents hydration mismatches between server and client
 */
const UIPreferencesHydrator = () => {
  const dispatch = useDispatch();

  useEffect(() => {
    // Only run on client after mount
    
    // dispatch(hydrateFromLocalStorage());
  }, [dispatch]);

  return null;
};

export default UIPreferencesHydrator;