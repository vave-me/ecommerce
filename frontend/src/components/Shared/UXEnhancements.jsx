/**
 * UNIFIED UX/UI ENHANCEMENT SYSTEM
 * Modern interaction patterns and accessibility improvements
 * - Touch-friendly interactions
 * - Improved loading states  
 * - Better error handling
 * - Enhanced accessibility
 * - Performance-optimized animations
 */
import React, { useState, useEffect, useCallback, useRef, memo } from 'react';
import styles from './UXEnhancements.module.css';
/**
 * Enhanced Loading Component with modern patterns
 */
export const EnhancedLoading = memo(({ 
  size = 'medium', 
  variant = 'spinner',
  text = '',
  fullscreen = false,
  overlay = false 
}) => {
  const sizeClasses = {
    small: styles.loadingSmall,
    medium: styles.loadingMedium,
    large: styles.loadingLarge
  };
  const variantClasses = {
    spinner: styles.spinnerVariant,
    dots: styles.dotsVariant,
    pulse: styles.pulseVariant,
    skeleton: styles.skeletonVariant
  };
  const containerClass = [
    styles.loadingContainer,
    sizeClasses[size],
    variantClasses[variant],
    fullscreen && styles.fullscreen,
    overlay && styles.overlay
  ].filter(Boolean).join(' ');
  return (
    <div className={containerClass} role="status" aria-label={text || "Loading"}>
      {variant === 'spinner' && (
        <div className={styles.spinner} aria-hidden="true">
          <svg viewBox="0 0 24 24" className={styles.spinnerSvg}>
            <circle cx="12" cy="12" r="10" className={styles.spinnerCircle} />
          </svg>
        </div>
      )}
      {variant === 'dots' && (
        <div className={styles.dots} aria-hidden="true">
          <div className={styles.dot}></div>
          <div className={styles.dot}></div>
          <div className={styles.dot}></div>
        </div>
      )}
      {variant === 'pulse' && (
        <div className={styles.pulse} aria-hidden="true"></div>
      )}
      {variant === 'skeleton' && (
        <div className={styles.skeleton} aria-hidden="true">
          <div className={styles.skeletonLine}></div>
          <div className={styles.skeletonLine}></div>
          <div className={styles.skeletonLine}></div>
        </div>
      )}
      {text && (
        <span className={styles.loadingText}>{text}</span>
      )}
    </div>
  );
});
/**
 * Enhanced Button with modern UX patterns
 */
export const EnhancedButton = memo(({
  children,
  variant = 'primary',
  size = 'medium',
  loading = false,
  disabled = false,
  icon = null,
  iconPosition = 'left',
  ripple = true,
  onClick,
  className = '',
  ...props
}) => {
  const [isPressed, setIsPressed] = useState(false);
  const [ripplePosition, setRipplePosition] = useState({ x: 0, y: 0 });
  const [showRipple, setShowRipple] = useState(false);
  const buttonRef = useRef(null);
  const handleClick = useCallback((e) => {
    if (loading || disabled) return;
    // Ripple effect
    if (ripple && buttonRef.current) {
      const rect = buttonRef.current.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;
      setRipplePosition({ x, y });
      setShowRipple(true);
      setTimeout(() => setShowRipple(false), 600);
    }
    onClick?.(e);
  }, [loading, disabled, ripple, onClick]);
  const buttonClasses = [
    styles.enhancedButton,
    styles[`variant${variant.charAt(0).toUpperCase() + variant.slice(1)}`],
    styles[`size${size.charAt(0).toUpperCase() + size.slice(1)}`],
    loading && styles.loading,
    disabled && styles.disabled,
    isPressed && styles.pressed,
    className
  ].filter(Boolean).join(' ');
  return (
    <button
      ref={buttonRef}
      className={buttonClasses}
      disabled={disabled || loading}
      onClick={handleClick}
      onMouseDown={() => setIsPressed(true)}
      onMouseUp={() => setIsPressed(false)}
      onMouseLeave={() => setIsPressed(false)}
      aria-busy={loading}
      {...props}
    >
      {showRipple && (
        <span 
          className={styles.ripple}
          style={{
            left: ripplePosition.x,
            top: ripplePosition.y
          }}
        />
      )}
      {loading && (
        <EnhancedLoading size="small" variant="spinner" />
      )}
      {!loading && (
        <>
          {icon && iconPosition === 'left' && (
            <span className={styles.iconLeft}>{icon}</span>
          )}
          <span className={styles.buttonText}>{children}</span>
          {icon && iconPosition === 'right' && (
            <span className={styles.iconRight}>{icon}</span>
          )}
        </>
      )}
    </button>
  );
});
/**
 * Enhanced Input with modern UX patterns
 */
export const EnhancedInput = memo(({
  label,
  error,
  helper,
  icon,
  iconPosition = 'left',
  loading = false,
  success = false,
  variant = 'default',
  size = 'medium',
  className = '',
  ...props
}) => {
  const [isFocused, setIsFocused] = useState(false);
  const [hasValue, setHasValue] = useState(props.value || props.defaultValue || '');
  const inputClasses = [
    styles.enhancedInput,
    styles[`inputVariant${variant.charAt(0).toUpperCase() + variant.slice(1)}`],
    styles[`inputSize${size.charAt(0).toUpperCase() + size.slice(1)}`],
    isFocused && styles.focused,
    error && styles.error,
    success && styles.success,
    hasValue && styles.hasValue,
    icon && styles.hasIcon,
    loading && styles.loading,
    className
  ].filter(Boolean).join(' ');
  const handleChange = useCallback((e) => {
    setHasValue(e.target.value);
    props.onChange?.(e);
  }, [props.onChange]);
  return (
    <div className={styles.inputWrapper}>
      {label && (
        <label className={styles.inputLabel}>
          {label}
          {props.required && <span className={styles.required}>*</span>}
        </label>
      )}
      <div className={styles.inputContainer}>
        {icon && iconPosition === 'left' && (
          <span className={styles.inputIconLeft}>{icon}</span>
        )}
        <input
          {...props}
          className={inputClasses}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          onChange={handleChange}
          aria-invalid={!!error}
          aria-describedby={error ? `${props.id}-error` : helper ? `${props.id}-helper` : undefined}
        />
        {loading && (
          <span className={styles.inputLoader}>
            <EnhancedLoading size="small" variant="spinner" />
          </span>
        )}
        {!loading && success && (
          <span className={styles.inputSuccess}>✓</span>
        )}
        {!loading && error && (
          <span className={styles.inputError}>!</span>
        )}
        {icon && iconPosition === 'right' && (
          <span className={styles.inputIconRight}>{icon}</span>
        )}
      </div>
      {error && (
        <div id={`${props.id}-error`} className={styles.errorMessage} role="alert">
          {error}
        </div>
      )}
      {!error && helper && (
        <div id={`${props.id}-helper`} className={styles.helperMessage}>
          {helper}
        </div>
      )}
    </div>
  );
});
/**
 * Enhanced Card with modern interaction patterns
 */
export const EnhancedCard = memo(({
  children,
  hoverable = false,
  clickable = false,
  loading = false,
  variant = 'default',
  padding = 'medium',
  onClick,
  className = '',
  ...props
}) => {
  const [isHovered, setIsHovered] = useState(false);
  const [isPressed, setIsPressed] = useState(false);
  const cardClasses = [
    styles.enhancedCard,
    styles[`cardVariant${variant.charAt(0).toUpperCase() + variant.slice(1)}`],
    styles[`cardPadding${padding.charAt(0).toUpperCase() + padding.slice(1)}`],
    hoverable && styles.hoverable,
    clickable && styles.clickable,
    loading && styles.cardLoading,
    isHovered && styles.hovered,
    isPressed && styles.cardPressed,
    className
  ].filter(Boolean).join(' ');
  const handleClick = useCallback((e) => {
    if (loading) return;
    onClick?.(e);
  }, [loading, onClick]);
  return (
    <div
      className={cardClasses}
      onClick={clickable ? handleClick : undefined}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => {
        setIsHovered(false);
        setIsPressed(false);
      }}
      onMouseDown={() => clickable && setIsPressed(true)}
      onMouseUp={() => setIsPressed(false)}
      role={clickable ? 'button' : undefined}
      tabIndex={clickable ? 0 : undefined}
      {...props}
    >
      {loading && (
        <div className={styles.cardLoadingOverlay}>
          <EnhancedLoading size="medium" variant="spinner" />
        </div>
      )}
      {children}
    </div>
  );
});
/**
 * Enhanced Toast Notification System
 */
export const EnhancedToast = memo(({
  message,
  type = 'info',
  duration = 4000,
  action = null,
  onClose,
  position = 'bottom-right'
}) => {
  const [isVisible, setIsVisible] = useState(true);
  const [isExiting, setIsExiting] = useState(false);
  useEffect(() => {
    if (duration > 0) {
      const timer = setTimeout(() => {
        setIsExiting(true);
        setTimeout(() => {
          setIsVisible(false);
          onClose?.();
        }, 300);
      }, duration);
      return () => clearTimeout(timer);
    }
  }, [duration, onClose]);
  if (!isVisible) return null;
  const toastClasses = [
    styles.enhancedToast,
    styles[`toast${type.charAt(0).toUpperCase() + type.slice(1)}`],
    styles[`toastPosition${position.split('-').map(p => 
      p.charAt(0).toUpperCase() + p.slice(1)
    ).join('')}`],
    isExiting && styles.toastExiting
  ].filter(Boolean).join(' ');
  return (
    <div className={toastClasses} role="alert">
      <div className={styles.toastContent}>
        <span className={styles.toastMessage}>{message}</span>
        {action && (
          <button className={styles.toastAction} onClick={action.onClick}>
            {action.label}
          </button>
        )}
        <button 
          className={styles.toastClose}
          onClick={() => {
            setIsExiting(true);
            setTimeout(() => {
              setIsVisible(false);
              onClose?.();
            }, 300);
          }}
          aria-label="Close notification"
        >
          ×
        </button>
      </div>
    </div>
  );
});
/**
 * Enhanced Modal with modern UX patterns
 */
export const EnhancedModal = memo(({
  isOpen,
  onClose,
  title,
  children,
  size = 'medium',
  closeOnOverlay = true,
  closeOnEscape = true,
  className = ''
}) => {
  const modalRef = useRef(null);
  const [isAnimating, setIsAnimating] = useState(false);
  useEffect(() => {
    if (isOpen) {
      setIsAnimating(true);
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);
  useEffect(() => {
    const handleEscape = (e) => {
      if (closeOnEscape && e.key === 'Escape') {
        onClose?.();
      }
    };
    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      return () => document.removeEventListener('keydown', handleEscape);
    }
  }, [isOpen, closeOnEscape, onClose]);
  const handleOverlayClick = useCallback((e) => {
    if (closeOnOverlay && e.target === e.currentTarget) {
      onClose?.();
    }
  }, [closeOnOverlay, onClose]);
  if (!isOpen && !isAnimating) return null;
  const modalClasses = [
    styles.enhancedModal,
    styles[`modalSize${size.charAt(0).toUpperCase() + size.slice(1)}`],
    isOpen && styles.modalOpen,
    className
  ].filter(Boolean).join(' ');
  return (
    <div 
      className={styles.modalOverlay}
      onClick={handleOverlayClick}
      onAnimationEnd={() => !isOpen && setIsAnimating(false)}
    >
      <div 
        ref={modalRef}
        className={modalClasses}
        role="dialog"
        aria-modal="true"
        aria-labelledby={title ? 'modal-title' : undefined}
      >
        {title && (
          <div className={styles.modalHeader}>
            <h2 id="modal-title" className={styles.modalTitle}>{title}</h2>
            <button 
              className={styles.modalCloseButton}
              onClick={onClose}
              aria-label="Close modal"
            >
              ×
            </button>
          </div>
        )}
        <div className={styles.modalContent}>
          {children}
        </div>
      </div>
    </div>
  );
});
// Export all components
export default {
  EnhancedLoading,
  EnhancedButton,
  EnhancedInput,
  EnhancedCard,
  EnhancedToast,
  EnhancedModal
}; 