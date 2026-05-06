"use client";
import React from 'react';
import { useTranslations } from 'next-intl';
import styles from './about.module.css';
import {
    Users,
    Target,
    Heart,
    Shield,
    Globe,
    Clock,
    Award,
    MessageSquare,
    Star
} from '@/icons';
export default function AboutPage() {
    const t = useTranslations('About');
    const sections = [
        {
            icon: <Users size={24} />,
            title: t('ourStory.title'),
            content: t('ourStory.content')
        },
        {
            icon: <Target size={24} />,
            title: t('mission.title'),
            content: t('mission.content')
        },
        {
            icon: <Heart size={24} />,
            title: t('values.title'),
            content: t('values.content')
        },
        {
            icon: <Shield size={24} />,
            title: t('security.title'),
            content: t('security.content')
        },
        {
            icon: <Globe size={24} />,
            title: t('global.title'),
            content: t('global.content')
        },
        {
            icon: <Clock size={24} />,
            title: t('support.title'),
            content: t('support.content')
        },
        {
            icon: <Award size={24} />,
            title: t('achievements.title'),
            content: t('achievements.content')
        },
        {
            icon: <Star size={24} />,
            title: t('future.title'),
            content: t('future.content')
        }
    ];
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>{t('pageTitle')}</h1>
                <p className={styles.subtitle}>{t('pageDescription')}</p>
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