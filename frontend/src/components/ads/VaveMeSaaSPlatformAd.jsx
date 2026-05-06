"use client"
import React from 'react';
import Image from 'next/image';
import CompactVaveAd from './CompactVaveAd';

// --- SVG Icon Sub-components for a clean, self-contained design ---
// (In a real app, these might be in their own files)

const CodeIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <polyline points="16 18 22 12 16 6"></polyline>
        <polyline points="8 6 2 12 8 18"></polyline>
    </svg>
);
const AiIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M12 8V4H8"></path>
        <rect x="4" y="8" width="8" height="12" rx="2"></rect>
        <path d="M2 14h2"></path>
        <path d="M20 14h2"></path>
        <path d="M15 13v2"></path>
        <path d="M15 5v2"></path>
    </svg>
);
const ChatIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
    </svg>
);
const ServerIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
        <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
        <line x1="6" y1="6" x2="6.01" y2="6"></line>
        <line x1="6" y1="18" x2="6.01" y2="18"></line>
    </svg>
);

// --- The Main Ad Component ---
const VaveMeSaaSPlatformAd = () => {
    // Styles are defined as JS objects for portability
    const styles = {
        container: {
            backgroundColor: '#f9fafb', // Light grey background
            padding: '64px 24px',
            fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif',
            color: '#1f2937',
        },
        contentWrapper: {
            maxWidth: '1100px',
            margin: '0 auto',
            textAlign: 'center',
        },
        headline: {
            fontSize: '42px',
            fontWeight: 'bold',
            lineHeight: '1.2',
            marginBottom: '16px',
        },
        subHeadline: {
            fontSize: '20px',
            color: '#4b5563',
            maxWidth: '700px',
            margin: '0 auto 48px auto',
        },
        sectionTitle: {
            fontSize: '32px',
            fontWeight: 'bold',
            marginBottom: '48px',
            marginTop: '64px',
        },
        featuresGrid: {
            display: 'flex',
            flexWrap: 'wrap',
            gap: '24px',
            justifyContent: 'center',
            textAlign: 'left',
        },
        featureCard: {
            backgroundColor: '#ffffff',
            padding: '24px',
            borderRadius: '12px',
            border: '1px solid #e5e7eb',
            width: '100%',
            maxWidth: '500px',
            boxSizing: 'border-box',
        },
        iconWrapper: {
            color: '#4f46e5',
            marginBottom: '16px',
        },
        cardTitle: {
            fontSize: '20px',
            fontWeight: 'bold',
            marginBottom: '8px',
        },
        cardText: {
            color: '#4b5563',
            lineHeight: '1.6',
        },
        tableContainer: {
            maxWidth: '800px',
            margin: '0 auto',
            border: '1px solid #e5e7eb',
            borderRadius: '12px',
            overflow: 'hidden',
        },
        table: {
            width: '100%',
            borderCollapse: 'collapse',
            textAlign: 'left',
        },
        th: {
            padding: '16px',
            backgroundColor: '#f3f4f6',
            fontWeight: 'bold',
        },
        td: {
            padding: '16px',
            borderTop: '1px solid #e5e7eb',
        },
        tr: {
            backgroundColor: '#ffffff',
        },
        ctaButton: {
            display: 'inline-block',
            backgroundColor: '#4f46e5',
            color: '#ffffff',
            border: 'none',
            borderRadius: '8px',
            padding: '16px 32px',
            fontSize: '18px',
            fontWeight: 'bold',
            cursor: 'pointer',
            textDecoration: 'none',
            transition: 'background-color 0.2s ease',
            marginTop: '24px',
        },
    };

    return (
        <div style={styles.container}>
            <div style={styles.contentWrapper}>

                {/* --- HERO SECTION --- */}
                <div style={{ marginBottom: '32px' }}>
                    <Image 
                        src="/images/logo-vaveme.png" 
                        alt="Vave.me Logo" 
                        width={120} 
                        height={60}
                        style={{ objectFit: 'contain' }}
                    />
                </div>
                <h1 style={styles.headline}>Enterprise-Grade Marketplace Infrastructure for Modern Commerce</h1>
                <p style={styles.subHeadline}>
                    Deploy white-label B2B, B2C, or hybrid marketplaces with enterprise security, compliance, and scalability. 
                    Trusted by Fortune 500 companies and high-growth startups to accelerate digital transformation.
                </p>

                {/* --- FEATURES SECTION --- */}
                <h2 style={styles.sectionTitle}>Everything You Need to Succeed, Right Out of the Box</h2>
                <div style={styles.featuresGrid}>
                    <div style={styles.featureCard}>
                        <div style={styles.iconWrapper}><CodeIcon/></div>
                        <h3 style={styles.cardTitle}>Industry-Specific Solutions</h3>
                        <p style={styles.cardText}>Pre-configured for Automotive, Real Estate, Healthcare, Manufacturing, and Professional Services.
                            Includes compliance frameworks, specialized workflows, and industry-standard integrations.</p>
                    </div>
                    <div style={styles.featureCard}>
                        <div style={styles.iconWrapper}><AiIcon/></div>
                        <h3 style={styles.cardTitle}>AI-Powered Business Intelligence</h3>
                        <p style={styles.cardText}>Advanced analytics, predictive insights, and automated workflows. Reduce operational costs by 40%
                            while improving vendor onboarding, fraud detection, and customer satisfaction metrics.</p>
                    </div>
                    <div style={styles.featureCard}>
                        <div style={styles.iconWrapper}><ChatIcon/></div>
                        <h3 style={styles.cardTitle}>B2B Commerce Features</h3>
                        <p style={styles.cardText}>Quote management, bulk ordering, custom pricing tiers, approval workflows,
                            and integration with ERP/CRM systems. Built for complex B2B transactions and relationships.</p>
                    </div>
                    <div style={styles.featureCard}>
                        <div style={styles.iconWrapper}><ServerIcon/></div>
                        <h3 style={styles.cardTitle}>Security & Compliance</h3>
                        <p style={styles.cardText}>SOC 2 Type II certified, GDPR compliant, PCI DSS Level 1. Enterprise SSO, role-based access,
                            audit trails, and data encryption. Meet the strictest corporate and regulatory requirements.</p>
                    </div>
                </div>

                {/* --- COMPARISON TABLE SECTION --- */}
                <h2 style={styles.sectionTitle}>The Smart Choice for Forward-Thinking Entrepreneurs</h2>
                <div style={styles.tableContainer}>
                    <table style={styles.table}>
                        <thead>
                        <tr>
                            <th style={styles.th}>Feature</th>
                            <th style={styles.th}>Traditional Approach</th>
                            <th style={styles.th}>With sfx-markt.de</th>
                        </tr>
                        </thead>
                        <tbody>
                        <tr style={styles.tr}>
                            <td style={styles.td}>Time to Market</td>
                            <td style={styles.td}>18-24 months of development</td>
                            <td style={styles.td}><strong>Launch in 4-6 weeks</strong> ✅</td>
                        </tr>
                        <tr style={{...styles.tr, backgroundColor: '#f9fafb'}}>
                            <td style={styles.td}>Initial Investment</td>
                            <td style={styles.td}>€500K - €2M upfront investment</td>
                            <td style={styles.td}><strong>Affordable monthly plans</strong> ✅</td>
                        </tr>
                        <tr style={styles.tr}>
                            <td style={styles.td}>AI Integration</td>
                            <td style={styles.td}>Months of complex integration</td>
                            <td style={styles.td}><strong>Instant AI capabilities</strong> ✅</td>
                        </tr>
                        <tr style={{...styles.tr, backgroundColor: '#f9fafb'}}>
                            <td style={styles.td}>Payment Processing</td>
                            <td style={styles.td}>6+ months implementation</td>
                            <td style={styles.td}><strong>Go live immediately</strong> ✅</td>
                        </tr>
                        </tbody>
                    </table>
                </div>

                {/* --- CTA SECTION --- */}
                <h2 style={styles.sectionTitle}>Your Success Story Starts Today</h2>
                <p style={styles.subHeadline}>Join thousands of thriving marketplace operators who've transformed their industries with sfx-markt.de.
                    Discover how our platform can amplify your impact and accelerate your success.</p>
                <a
                    href="https://sfx-markt.de/contact-sales" // This would link to your contact or booking page
                    style={styles.ctaButton}
                    onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#4338ca'}
                    onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#4f46e5'}
                >
                    Experience the Future Today
                </a>

                {/* Compact Ad Component */}
                <div style={{ marginTop: '64px', display: 'flex', justifyContent: 'center' }}>
                    <CompactVaveAd />
                </div>

            </div>
        </div>
    );
};

export default VaveMeSaaSPlatformAd;

