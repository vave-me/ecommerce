import React from 'react';

export const Button = ({ 
  children, 
  variant = 'default', 
  onClick, 
  disabled = false,
  className = '' 
}) => (
  <button 
    className={`button button-${variant} ${className}`}
    onClick={onClick}
    disabled={disabled}
  >
    {children}
  </button>
);