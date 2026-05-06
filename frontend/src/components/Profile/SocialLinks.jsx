// src/components/UserProfile/SocialLinks.jsx
"use client"
import React, { useEffect, useState, memo } from 'react';
import PropTypes from 'prop-types';
import { Globe } from '@/icons';
import { FaLinkedin, FaTwitter, FaGithub, FaInstagram, FaFacebook } from '../../utils/iconImports';
import { getSocialLinks, getUserById } from "../../api/client/userApi";
import styles from './SocialLinks.module.css';
const SocialLinks = memo(({ userId, user }) => {
    const [links, setLinks] = useState({});
    const [loading, setLoading] = useState(true);
    useEffect(() => {
        const fetchSocialLinks = async () => {
            try {
                setLoading(true);
                // First try to get real user data that might contain social links in the future
                if (!user) {
                    try {
                        const userData = await getUserById(userId);
                        // If API starts supporting social links, extract them here
                    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
                }
                // Get social links from API
                const data = await getSocialLinks(userId);
                setLinks(data.links);
            } catch (error) {
                // Return empty links on error instead of dummy data
                setLinks({});
            } finally {
                setLoading(false);
            }
        };
        fetchSocialLinks();
    }, [userId, user]);
    if (loading) {
        return <div className={styles.loading}>Loading social links...</div>;
    }
    // If no links available, show message
    if (!links || Object.keys(links).length === 0) {
        return <div className={styles.noLinks}>No social links available</div>;
    }
    return (
        <div className={styles.socialContainer}>
            {links.linkedin && (
                <a
                    className={styles.socialLink}
                    href={links.linkedin}
                    target="_blank"
                    rel="noopener noreferrer"
                    aria-label="LinkedIn"
                >
                    <FaLinkedin aria-hidden="true" />
                </a>
            )}
            {links.twitter && (
                <a
                    className={styles.socialLink}
                    href={links.twitter}
                    target="_blank"
                    rel="noopener noreferrer"
                    aria-label="Twitter"
                >
                    <FaTwitter aria-hidden="true" />
                </a>
            )}
            {links.github && (
                <a
                    className={styles.socialLink}
                    href={links.github}
                    target="_blank"
                    rel="noopener noreferrer"
                    aria-label="GitHub"
                >
                    <FaGithub aria-hidden="true" />
                </a>
            )}
            {links.instagram && (
                <a
                    className={styles.socialLink}
                    href={links.instagram}
                    target="_blank"
                    rel="noopener noreferrer"
                    aria-label="Instagram"
                >
                    <FaInstagram aria-hidden="true" />
                </a>
            )}
            {links.facebook && (
                <a
                    className={styles.socialLink}
                    href={links.facebook}
                    target="_blank"
                    rel="noopener noreferrer"
                    aria-label="Facebook"
                >
                    <FaFacebook aria-hidden="true" />
                </a>
            )}
            {links.website && (
                <a
                    className={styles.socialLink}
                    href={links.website}
                    target="_blank"
                    rel="noopener noreferrer"
                    aria-label="Website"
                >
                    <Globe aria-hidden="true" />
                </a>
            )}
        </div>
    );
});
SocialLinks.displayName = 'SocialLinks';
SocialLinks.propTypes = {
    userId: PropTypes.string.isRequired,
    user: PropTypes.object
};
export default SocialLinks;
