"use client";
import React, { memo } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { Bell, Heart, MessageCircle, ShoppingBag, DoorOpen } from '@/icons';
import { useAuth } from '../../context/AuthContext';
import ModeSwitcher from './ModeSwitcher';
import styles from './Header.module.css';
// This is a new, simplified navigation component for the classic desktop header.
function ClassicNav() {
    const t = useTranslations('Header');
    const pathname = usePathname();
    const router = useRouter();
    const { user } = useAuth();
    // Hardcoded badge counts for demonstration. Replace with actual data from a store or context.
    const badgeCounts = {
        notifications: 1,
        messages: 5,
        wishlist: 0,
        cartItems: 0, // Cart items count - replace with actual cart state
    };
    // Only show cart if user has items in cart
    const hasCartItems = badgeCounts.cartItems > 0;
    // If user is not logged in, show login button
    if (!user) {
        return (
            <div className={styles.desktopNav}>
                <button
                    className={`${styles.navItem} ${styles.loginButton}`}
                    onClick={() => router.push('/login')}
                    aria-label={t('loginButtonAriaLabel')}
                >
                    <DoorOpen size={16} aria-hidden="true" />
                    <span className={styles.loginText}>{t('loginButton')}</span>
                </button>
            </div>
        );
    }
    // If user is logged in, show regular navigation icons
    return (
        <div className={styles.desktopNav}>
            <button
                className={styles.navItem}
                onClick={() => router.push('/notifications')}
                aria-label={t('notificationsButtonAriaLabel', { count: badgeCounts.notifications || 0 })}
            >
                <Bell size={16} aria-hidden="true" />
            </button>
            <button
                className={styles.navItem}
                onClick={() => router.push('/messages')}
                aria-label={t('messagesButtonAriaLabel', { count: badgeCounts.messages || 0 })}
            >
                <MessageCircle size={16} aria-hidden="true" />
            </button>
            <button
                className={styles.navItem}
                onClick={() => router.push('/wishlist')}
                aria-label={t('wishlistButtonAriaLabel', { count: badgeCounts.wishlist || 0 })}
            >
                <Heart size={16} aria-hidden="true" />
            </button>
            {hasCartItems && (
                <button
                    className={styles.navItem}
                    onClick={() => router.push('/cart')}
                    aria-label={t('cartButtonAriaLabel', { count: badgeCounts.cartItems })}
                >
                    <ShoppingBag size={16} aria-hidden="true" />
                </button>
            )}
        </div>
    );
}
export default memo(ClassicNav); 