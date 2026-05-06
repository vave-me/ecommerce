import React from 'react';

export const Alert = ({ children, variant = 'default', className = '' }) => (
  <div className={`alert alert-${variant} ${className}`}>
    {children}
  </div>
);

export const AlertDescription = ({ children, className = '' }) => (
  <div className={`alert-description ${className}`}>
    {children}
  </div>
);