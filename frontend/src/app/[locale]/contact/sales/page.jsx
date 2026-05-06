"use client";

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  Phone,
  Mail,
  MapPin,
  Clock,
  CheckCircle,
  ArrowRight,
  Building2,
  Users,
  Target,
  Zap,
  Shield,
  Globe,
  Calendar,
  MessageSquare
} from 'lucide-react';
import styles from './ContactSales.module.css';

export default function ContactSalesPage() {
  const t = useTranslations('contactSales');
  const router = useRouter();
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    phone: '',
    company: '',
    title: '',
    companySize: '',
    industry: '',
    country: '',
    message: '',
    interests: []
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleInterestToggle = (interest) => {
    setFormData(prev => ({
      ...prev,
      interests: prev.interests.includes(interest)
        ? prev.interests.filter(i => i !== interest)
        : [...prev.interests, interest]
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    // Simulate form submission
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    setSubmitted(true);
    setIsSubmitting(false);
  };

  const interests = [
    'Platform Development',
    'API Integration',
    'Custom Solutions',
    'AI/ML Features',
    'Security & Compliance',
    'Migration Services'
  ];

  const benefits = [
    {
      icon: Zap,
      title: t('benefit1Title', 'Lightning Fast Implementation'),
      description: t('benefit1Desc', 'Get up and running in weeks, not months')
    },
    {
      icon: Shield,
      title: t('benefit2Title', 'Enterprise Security'),
      description: t('benefit2Desc', 'Bank-grade security with SOC2 compliance')
    },
    {
      icon: Globe,
      title: t('benefit3Title', 'Global Scale'),
      description: t('benefit3Desc', 'Infrastructure that grows with your business')
    },
    {
      icon: Users,
      title: t('benefit4Title', 'Dedicated Support'),
      description: t('benefit4Desc', '24/7 technical support with SLA guarantees')
    }
  ];

  const offices = [
    {
      location: t('office1Location', 'Berlin, Germany'),
      address: t('office1Address', 'Friedrichstraße 123, 10117 Berlin'),
      phone: t('office1Phone', '+49 30 1234 5678'),
      email: 'redacted-email@example.com'
    },
    {
      location: t('office2Location', 'Milan, Italy'),
      address: t('office2Address', 'Via Monte Napoleone 8, 20121 Milano'),
      phone: t('office2Phone', '+39 02 1234 5678'),
      email: 'redacted-email@example.com'
    },
    {
      location: t('office3Location', 'Warsaw, Poland'),
      address: t('office3Address', 'ul. Marszałkowska 111, 00-102 Warszawa'),
      phone: t('office3Phone', '+48 22 1234 5678'),
      email: 'redacted-email@example.com'
    }
  ];

  if (submitted) {
    return (
      <div className={styles.container}>
        <div className={styles.successMessage}>
          <CheckCircle className={styles.successIcon} />
          <h1>{t('successTitle', 'Thank You for Your Interest!')}</h1>
          <p>{t('successMessage', 'Our enterprise sales team will contact you within 24 hours.')}</p>
          <button 
            onClick={() => router.push('/resources/whitepaper')}
            className={styles.primaryButton}
          >
            {t('downloadWhitepaper', 'Download Whitepaper')}
            <ArrowRight size={16} />
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Hero Section */}
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <h1 className={styles.heroTitle}>
            {t('heroTitle', 'Talk to Our Enterprise Sales Team')}
          </h1>
          <p className={styles.heroSubtitle}>
            {t('heroSubtitle', 'Discover how our platform can transform your business with custom solutions tailored to your needs')}
          </p>
        </div>
      </section>

      <div className={styles.mainContent}>
        <div className={styles.formSection}>
          <div className={styles.formHeader}>
            <h2>{t('formTitle', 'Request a Demo')}</h2>
            <p>{t('formSubtitle', 'Fill out the form below and our team will contact you shortly')}</p>
          </div>

          <form onSubmit={handleSubmit} className={styles.contactForm}>
            <div className={styles.formRow}>
              <div className={styles.formGroup}>
                <label htmlFor="firstName">{t('firstName', 'First Name')} *</label>
                <input
                  type="text"
                  id="firstName"
                  name="firstName"
                  value={formData.firstName}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="lastName">{t('lastName', 'Last Name')} *</label>
                <input
                  type="text"
                  id="lastName"
                  name="lastName"
                  value={formData.lastName}
                  onChange={handleChange}
                  required
                />
              </div>
            </div>

            <div className={styles.formRow}>
              <div className={styles.formGroup}>
                <label htmlFor="email">{t('email', 'Business Email')} *</label>
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
                <label htmlFor="phone">{t('phone', 'Phone Number')}</label>
                <input
                  type="tel"
                  id="phone"
                  name="phone"
                  value={formData.phone}
                  onChange={handleChange}
                />
              </div>
            </div>

            <div className={styles.formRow}>
              <div className={styles.formGroup}>
                <label htmlFor="company">{t('company', 'Company Name')} *</label>
                <input
                  type="text"
                  id="company"
                  name="company"
                  value={formData.company}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="title">{t('title', 'Job Title')} *</label>
                <input
                  type="text"
                  id="title"
                  name="title"
                  value={formData.title}
                  onChange={handleChange}
                  required
                />
              </div>
            </div>

            <div className={styles.formRow}>
              <div className={styles.formGroup}>
                <label htmlFor="companySize">{t('companySize', 'Company Size')} *</label>
                <select
                  id="companySize"
                  name="companySize"
                  value={formData.companySize}
                  onChange={handleChange}
                  required
                >
                  <option value="">{t('selectSize', 'Select company size')}</option>
                  <option value="1-50">1-50 {t('employees', 'employees')}</option>
                  <option value="51-200">51-200 {t('employees', 'employees')}</option>
                  <option value="201-500">201-500 {t('employees', 'employees')}</option>
                  <option value="501-1000">501-1000 {t('employees', 'employees')}</option>
                  <option value="1000+">1000+ {t('employees', 'employees')}</option>
                </select>
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="industry">{t('industry', 'Industry')} *</label>
                <select
                  id="industry"
                  name="industry"
                  value={formData.industry}
                  onChange={handleChange}
                  required
                >
                  <option value="">{t('selectIndustry', 'Select industry')}</option>
                  <option value="retail">{t('retail', 'Retail & E-commerce')}</option>
                  <option value="manufacturing">{t('manufacturing', 'Manufacturing')}</option>
                  <option value="technology">{t('technology', 'Technology')}</option>
                  <option value="healthcare">{t('healthcare', 'Healthcare')}</option>
                  <option value="finance">{t('finance', 'Finance')}</option>
                  <option value="energy">{t('energy', 'Energy & Utilities')}</option>
                  <option value="other">{t('other', 'Other')}</option>
                </select>
              </div>
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="country">{t('country', 'Country')} *</label>
              <input
                type="text"
                id="country"
                name="country"
                value={formData.country}
                onChange={handleChange}
                required
              />
            </div>

            <div className={styles.formGroup}>
              <label>{t('interests', 'Areas of Interest')}</label>
              <div className={styles.interestsGrid}>
                {interests.map((interest) => (
                  <label key={interest} className={styles.interestItem}>
                    <input
                      type="checkbox"
                      checked={formData.interests.includes(interest)}
                      onChange={() => handleInterestToggle(interest)}
                    />
                    <span>{t(interest.toLowerCase().replace(/\s+/g, ''), interest)}</span>
                  </label>
                ))}
              </div>
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="message">{t('message', 'Tell us about your project')}</label>
              <textarea
                id="message"
                name="message"
                value={formData.message}
                onChange={handleChange}
                rows={4}
                placeholder={t('messagePlaceholder', 'Please describe your business needs and goals...')}
              />
            </div>

            <button 
              type="submit" 
              className={styles.submitButton}
              disabled={isSubmitting}
            >
              {isSubmitting ? t('submitting', 'Submitting...') : t('submit', 'Request Demo')}
              <ArrowRight size={16} />
            </button>
          </form>
        </div>

        <div className={styles.sidebarSection}>
          {/* Benefits */}
          <div className={styles.benefitsCard}>
            <h3>{t('whyChooseUs', 'Why Choose Our Platform')}</h3>
            <div className={styles.benefitsList}>
              {benefits.map((benefit, index) => (
                <div key={index} className={styles.benefitItem}>
                  <benefit.icon className={styles.benefitIcon} />
                  <div>
                    <h4>{benefit.title}</h4>
                    <p>{benefit.description}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Contact Information */}
          <div className={styles.contactCard}>
            <h3>{t('contactInfo', 'Direct Contact')}</h3>
            <div className={styles.contactItem}>
              <Phone className={styles.contactIcon} />
              <div>
                <p className={styles.contactLabel}>{t('callUs', 'Call Us')}</p>
                <a href="tel:+498001234567">+49 800 123 4567</a>
              </div>
            </div>
            <div className={styles.contactItem}>
              <Mail className={styles.contactIcon} />
              <div>
                <p className={styles.contactLabel}>{t('emailUs', 'Email Us')}</p>
                <a href="mailto:redacted-email@example.com">redacted-email@example.com</a>
              </div>
            </div>
            <div className={styles.contactItem}>
              <Clock className={styles.contactIcon} />
              <div>
                <p className={styles.contactLabel}>{t('businessHours', 'Business Hours')}</p>
                <p>{t('hours', 'Mon-Fri 9:00-18:00 CET')}</p>
              </div>
            </div>
          </div>

          {/* Schedule Call */}
          <div className={styles.scheduleCard}>
            <Calendar className={styles.scheduleIcon} />
            <h3>{t('scheduleCall', 'Schedule a Call')}</h3>
            <p>{t('scheduleDesc', 'Book a 30-minute consultation with our experts')}</p>
            <button 
              onClick={() => window.open('https://calendly.com/platform-sales', '_blank')}
              className={styles.scheduleButton}
            >
              {t('bookMeeting', 'Book Meeting')}
              <ArrowRight size={16} />
            </button>
          </div>
        </div>
      </div>

      {/* Office Locations */}
      <section className={styles.offices}>
        <h2 className={styles.sectionTitle}>{t('officesTitle', 'Our Offices')}</h2>
        <div className={styles.officesGrid}>
          {offices.map((office, index) => (
            <div key={index} className={styles.officeCard}>
              <MapPin className={styles.officeIcon} />
              <h3>{office.location}</h3>
              <p>{office.address}</p>
              <a href={`tel:${office.phone}`}>{office.phone}</a>
              <a href={`mailto:${office.email}`}>{office.email}</a>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}