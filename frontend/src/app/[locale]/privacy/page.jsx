"use client";
import React from 'react';
import {useTranslations} from 'next-intl';
import styles from './page.module.css';
import {Shield, Lock, Eye, UserCheck, Server, Bell, RefreshCw, Trash2} from '@/icons';
export default function PrivacyPolicy() {
    const t = useTranslations('Privacy');
    const sections = [
        {
            icon: <Shield size={24}/>,
            title: t('dataCollection.title'),
            content: t('dataCollection.content')
        },
        {
            icon: <Lock size={24}/>,
            title: t('dataSecurity.title'),
            content: t('dataSecurity.content')
        },
        {
            icon: <Eye size={24}/>,
            title: t('dataUsage.title'),
            content: t('dataUsage.content')
        },
        {
            icon: <UserCheck size={24}/>,
            title: t('userRights.title'),
            content: t('userRights.content')
        },
        {
            icon: <Server size={24}/>,
            title: t('dataStorage.title'),
            content: t('dataStorage.content')
        },
        {
            icon: <Bell size={24}/>,
            title: t('communications.title'),
            content: t('communications.content')
        },
        {
            icon: <RefreshCw size={24}/>,
            title: t('policyUpdates.title'),
            content: t('policyUpdates.content')
        },
        {
            icon: <Trash2 size={24}/>,
            title: t('dataDeletion.title'),
            content: t('dataDeletion.content')
        }
    ];
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>{t('pageTitle')}</h1>
                <p className={styles.subtitle}>{t('pageDescription')}</p>
                <div className={styles.lastUpdated}>
                    {t('lastUpdated', {date: new Date().toLocaleDateString()})}
                </div>
            </div>
            <div className={styles.content}>
                <div className={styles.introduction}>
                    <p>{t('introduction')}</p>
                </div>
                <div className={styles.sections}>
                    {sections.map((section, index) => (
                        <div key={index} className={styles.section}>
                            <div className={styles.sectionHeader}>
                                <div className={styles.iconWrapper}>
                                    {section.icon}
                                </div>
                                <h2 className={styles.sectionTitle}>{section.title}</h2>
                            </div>
                            <div className={styles.sectionContent}>
                                <p>{section.content}</p>
                            </div>
                        </div>
                    ))}
                </div>
                <div className={styles.contact}>
                    <h2 className={styles.contactTitle}>{t('contact.title')}</h2>
                    <p className={styles.contactText}>{t('contact.content')}</p>
                    <a href="mailto:redacted-email@example.com" className={styles.contactLink}>
                        redacted-email@example.com
                    </a>
                </div>
            </div>
        </div>
    );
}