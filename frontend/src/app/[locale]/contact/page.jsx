"use client";
import React, {useState} from 'react';
import {useTranslations} from 'next-intl';
import styles from './contact.module.css';
import {
    Mail,
    Phone,
    MapPin,
    Clock,
    Send,
    AlertCircle,
    CheckCircle
} from '@/icons';
export default function ContactPage() {
    const t = useTranslations('Contact');
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        subject: '',
        message: ''
    });
    const [status, setStatus] = useState({type: '', message: ''});
    const handleChange = (e) => {
        const {name, value} = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };
    const handleSubmit = async (e) => {
        e.preventDefault();
        setStatus({type: 'loading', message: t('form.submitting')});
        try {
            // Here you would typically send the form data to your backend
            await new Promise(resolve => setTimeout(resolve, 1000)); // Simulated API call
            setStatus({type: 'success', message: t('form.success')});
            setFormData({name: '', email: '', subject: '', message: ''});
        } catch (error) {
            setStatus({type: 'error', message: t('form.error')});
        }
    };
    const contactInfo = [
        {
            icon: <Mail size={24}/>,
            title: t('info.email.title'),
            content: t('info.email.content'),
            link: 'mailto:redacted-email@example.com'
        },
        {
            icon: <Phone size={24}/>,
            title: t('info.phone.title'),
            content: t('info.phone.content'),
            link: 'tel:+1234567890'
        },
        {
            icon: <MapPin size={24}/>,
            title: t('info.address.title'),
            content: t('info.address.content')
        },
        {
            icon: <Clock size={24}/>,
            title: t('info.hours.title'),
            content: t('info.hours.content')
        }
    ];
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>{t('pageTitle')}</h1>
                <p className={styles.subtitle}>{t('pageDescription')}</p>
            </div>
            <div className={styles.content}>
                <div className={styles.contactInfo}>
                    {contactInfo.map((info, index) => (
                        <div key={index} className={styles.infoCard}>
                            <div className={styles.iconWrapper}>
                                {info.icon}
                            </div>
                            <h3 className={styles.infoTitle}>{info.title}</h3>
                            {info.link ? (
                                <a href={info.link} className={styles.infoLink}>
                                    {info.content}
                                </a>
                            ) : (
                                <p className={styles.infoContent}>{info.content}</p>
                            )}
                        </div>
                    ))}
                </div>
                <div className={styles.formContainer}>
                    <h2 className={styles.formTitle}>{t('form.title')}</h2>
                    <form onSubmit={handleSubmit} className={styles.form}>
                        <div className={styles.formGroup}>
                            <label htmlFor="name">{t('form.name')}</label>
                            <input
                                type="text"
                                id="name"
                                name="name"
                                value={formData.name}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className={styles.formGroup}>
                            <label htmlFor="email">{t('form.email')}</label>
                            <input
                                type="email"
                                id="email"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className={styles.formGroup}>
                            <label htmlFor="subject">{t('form.subject')}</label>
                            <input
                                type="text"
                                id="subject"
                                name="subject"
                                value={formData.subject}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className={styles.formGroup}>
                            <label htmlFor="message">{t('form.message')}</label>
                            <textarea
                                id="message"
                                name="message"
                                value={formData.message}
                                onChange={handleChange}
                                required
                                rows="5"
                            />
                        </div>
                        {status.type && (
                            <div className={`${styles.status} ${styles[status.type]}`}>
                                {status.type === 'error' ? <AlertCircle size={20}/> : <CheckCircle size={20}/>}
                                <span>{status.message}</span>
                            </div>
                        )}
                        <button type="submit" className={styles.submitButton} disabled={status.type === 'loading'}>
                            <Send size={20}/>
                            {t('form.submit')}
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
}