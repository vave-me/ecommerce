/**
 * Unified User Avatar Component
 * Displays user profile images with fallback support
 */
import React, { memo } from 'react';
import Image from 'next/image';
import styles from './UserAvatar.module.css';

const UserAvatar = memo(({ 
  src, 
  alt = 'User Avatar',
  width = 40, 
  height = 40,
  size = 'medium',
  className = '',
  priority = false,
  onClick,
  ...props 
}) => {
  // Handle size presets
  const dimensions = {
    small: { width: 32, height: 32 },
    medium: { width: 40, height: 40 },
    large: { width: 56, height: 56 },
    xlarge: { width: 80, height: 80 }
  };
  
  const { width: finalWidth, height: finalHeight } = dimensions[size] || { width, height };
  
  const wrapperClass = `
    ${styles.avatarWrapper} 
    ${styles[size]} 
    ${onClick ? styles.clickable : ''} 
    ${className}
  `.trim();
  
  return (
    <div 
      className={wrapperClass} 
      onClick={onClick}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
      {...props}
    >
      <Image 
        src={src || '/images/default-avatar.png'} 
        alt={alt}
        width={finalWidth} 
        height={finalHeight} 
        className={styles.avatar}
        priority={priority}
        sizes={`${finalWidth}px`}
        onError={(e) => {
          e.currentTarget.src = '/images/default-avatar.png';
        }}
      />
    </div>
  );
});

UserAvatar.displayName = 'UserAvatar';

export default UserAvatar;