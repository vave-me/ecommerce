// src/app/[locale]/settings/page.jsx
"use client"
import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useAuth } from '../../../context/AuthContext';
import {
    UserCircle,
    Bell,
    Shield,
    Eye,
    Palette,
    CreditCard,
    Settings as SettingsIcon,
    ChevronRight,
    Menu,
    X,
    RefreshCw
} from '@/icons';
import styles from './Settings.module.css';
import AccountSettings from "../../../components/Settings/AccountSettings";
import NotificationSettings from "../../../components/Settings/NotificationSettings";
import PrivacySettings from "../../../components/Settings/PrivacySettings";
import SubscriptionSettings from "../../../components/Settings/SubscriptionSettings";
import DisplaySettings from "../../../components/Settings/DisplaySettings";
import SecuritySettings from "../../../components/Settings/SecuritySettings";
/**
 * Modern Settings Page Component
 * Elegant design matching notifications and wishlists pages
 */
const Settings = () => {
    const [activeTab, setActiveTab] = useState('Account');
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
    const { user, isLoading: authLoading } = useAuth();
    const t = useTranslations('Settings');
    const [loading, setLoading] = useState(false);
    const [refreshing, setRefreshing] = useState(false);
    // Redirect to login if user is not authenticated
    useEffect(() => {
        if (!authLoading && !user && typeof window !== 'undefined') {
            window.location.href = '/login';
        }
    }, [user, authLoading]);
    // Close mobile menu when tab changes
    useEffect(() => {
        setIsMobileMenuOpen(false);
    }, [activeTab]);
    // Enhanced tabs configuration with descriptions
    const TABS = [
        { 
            id: 'Account', 
            icon: UserCircle, 
            label: t('tabs.account'),
            description: t('tabs.accountDesc')
        },
        { 
            id: 'Notifications', 
            icon: Bell, 
            label: t('tabs.notifications'),
            description: t('tabs.notificationsDesc')
        },
        { 
            id: 'Privacy', 
            icon: Eye, 
            label: t('tabs.privacy'),
            description: t('tabs.privacyDesc')
        },
        { 
            id: 'Security', 
            icon: Shield, 
            label: t('tabs.security'),
            description: t('tabs.securityDesc')
        },
        { 
            id: 'Display', 
            icon: Palette, 
            label: t('tabs.display'),
            description: t('tabs.displayDesc')
        },
        { 
            id: 'Subscriptions', 
            icon: CreditCard, 
            label: t('tabs.subscriptions'),
            description: t('tabs.subscriptionsDesc')
        },
    ];
    const handleRefresh = async () => {
        setRefreshing(true);
        // Simulate refresh - in real implementation, this would reload settings data
        setTimeout(() => {
            setRefreshing(false);
        }, 1000);
    };
    const renderActiveTab = () => {
        if (loading) {
            return (
                <div className={styles.loadingState}>
                    <div className={styles.spinner}></div>
                    <h3>{t('loading')}</h3>
                    <p>{t('loadingDesc')}</p>
                </div>
            );
        }
        switch (activeTab) {
            case 'Account':
                return <AccountSettings userId={user?.userId} />;
            case 'Notifications':
                return <NotificationSettings userId={user?.userId} />;
            case 'Privacy':
                return <PrivacySettings userId={user?.userId} />;
            case 'Security':
                return <SecuritySettings userId={user?.userId} />;
            case 'Display':
                return <DisplaySettings userId={user?.userId} />;
            case 'Subscriptions':
                return <SubscriptionSettings userId={user?.userId} />;
            default:
                return <AccountSettings userId={user?.userId} />;
        }
    };
    // Show loading if user is not loaded yet or still loading
    if (authLoading || !user) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingContainer}>
                    <div className={styles.spinner}></div>
                    <h3>{t('authenticating')}</h3>
                    <p>{t('authenticatingDesc')}</p>
                </div>
            </div>
        );
    }
    const currentTab = TABS.find(tab => tab.id === activeTab);
    return (
        <div className={styles.container}>
            {/* Elegant Header */}
            <div className={styles.header}>
                <div className={styles.headerContent}>
                    <div className={styles.titleSection}>
                        <SettingsIcon size={28} className={styles.titleIcon} />
                        <div>
                            <h1 className={styles.pageTitle}>{t('pageTitle')}</h1>
                            <p className={styles.pageSubtitle}>
                                {currentTab?.description || t('pageDesc')}
                            </p>
                        </div>
                    </div>
                    <div className={styles.headerActions}>
                        <button
                            onClick={handleRefresh}
                            disabled={refreshing}
                            className={styles.refreshButton}
                            title={t('refresh')}
                        >
                            <RefreshCw 
                                size={16} 
                                className={refreshing ? styles.spinning : ''} 
                            />
                            <span className={styles.buttonText}>{t('refresh')}</span>
                        </button>
                        {/* Mobile menu toggle */}
                        <button
                            className={styles.mobileMenuToggle}
                            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                            aria-label={t('toggleMenu')}
                        >
                            {isMobileMenuOpen ? <X size={20} /> : <Menu size={20} />}
                        </button>
                    </div>
                </div>
                {/* Breadcrumb */}
                <div className={styles.breadcrumb}>
                    <span className={styles.breadcrumbItem}>{t('pageTitle')}</span>
                    <ChevronRight size={14} className={styles.breadcrumbSeparator} />
                    <span className={styles.breadcrumbCurrent}>{currentTab?.label}</span>
                </div>
            </div>
            <div className={styles.content}>
                {/* Elegant Sidebar */}
                <div className={`${styles.sidebar} ${isMobileMenuOpen ? styles.sidebarOpen : ''}`}>
                    <nav className={styles.navigation} role="tablist">
                        {TABS.map((tab) => (
                            <button
                                key={tab.id}
                                className={`${styles.navItem} ${activeTab === tab.id ? styles.active : ''}`}
                                onClick={() => setActiveTab(tab.id)}
                                aria-label={t('tabAriaLabel', { tabName: tab.label })}
                                aria-selected={activeTab === tab.id}
                                role="tab"
                            >
                                <div className={styles.navItemContent}>
                                    <tab.icon size={20} className={styles.navIcon} />
                                    <div className={styles.navText}>
                                        <span className={styles.navLabel}>{tab.label}</span>
                                        <span className={styles.navDescription}>{tab.description}</span>
                                    </div>
                                </div>
                                {activeTab === tab.id && (
                                    <div className={styles.activeIndicator} />
                                )}
                                <ChevronRight size={16} className={styles.navChevron} />
                            </button>
                        ))}
                    </nav>
                </div>
                {/* Main Content Area */}
                <main className={styles.mainContent}>
                    <div className={styles.contentWrapper}>
                        {renderActiveTab()}
                    </div>
                </main>
            </div>
            {/* Mobile overlay */}
            {isMobileMenuOpen && (
                <div 
                    className={styles.mobileOverlay}
                    onClick={() => setIsMobileMenuOpen(false)}
                />
            )}
        </div>
    );
};
export default Settings;