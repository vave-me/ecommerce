// src/components/Settings.jsx
"use client"
import React, { useState, memo } from 'react';
import { Settings as SettingsIcon, User, Bell, Shield, CreditCard, HelpCircle } from '@/icons';
import { FaCheckCircle } from '../../utils/iconImports';
import DisplaySettings from './DisplaySettings';
import NotificationSettings from './NotificationSettings';
import PrivacySettings from './PrivacySettings';
import SubscriptionSettings from './SubscriptionSettings';
import styles from './Settings.module.css';
import AccountSettings from './AccountSettings';
import SecuritySettings from "./SecuritySettings";
const Settings = memo(() => {
    const [activeTab, setActiveTab] = useState('Account');
    const renderActiveTab = () => {
        switch (activeTab) {
            case 'Account':
                return <AccountSettings/>;
            case 'Notifications':
                return <NotificationSettings/>;
            case 'Privacy':
                return <PrivacySettings/>;
            case 'Subscriptions':
                return <SubscriptionSettings/>;
            case 'Display':
                return <DisplaySettings/>;
            case 'Security':
                return <SecuritySettings/>;
            default:
                return <AccountSettings/>;
        }
    };
    return (
        <div className={styles.container}>
            <div className={styles.content}>
                <div className={styles.sidebar}>
                    <h2 className={styles.settingsTitle}>Settings</h2>
                    <button
                        className={`${styles.tab} ${activeTab === 'Account' ? styles.active : ''}`}
                        onClick={() => setActiveTab('Account')}
                        aria-label="Account Settings"
                    >
                        <FaCheckCircle className={styles.tabIcon}/>
                        <span className={styles.tabText}>Account</span>
                    </button>
                    <button
                        className={`${styles.tab} ${activeTab === 'Notifications' ? styles.active : ''}`}
                        onClick={() => setActiveTab('Notifications')}
                        aria-label="Notification Settings"
                    >
                        <Bell className={styles.tabIcon}/>
                        <span className={styles.tabText}>Notifications</span>
                    </button>
                    <button
                        className={`${styles.tab} ${activeTab === 'Privacy' ? styles.active : ''}`}
                        onClick={() => setActiveTab('Privacy')}
                        aria-label="Privacy Settings"
                    >
                        <HelpCircle className={styles.tabIcon}/>
                        <span className={styles.tabText}>Privacy</span>
                    </button>
                    <button
                        className={`${styles.tab} ${activeTab === 'Subscriptions' ? styles.active : ''}`}
                        onClick={() => setActiveTab('Subscriptions')}
                        aria-label="Subscription Settings"
                    >
                        <CreditCard className={styles.tabIcon}/>
                        <span className={styles.tabText}>Subscriptions</span>
                    </button>
                    <button
                        className={`${styles.tab} ${activeTab === 'Display' ? styles.active : ''}`}
                        onClick={() => setActiveTab('Display')}
                        aria-label="Display Settings"
                    >
                        <SettingsIcon className={styles.tabIcon}/>
                        <span className={styles.tabText}>Display</span>
                    </button>
                    <button
                        className={`${styles.tab} ${activeTab === 'Security' ? styles.active : ''}`}
                        onClick={() => setActiveTab('Security')}
                        aria-label="Security Settings"
                    >
                        <Shield className={styles.tabIcon}/>
                        <span className={styles.tabText}>Security</span>
                    </button>
                </div>
                <div className={styles.mainContent}>
                    {renderActiveTab()}
                </div>
            </div>
        </div>
    );
});
Settings.displayName = 'Settings';
export default Settings;