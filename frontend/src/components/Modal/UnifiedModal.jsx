"use client";
import React, { useEffect, useRef, useCallback } from 'react';
import { createPortal } from 'react-dom';
import { X } from 'lucide-react';
import PropTypes from 'prop-types';
import styles from './UnifiedModal.module.css';

/**
 * UnifiedModal - WCAG 2.1 Compliant Modal Component
 * 
 * Features:
 * - Keyboard navigation (Escape to close)
 * - Focus trapping
 * - Screen reader announcements
 * - Proper ARIA attributes
 * - Responsive design
 * - Dark mode support via CSS variables
 */
const UnifiedModal = ({
  isOpen,
  onClose,
  title,
  children,
  footer,
  size = 'medium',
  closeOnOverlayClick = true,
  closeOnEscape = true,
  showCloseButton = true,
  className = '',
  loading = false,
  ariaLabel,
  ariaDescribedBy,
}) => {
  const modalRef = useRef(null);
  const previousActiveElement = useRef(null);
  const firstFocusableElement = useRef(null);
  const lastFocusableElement = useRef(null);

  // Handle escape key
  const handleEscape = useCallback((event) => {
    if (closeOnEscape && event.key === 'Escape') {
      onClose();
    }
  }, [closeOnEscape, onClose]);

  // Handle overlay click
  const handleOverlayClick = useCallback((event) => {
    if (closeOnOverlayClick && event.target === event.currentTarget) {
      onClose();
    }
  }, [closeOnOverlayClick, onClose]);

  // Focus trap
  const handleTabKey = useCallback((event) => {
    if (event.key !== 'Tab') return;

    const focusableElements = modalRef.current?.querySelectorAll(
      'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])'
    );

    if (!focusableElements || focusableElements.length === 0) return;

    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];

    if (event.shiftKey && document.activeElement === firstElement) {
      event.preventDefault();
      lastElement.focus();
    } else if (!event.shiftKey && document.activeElement === lastElement) {
      event.preventDefault();
      firstElement.focus();
    }
  }, []);

  // Setup focus management
  useEffect(() => {
    if (isOpen) {
      // Store current focus
      previousActiveElement.current = document.activeElement;

      // Focus first focusable element or modal
      setTimeout(() => {
        const focusableElements = modalRef.current?.querySelectorAll(
          'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])'
        );

        if (focusableElements && focusableElements.length > 0) {
          firstFocusableElement.current = focusableElements[0];
          lastFocusableElement.current = focusableElements[focusableElements.length - 1];
          focusableElements[0].focus();
        } else {
          modalRef.current?.focus();
        }
      }, 100);

      // Add event listeners
      document.addEventListener('keydown', handleEscape);
      document.addEventListener('keydown', handleTabKey);

      // Prevent body scroll
      document.body.style.overflow = 'hidden';

      return () => {
        document.removeEventListener('keydown', handleEscape);
        document.removeEventListener('keydown', handleTabKey);
        document.body.style.overflow = '';

        // Restore focus
        if (previousActiveElement.current) {
          previousActiveElement.current.focus();
        }
      };
    }
  }, [isOpen, handleEscape, handleTabKey]);

  if (!isOpen) return null;

  // Determine size class
  const sizeClass = {
    small: styles.modalSmall,
    medium: styles.modalMedium,
    large: styles.modalLarge,
    fullWidth: styles.modalFullWidth,
  }[size] || styles.modalMedium;

  const modalContent = (
    <div 
      className={styles.overlay}
      onClick={handleOverlayClick}
      role="presentation"
    >
      <div
        ref={modalRef}
        className={`${styles.modal} ${sizeClass} ${loading ? styles.loading : ''} ${className}`}
        role="dialog"
        aria-modal="true"
        aria-label={ariaLabel || title}
        aria-describedby={ariaDescribedBy}
        tabIndex={-1}
      >
        {/* Header */}
        {(title || showCloseButton) && (
          <div className={styles.header}>
            {title && <h2 className={styles.title}>{title}</h2>}
            {showCloseButton && (
              <button
                className={styles.closeButton}
                onClick={onClose}
                aria-label="Close modal"
                type="button"
              >
                <X size={20} />
              </button>
            )}
          </div>
        )}

        {/* Body */}
        <div className={styles.body} id={ariaDescribedBy}>
          {children}
        </div>

        {/* Footer */}
        {footer && (
          <div className={styles.footer}>
            {footer}
          </div>
        )}

        {/* Loading overlay */}
        {loading && (
          <div className={styles.loadingOverlay}>
            <div className={styles.spinner} role="status" aria-label="Loading" />
          </div>
        )}
      </div>
    </div>
  );

  // Portal render
  if (typeof window !== 'undefined') {
    return createPortal(modalContent, document.body);
  }

  return null;
};

UnifiedModal.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  title: PropTypes.string,
  children: PropTypes.node.isRequired,
  footer: PropTypes.node,
  size: PropTypes.oneOf(['small', 'medium', 'large', 'fullWidth']),
  closeOnOverlayClick: PropTypes.bool,
  closeOnEscape: PropTypes.bool,
  showCloseButton: PropTypes.bool,
  className: PropTypes.string,
  loading: PropTypes.bool,
  ariaLabel: PropTypes.string,
  ariaDescribedBy: PropTypes.string,
};

// Button component for modal actions
export const ModalButton = ({ 
  variant = 'secondary', 
  children, 
  onClick, 
  disabled = false,
  loading = false,
  ...props 
}) => {
  const variantClass = {
    primary: styles.buttonPrimary,
    secondary: styles.buttonSecondary,
    danger: styles.buttonDanger,
  }[variant] || styles.buttonSecondary;

  return (
    <button
      className={`${styles.button} ${variantClass}`}
      onClick={onClick}
      disabled={disabled || loading}
      {...props}
    >
      {loading && <span className={styles.spinner} />}
      {children}
    </button>
  );
};

ModalButton.propTypes = {
  variant: PropTypes.oneOf(['primary', 'secondary', 'danger']),
  children: PropTypes.node.isRequired,
  onClick: PropTypes.func,
  disabled: PropTypes.bool,
  loading: PropTypes.bool,
};

export default UnifiedModal;