'use client';
import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { showUnifiedComposer, selectShowUnifiedComposer } from '../../redux/slices/uiPreferencesSlice';
import { Edit3, Plus } from '@/icons';
import styles from './ShowComposerButton.module.css';

/**
 * Small floating button to show the composer when it's hidden
 * Can be placed anywhere in the UI
 */
const ShowComposerButton = ({ 
  variant = 'floating', // 'floating', 'inline', 'minimal'
  position = 'bottom-right', // for floating variant
  text = 'Create',
  icon = 'edit'
}) => {
  const dispatch = useDispatch();
  const isComposerVisible = useSelector(selectShowUnifiedComposer);
  const [forceUpdate, setForceUpdate] = useState(0);

  // Force re-render when Redux state changes
  useEffect(() => {
    
  }, [isComposerVisible]);

  // Don't show if composer is already visible
  if (isComposerVisible) {
    return null;
  }

  const handleClick = () => {

    try {
      const action = showUnifiedComposer();
      
      dispatch(action);

      // Force a re-check after a small delay
      setTimeout(() => {
        const state = dispatch((state) => state);

      }, 100);
    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
  };

  const Icon = icon === 'plus' ? Plus : Edit3;

  // Floating variant
  if (variant === 'floating') {
    return (
      <button
        onClick={handleClick}
        className={`${styles.floatingButton} ${styles[position]}`}
        title="Show composer"
        aria-label="Show composer"
      >
        <Icon size={24} />
        <span className={styles.floatingText}>{text}</span>
      </button>
    );
  }

  // Inline variant
  if (variant === 'inline') {
    return (
      <button
        onClick={handleClick}
        className={styles.inlineButton}
        title="Show composer"
      >
        <Icon size={20} />
        <span>{text}</span>
      </button>
    );
  }

  // Minimal variant
  return (
    <button
      onClick={handleClick}
      className={styles.minimalButton}
      title="Show composer"
      aria-label="Show composer"
    >
      <Icon size={18} />
    </button>
  );
};

export default ShowComposerButton;