"use client";

import React, { useState, useCallback } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { 
  Settings,
  Globe,
  Mail,
  Shield,
  Bell,
  Database,
  Zap,
  Users,
  CreditCard,
  Save,
  RefreshCw
} from 'lucide-react';
import { toast } from 'react-toastify';
import styles from './AdminSettings.module.css';

const AdminSettings = () => {
  const t = useTranslations('AdminSettings');
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState('general');
  const [hasChanges, setHasChanges] = useState(false);

  // Mock settings data - replace with API
  const [settings, setSettings] = useState({
    general: {
      siteName: 'Classified Platform',
      siteUrl: 'https://example.com',
      supportEmail: 'redacted-email@example.com',
      defaultLanguage: 'en',
      timeZone: 'UTC',
      maintenanceMode: false
    },
    features: {
      userRegistration: true,
      emailVerification: true,
      socialLogin: true,
      guestCheckout: false,
      multiCurrency: true,
      aiRecommendations: true
    },
    email: {
      smtpHost: 'smtp.example.com',
      smtpPort: '587',
      smtpUser: 'redacted-email@example.com',
      smtpSecure: true,
      emailFrom: 'Classified Platform <redacted-email@example.com>'
    },
    security: {
      twoFactorAuth: true,
      passwordMinLength: 8,
      sessionTimeout: 3600,
      maxLoginAttempts: 5,
      ipWhitelist: [],
      requireCaptcha: true
    },
    notifications: {
      newUserAlert: true,
      lowStockAlert: true,
      paymentFailureAlert: true,
      systemErrorAlert: true,
      dailyReportEmail: true
    },
    payments: {
      stripeEnabled: true,
      paypalEnabled: false,
      bankTransferEnabled: true,
      transactionFee: 2.5,
      minOrderAmount: 10
    }
  });

  const saveMutation = useMutation({
    mutationFn: async (updatedSettings) => {
      // TODO: Replace with actual API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      return { success: true };
    },
    onSuccess: () => {
      toast.success(t('saveSuccess'));
      setHasChanges(false);
      queryClient.invalidateQueries(['admin-settings']);
    },
    onError: () => {
      toast.error(t('saveError'));
    }
  });

  const handleSettingChange = useCallback((category, field, value) => {
    setSettings(prev => ({
      ...prev,
      [category]: {
        ...prev[category],
        [field]: value
      }
    }));
    setHasChanges(true);
  }, []);

  const handleSave = useCallback(() => {
    saveMutation.mutate(settings);
  }, [settings, saveMutation]);

  const tabs = [
    { id: 'general', label: t('tabs.general'), icon: Settings },
    { id: 'features', label: t('tabs.features'), icon: Zap },
    { id: 'email', label: t('tabs.email'), icon: Mail },
    { id: 'security', label: t('tabs.security'), icon: Shield },
    { id: 'notifications', label: t('tabs.notifications'), icon: Bell },
    { id: 'payments', label: t('tabs.payments'), icon: CreditCard }
  ];

  const renderGeneralSettings = () => (
    <div className={styles.settingsGroup}>
      <h3>{t('general.title')}</h3>
      <div className={styles.field}>
        <label>{t('general.siteName')}</label>
        <input
          type="text"
          value={settings.general.siteName}
          onChange={(e) => handleSettingChange('general', 'siteName', e.target.value)}
        />
      </div>
      <div className={styles.field}>
        <label>{t('general.siteUrl')}</label>
        <input
          type="url"
          value={settings.general.siteUrl}
          onChange={(e) => handleSettingChange('general', 'siteUrl', e.target.value)}
        />
      </div>
      <div className={styles.field}>
        <label>{t('general.supportEmail')}</label>
        <input
          type="email"
          value={settings.general.supportEmail}
          onChange={(e) => handleSettingChange('general', 'supportEmail', e.target.value)}
        />
      </div>
      <div className={styles.field}>
        <label>{t('general.defaultLanguage')}</label>
        <select
          value={settings.general.defaultLanguage}
          onChange={(e) => handleSettingChange('general', 'defaultLanguage', e.target.value)}
        >
          <option value="en">English</option>
          <option value="de">Deutsch</option>
          <option value="it">Italiano</option>
          <option value="pl">Polski</option>
        </select>
      </div>
      <div className={styles.field}>
        <label className={styles.toggleLabel}>
          <input
            type="checkbox"
            checked={settings.general.maintenanceMode}
            onChange={(e) => handleSettingChange('general', 'maintenanceMode', e.target.checked)}
          />
          <span>{t('general.maintenanceMode')}</span>
        </label>
      </div>
    </div>
  );

  const renderFeaturesSettings = () => (
    <div className={styles.settingsGroup}>
      <h3>{t('features.title')}</h3>
      {Object.entries(settings.features).map(([key, value]) => (
        <div key={key} className={styles.field}>
          <label className={styles.toggleLabel}>
            <input
              type="checkbox"
              checked={value}
              onChange={(e) => handleSettingChange('features', key, e.target.checked)}
            />
            <span>{t(`features.${key}`)}</span>
          </label>
        </div>
      ))}
    </div>
  );

  const renderEmailSettings = () => (
    <div className={styles.settingsGroup}>
      <h3>{t('email.title')}</h3>
      <div className={styles.field}>
        <label>{t('email.smtpHost')}</label>
        <input
          type="text"
          value={settings.email.smtpHost}
          onChange={(e) => handleSettingChange('email', 'smtpHost', e.target.value)}
        />
      </div>
      <div className={styles.field}>
        <label>{t('email.smtpPort')}</label>
        <input
          type="number"
          value={settings.email.smtpPort}
          onChange={(e) => handleSettingChange('email', 'smtpPort', e.target.value)}
        />
      </div>
      <div className={styles.field}>
        <label>{t('email.smtpUser')}</label>
        <input
          type="text"
          value={settings.email.smtpUser}
          onChange={(e) => handleSettingChange('email', 'smtpUser', e.target.value)}
        />
      </div>
      <div className={styles.field}>
        <label className={styles.toggleLabel}>
          <input
            type="checkbox"
            checked={settings.email.smtpSecure}
            onChange={(e) => handleSettingChange('email', 'smtpSecure', e.target.checked)}
          />
          <span>{t('email.smtpSecure')}</span>
        </label>
      </div>
    </div>
  );

  const renderSecuritySettings = () => (
    <div className={styles.settingsGroup}>
      <h3>{t('security.title')}</h3>
      <div className={styles.field}>
        <label className={styles.toggleLabel}>
          <input
            type="checkbox"
            checked={settings.security.twoFactorAuth}
            onChange={(e) => handleSettingChange('security', 'twoFactorAuth', e.target.checked)}
          />
          <span>{t('security.twoFactorAuth')}</span>
        </label>
      </div>
      <div className={styles.field}>
        <label>{t('security.passwordMinLength')}</label>
        <input
          type="number"
          value={settings.security.passwordMinLength}
          onChange={(e) => handleSettingChange('security', 'passwordMinLength', parseInt(e.target.value))}
          min="6"
          max="32"
        />
      </div>
      <div className={styles.field}>
        <label>{t('security.sessionTimeout')}</label>
        <input
          type="number"
          value={settings.security.sessionTimeout}
          onChange={(e) => handleSettingChange('security', 'sessionTimeout', parseInt(e.target.value))}
          min="300"
          step="300"
        />
        <span className={styles.fieldHint}>{t('security.sessionTimeoutHint')}</span>
      </div>
      <div className={styles.field}>
        <label>{t('security.maxLoginAttempts')}</label>
        <input
          type="number"
          value={settings.security.maxLoginAttempts}
          onChange={(e) => handleSettingChange('security', 'maxLoginAttempts', parseInt(e.target.value))}
          min="3"
          max="10"
        />
      </div>
    </div>
  );

  const renderNotificationSettings = () => (
    <div className={styles.settingsGroup}>
      <h3>{t('notifications.title')}</h3>
      {Object.entries(settings.notifications).map(([key, value]) => (
        <div key={key} className={styles.field}>
          <label className={styles.toggleLabel}>
            <input
              type="checkbox"
              checked={value}
              onChange={(e) => handleSettingChange('notifications', key, e.target.checked)}
            />
            <span>{t(`notifications.${key}`)}</span>
          </label>
        </div>
      ))}
    </div>
  );

  const renderPaymentSettings = () => (
    <div className={styles.settingsGroup}>
      <h3>{t('payments.title')}</h3>
      <div className={styles.field}>
        <label className={styles.toggleLabel}>
          <input
            type="checkbox"
            checked={settings.payments.stripeEnabled}
            onChange={(e) => handleSettingChange('payments', 'stripeEnabled', e.target.checked)}
          />
          <span>{t('payments.stripeEnabled')}</span>
        </label>
      </div>
      <div className={styles.field}>
        <label className={styles.toggleLabel}>
          <input
            type="checkbox"
            checked={settings.payments.paypalEnabled}
            onChange={(e) => handleSettingChange('payments', 'paypalEnabled', e.target.checked)}
          />
          <span>{t('payments.paypalEnabled')}</span>
        </label>
      </div>
      <div className={styles.field}>
        <label>{t('payments.transactionFee')}</label>
        <input
          type="number"
          value={settings.payments.transactionFee}
          onChange={(e) => handleSettingChange('payments', 'transactionFee', parseFloat(e.target.value))}
          min="0"
          max="10"
          step="0.1"
        />
        <span className={styles.fieldHint}>%</span>
      </div>
      <div className={styles.field}>
        <label>{t('payments.minOrderAmount')}</label>
        <input
          type="number"
          value={settings.payments.minOrderAmount}
          onChange={(e) => handleSettingChange('payments', 'minOrderAmount', parseFloat(e.target.value))}
          min="0"
          step="1"
        />
      </div>
    </div>
  );

  const renderContent = () => {
    switch (activeTab) {
      case 'general':
        return renderGeneralSettings();
      case 'features':
        return renderFeaturesSettings();
      case 'email':
        return renderEmailSettings();
      case 'security':
        return renderSecuritySettings();
      case 'notifications':
        return renderNotificationSettings();
      case 'payments':
        return renderPaymentSettings();
      default:
        return null;
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div>
          <h1 className={styles.title}>{t('title')}</h1>
          <p className={styles.subtitle}>{t('subtitle')}</p>
        </div>
        <div className={styles.actions}>
          <button 
            className={styles.saveButton}
            onClick={handleSave}
            disabled={!hasChanges || saveMutation.isLoading}
          >
            {saveMutation.isLoading ? (
              <RefreshCw size={20} className={styles.spinning} />
            ) : (
              <Save size={20} />
            )}
            {t('save')}
          </button>
        </div>
      </div>

      <div className={styles.content}>
        <div className={styles.sidebar}>
          {tabs.map(tab => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                className={`${styles.tabButton} ${activeTab === tab.id ? styles.active : ''}`}
                onClick={() => setActiveTab(tab.id)}
              >
                <Icon size={20} />
                {tab.label}
              </button>
            );
          })}
        </div>

        <div className={styles.mainContent}>
          {renderContent()}
        </div>
      </div>
    </div>
  );
};

export default AdminSettings;