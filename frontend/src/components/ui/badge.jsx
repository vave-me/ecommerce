import React from 'react';

export const Badge = ({ children, variant = 'default', className = '' }) => (
  <span className={`badge badge-${variant} ${className}`}>
    {children}
  </span>
);