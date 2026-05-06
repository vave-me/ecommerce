// src/components/UserProfile/ContactInfo.jsx
"use client"
import React, { useEffect, useState, memo } from 'react';
import PropTypes from 'prop-types';
import { Mail, Phone, MapPin } from '@/icons';
import { getUserById } from "../../api/client/userApi";
import styles from './ContactInfo.module.css';
const ContactInfo = memo(({ userId, user }) => {
    const [contactInfo, setContactInfo] = useState({});
    const [loading, setLoading] = useState(true);
    useEffect(() => {
        const fetchContactInfo = async () => {
            try {
                setLoading(true);
                // Get user data from API
                const userData = user || await getUserById(userId);
                // Extract contact info from user data
                const info = {
                    email: userData.email || null,
                    phone: userData.phone || null,
                    location: userData.location || userData.city || null
                };
                setContactInfo(info);
            } catch (error) {
                // Return empty contact info on error
                setContactInfo({});
            } finally {
                setLoading(false);
            }
        };
        fetchContactInfo();
    }, [userId, user]);
    if (loading) {
        return <div className={styles.loading}>Loading contact info...</div>;
    }
    // If no contact info available, show message
    const hasContactInfo = contactInfo.email || contactInfo.phone || contactInfo.location;
    if (!hasContactInfo) {
        return <div className={styles.noContact}>No contact information available</div>;
    }
    return (
        <div className={styles.contactContainer}>
            {contactInfo.email && (
                <div className={styles.contactItem}>
                    <Mail className={styles.contactIcon} size={16} />
                    <a href={`mailto:${contactInfo.email}`} className={styles.contactLink}>
                        {contactInfo.email}
                    </a>
                </div>
            )}
            {contactInfo.phone && (
                <div className={styles.contactItem}>
                    <Phone className={styles.contactIcon} size={16} />
                    <a href={`tel:${contactInfo.phone}`} className={styles.contactLink}>
                        {contactInfo.phone}
                    </a>
                </div>
            )}
            {contactInfo.location && (
                <div className={styles.contactItem}>
                    <MapPin className={styles.contactIcon} size={16} />
                    <span className={styles.contactText}>{contactInfo.location}</span>
                </div>
            )}
        </div>
    );
});
ContactInfo.displayName = 'ContactInfo';
ContactInfo.propTypes = {
    userId: PropTypes.string.isRequired,
    user: PropTypes.object
};
export default ContactInfo;
