"use client";
import React, { useState, useCallback, useRef, useEffect, useMemo, memo } from 'react';
import { useTranslations } from 'next-intl';
import { useIsMobile } from '../../hooks/useMobileDetection';
import { useUserRole } from '../../hooks/useUserRole';
import { roleMenuItems } from '../../config/roleMenuConfig';
import * as Dialog from '@radix-ui/react-dialog';
import { 
  User, 
  List, 
  ShoppingBag, 
  Tag, 
  Headphones, 
  Activity, 
  Star, 
  CreditCard, 
  Mail, 
  Settings, 
  LogOut,
  X
} from '@/icons';
import styles from './SignOutMenu.module.css';
/**
 * SignOutMenu Component - Optimized to prevent multiple re-renders
 * Uses React.memo and memoized menu structure
 */
const SignOutMenu = memo(({ 
  id, 
  className, 
  isOpen = true, 
  onClose, 
  user, 
  handleSignOut, 
  onSignOut, 
  isMobile,
  triggerRef
}) => {
  const t = useTranslations('UserMenu');
  const isMobileDetected = useIsMobile();
  const effectiveMobile = isMobile ?? isMobileDetected;
  const [isClosing, setIsClosing] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ top: 70, right: 16, left: 'auto' });
  const menuRef = useRef(null);
  
  // Get user role
  const { role } = useUserRole();
  
  // Get role-specific menu items
  const roleSpecificItems = useMemo(() => {
    return roleMenuItems[role] || roleMenuItems.customer;
  }, [role]);
  
  // Memoized menu structure to prevent recreation on every render
  const menuSections = useMemo(() => {
    // For admin and business roles, show only role-specific items
    if (role === 'admin' || role === 'business') {
      return [
        {
          id: 'role-specific',
          title: role === 'admin' ? t('section_admin_title', { defaultValue: 'Admin Tools' }) : t('section_business_title', { defaultValue: 'Business Tools' }),
          items: roleSpecificItems
        }
      ];
    }
    
    // For customers, show the original menu structure plus role-specific items
    return [
      {
        id: 'account',
        title: t('section_account_title'),
        items: roleSpecificItems.filter(item => 
          ['profile', 'orders', 'wishlist', 'settings'].includes(item.id)
        )
      },
      {
        id: 'support',
        title: t('section_support_title'),
        items: [
          {
            id: 'support',
            label: t('item_support_label'),
            icon: Headphones,
            href: '/support'
          },
          {
            id: 'activity',
            label: t('item_activity_label'),
            icon: Activity,
            href: '/activity'
          },
          {
            id: 'reviews',
            label: t('item_reviews_label'),
            icon: Star,
            href: '/reviews'
          }
        ]
      },
      {
        id: 'preferences',
        title: t('section_preferences_title'),
        items: [
          {
            id: 'payments',
            label: t('item_payments_label'),
            icon: CreditCard,
            href: '/payments'
          },
          {
            id: 'newsletters',
            label: t('item_newsletters_label'),
            icon: Mail,
            href: '/newsletters'
          }
        ]
      }
    ];
  }, [t, role, roleSpecificItems]); // Re-create when translations or role change
  // Calculate menu position based on trigger button location
  useEffect(() => {
    if (!isOpen || effectiveMobile || !triggerRef?.current) return;
    const calculatePosition = () => {
      const triggerElement = triggerRef.current;
      if (!triggerElement) return;
      const triggerRect = triggerElement.getBoundingClientRect();
      const viewportWidth = window.innerWidth;
      const menuWidth = 280;
      // Calculate top position (below the trigger button)
      const top = triggerRect.bottom + 8;
      // Calculate horizontal position
      let right = 'auto';
      let left = 'auto';
      // Default: align right edge of menu with right edge of button
      const defaultRight = viewportWidth - triggerRect.right;
      // Check if menu would overflow viewport on the right
      const wouldOverflowRight = (viewportWidth - triggerRect.right) < menuWidth && triggerRect.right > menuWidth;
      // Check if menu would overflow viewport on the left  
      const wouldOverflowLeft = triggerRect.left < 16;
      if (wouldOverflowRight && !wouldOverflowLeft) {
        // Position menu to the left of trigger
        right = viewportWidth - triggerRect.left;
      } else if (wouldOverflowLeft && wouldOverflowRight) {
        // Center the menu
        left = Math.max(16, (viewportWidth - menuWidth) / 2);
        right = 'auto';
      } else {
        // Default positioning (right-aligned with trigger)
        right = Math.max(16, defaultRight);
      }
          // Debug logging to understand the positioning issue (development only)
    if (process.env.NODE_ENV === 'development') {
      }
      setMenuPosition({ top, right, left });
    };
    // Initial calculation
    calculatePosition();
    const handleResize = () => calculatePosition();
    window.addEventListener('resize', handleResize, { passive: true });
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [isOpen, effectiveMobile, triggerRef]);
  // Stable callback references
  const handleClose = useCallback(() => {
    if (isClosing) return; // Prevent multiple rapid calls
    setIsClosing(true);
    onClose(); // Call immediately, no setTimeout needed
  }, [onClose, isClosing]);
  const handleSignOutClick = useCallback(() => {
    if (isClosing) return; // Prevent multiple rapid calls
    handleClose();
    setTimeout(() => {
      (handleSignOut || onSignOut)?.();
    }, 250);
  }, [handleClose, handleSignOut, onSignOut, isClosing]);
  const handleMenuItemClick = useCallback((item) => {
    handleClose();
    // Navigate to href
    if (item.href) {
      window.location.href = item.href;
    }
  }, [handleClose]);
  // Optimized keyboard handling with stable reference
  useEffect(() => {
    if (!isOpen) return;
    const handleKeyDown = (e) => {
      if (e.key === 'Escape') {
        handleClose();
      }
    };
    document.addEventListener('keydown', handleKeyDown);
    // Prevent body scroll on mobile
    if (effectiveMobile) {
      document.body.style.overflow = 'hidden';
    }
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.body.style.overflow = '';
    };
  }, [isOpen, handleClose, effectiveMobile]);
  if (!isOpen) return null;
  // Mobile bottom sheet
  if (effectiveMobile) {
    return (
      <Dialog.Root open={isOpen && !isClosing} onOpenChange={(open) => {
        if (!open && !isClosing) handleClose();
      }}>
        <Dialog.Portal>
          <Dialog.Overlay className={`${styles.overlay} ${isClosing ? styles.overlayClosing : ''}`} />
          <Dialog.Content 
            className={`${styles.mobileMenu} ${isClosing ? styles.mobileMenuClosing : ''}`}
            ref={menuRef}
          >
            <Dialog.Title className={styles.visuallyHidden}>
              {t('menuAriaLabel')}
            </Dialog.Title>
            {/* Drag handle */}
            <div className={styles.dragHandle} />
            {/* Close button */}
            <button 
              className={styles.closeButton}
              onClick={(e) => {
                e.stopPropagation();
                handleClose();
              }}
              aria-label="Close menu"
            >
              <X size={20} />
            </button>
            {/* User info */}
            <div className={styles.userSection}>
              <div className={styles.avatarWrapper}>
                {user?.avatar ? (
                  <img 
                    src={user.avatar} 
                    alt={user.name || t('userFallbackName')}
                    className={styles.avatar}
                  />
                ) : (
                  <User size={24} className={styles.userIcon} />
                )}
              </div>
              <div className={styles.userInfo}>
                <div className={styles.userName}>
                  {user?.name || t('userFallbackName')}
                </div>
                {user?.email && (
                  <div className={styles.userEmail}>{user.email}</div>
                )}
              </div>
            </div>
            {/* Menu sections */}
            <div className={styles.menuContent}>
              {menuSections.map((section, sectionIndex) => (
                <div key={section.id} className={styles.menuSection}>
                  <div className={styles.sectionTitle}>{section.title}</div>
                  {section.items.map((item) => (
                    <button
                      key={item.id}
                      className={styles.menuItem}
                      onClick={() => handleMenuItemClick(item)}
                    >
                      <item.icon size={20} className={styles.menuIcon} />
                      <span className={styles.menuLabel}>{item.label}</span>
                    </button>
                  ))}
                </div>
              ))}
              {/* Sign out section */}
              <div className={styles.signOutSection}>
                <button 
                  className={styles.signOutButton}
                  onClick={handleSignOutClick}
                >
                  <LogOut size={20} className={styles.signOutIcon} />
                  <span>{t('item_signOut_label')}</span>
                </button>
              </div>
            </div>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    );
  }
  // Desktop dropdown
  return (
    <Dialog.Root open={isOpen && !isClosing} onOpenChange={(open) => {
      if (!open && !isClosing) handleClose();
    }}>
      <Dialog.Portal>
        <Dialog.Overlay className={styles.desktopOverlay} onClick={(e) => {
          e.stopPropagation();
          handleClose();
        }} />
        <Dialog.Content 
          className={`${styles.desktopMenu} ${isClosing ? styles.desktopMenuClosing : ''}`}
          ref={menuRef}
          style={{
            '--menu-top': `${menuPosition.top}px`,
            '--menu-right': menuPosition.right === 'auto' ? 'auto' : `${menuPosition.right}px`,
            '--menu-left': menuPosition.left === 'auto' ? 'auto' : `${menuPosition.left}px`
          }}
        >
          <Dialog.Title className={styles.visuallyHidden}>
            {t('menuAriaLabel')}
          </Dialog.Title>
          {/* User info */}
          <div className={styles.desktopUserSection}>
            <div className={styles.avatarWrapper}>
              {user?.avatar ? (
                <img 
                  src={user.avatar} 
                  alt={user.name || t('userFallbackName')}
                  className={styles.avatar}
                />
              ) : (
                <User size={20} className={styles.userIcon} />
              )}
            </div>
            <div className={styles.userInfo}>
              <div className={styles.userName}>
                {user?.name || t('userFallbackName')}
              </div>
              {user?.email && (
                <div className={styles.userEmail}>{user.email}</div>
              )}
            </div>
          </div>
          {/* Menu sections */}
          <div className={styles.desktopMenuContent}>
            {menuSections.map((section) => (
              <div key={section.id} className={styles.desktopMenuSection}>
                <div className={styles.desktopSectionTitle}>{section.title}</div>
                {section.items.map((item) => (
                  <button
                    key={item.id}
                    className={styles.desktopMenuItem}
                    onClick={() => handleMenuItemClick(item)}
                  >
                    <item.icon size={16} className={styles.desktopMenuIcon} />
                    <span>{item.label}</span>
                  </button>
                ))}
              </div>
            ))}
            {/* Sign out */}
            <div className={styles.desktopSignOutSection}>
              <button 
                className={styles.desktopSignOutButton}
                onClick={handleSignOutClick}
              >
                <LogOut size={16} className={styles.signOutIcon} />
                <span>{t('item_signOut_label')}</span>
              </button>
            </div>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}, (prevProps, nextProps) => {
  // Custom comparison for React.memo to prevent unnecessary re-renders
  return (
    prevProps.isOpen === nextProps.isOpen &&
    prevProps.isMobile === nextProps.isMobile &&
    prevProps.user?.userId === nextProps.user?.userId &&
    prevProps.user?.name === nextProps.user?.name &&
    prevProps.user?.email === nextProps.user?.email &&
    prevProps.user?.avatar === nextProps.user?.avatar &&
    prevProps.onClose === nextProps.onClose &&
    prevProps.handleSignOut === nextProps.handleSignOut &&
    prevProps.onSignOut === nextProps.onSignOut
  );
});
SignOutMenu.displayName = 'SignOutMenu';
export default SignOutMenu; 