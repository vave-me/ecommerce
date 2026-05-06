'use client';
import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { showUnifiedComposer, hideUnifiedComposer, selectShowUnifiedComposer } from '../../redux/slices/uiPreferencesSlice';

const ReduxDebugger = () => {
  const dispatch = useDispatch();
  const showComposer = useSelector(selectShowUnifiedComposer);
  const fullState = useSelector(state => state);

  return (
    <div style={{ 
      position: 'fixed', 
      bottom: 20, 
      left: 20, 
      background: 'white', 
      border: '2px solid red', 
      padding: 10, 
      zIndex: 9999 
    }}>
      <h4>Redux Debug</h4>
      <p>showComposer: {String(showComposer)}</p>
      <button onClick={() => {
        
        dispatch(showUnifiedComposer());
      }}>Show Composer</button>
      <button onClick={() => {
        
        dispatch(hideUnifiedComposer());
      }}>Hide Composer</button>
      <details>
        <summary>Full State</summary>
        <pre>{JSON.stringify(fullState, null, 2)}</pre>
      </details>
    </div>
  );
};

export default ReduxDebugger;