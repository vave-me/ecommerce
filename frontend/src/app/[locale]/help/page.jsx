"use client";
import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import styles from './help.module.css';
import { 
    Search, 
    ChevronDown, 
    ChevronUp,
    HelpCircle,
    MessageSquare,
    Shield,
    CreditCard,
    User,
    Settings,
    AlertCircle
} from '@/icons';
export default function HelpPage() {
    const t = useTranslations('Help');
    const [searchQuery, setSearchQuery] = useState('');
    const [expandedFaq, setExpandedFaq] = useState(null);
    const categories = [
        {
            icon: <HelpCircle size={24} />,
            title: t('categories.general.title'),
            description: t('categories.general.description')
        },
        {
            icon: <User size={24} />,
            title: t('categories.account.title'),
            description: t('categories.account.description')
        },
        {
            icon: <CreditCard size={24} />,
            title: t('categories.payments.title'),
            description: t('categories.payments.description')
        },
        {
            icon: <Shield size={24} />,
            title: t('categories.security.title'),
            description: t('categories.security.description')
        },
        {
            icon: <MessageSquare size={24} />,
            title: t('categories.communication.title'),
            description: t('categories.communication.description')
        },
        {
            icon: <Settings size={24} />,
            title: t('categories.technical.title'),
            description: t('categories.technical.description')
        }
    ];
    const faqs = [
        {
            question: t('faq.howToCreateAccount.question'),
            answer: t('faq.howToCreateAccount.answer')
        },
        {
            question: t('faq.howToResetPassword.question'),
            answer: t('faq.howToResetPassword.answer')
        },
        {
            question: t('faq.paymentMethods.question'),
            answer: t('faq.paymentMethods.answer')
        },
        {
            question: t('faq.securityMeasures.question'),
            answer: t('faq.securityMeasures.answer')
        },
        {
            question: t('faq.contactSupport.question'),
            answer: t('faq.contactSupport.answer')
        },
        {
            question: t('faq.technicalIssues.question'),
            answer: t('faq.technicalIssues.answer')
        }
    ];
    const filteredFaqs = faqs.filter(faq => 
        faq.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
        faq.answer.toLowerCase().includes(searchQuery.toLowerCase())
    );
    const toggleFaq = (index) => {
        setExpandedFaq(expandedFaq === index ? null : index);
    };
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>{t('pageTitle')}</h1>
                <p className={styles.subtitle}>{t('pageDescription')}</p>
            </div>
            <div className={styles.searchContainer}>
                <div className={styles.searchWrapper}>
                    <Search size={20} className={styles.searchIcon} />
                    <input
                        type="text"
                        placeholder={t('searchPlaceholder')}
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className={styles.searchInput}
                    />
                </div>
            </div>
            <div className={styles.categories}>
                {categories.map((category, index) => (
                    <div key={index} className={styles.categoryCard}>
                        <div className={styles.iconWrapper}>
                            {category.icon}
                        </div>
                        <h3 className={styles.categoryTitle}>{category.title}</h3>
                        <p className={styles.categoryDescription}>{category.description}</p>
                    </div>
                ))}
            </div>
            <div className={styles.faqSection}>
                <h2 className={styles.faqTitle}>{t('faqTitle')}</h2>
                <div className={styles.faqList}>
                    {filteredFaqs.map((faq, index) => (
                        <div key={index} className={styles.faqItem}>
                            <button
                                className={styles.faqQuestion}
                                onClick={() => toggleFaq(index)}
                            >
                                <span>{faq.question}</span>
                                {expandedFaq === index ? <ChevronUp size={20} /> : <ChevronDown size={20} />}
                            </button>
                            {expandedFaq === index && (
                                <div className={styles.faqAnswer}>
                                    <p>{faq.answer}</p>
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            </div>
            <div className={styles.contactSection}>
                <div className={styles.contactCard}>
                    <AlertCircle size={24} className={styles.contactIcon} />
                    <h3 className={styles.contactTitle}>{t('contact.title')}</h3>
                    <p className={styles.contactText}>{t('contact.description')}</p>
                    <a href="/contact" className={styles.contactButton}>
                        {t('contact.button')}
                    </a>
                </div>
            </div>
        </div>
    );
} 