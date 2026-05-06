import React from 'react';

// --- SVG Icon Sub-components for this section ---
const TransformIcon = () => ( <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21.174 6.786a2.433 2.433 0 0 0-3.44 0L6.786 17.734a2.433 2.433 0 0 0 0 3.44l.172.172a2.433 2.433 0 0 0 3.44 0L21.346 10.4a2.433 2.433 0 0 0 0-3.44l-.172-.172zM2.654 10.4a2.433 2.433 0 0 1 0-3.44l.172-.172a2.433 2.433 0 0 1 3.44 0L17.214 17.734a2.433 2.433 0 0 1 0 3.44l-.172.172a2.433 2.433 0 0 1-3.44 0L2.654 10.4z"></path></svg> );
const MigrateIcon = () => ( <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 17v-4a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v4M12 11V3M9 6l3-3 3 3M21 17H3"></path></svg> );
const LegacyIcon = () => ( <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.72"></path><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.72-1.72"></path></svg> );
const SaveIcon = () => ( <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg> );

const TransformBusinessSection = () => {

    const styles = {
        section: {
            maxWidth: '1200px',
            margin: '0 auto',
            padding: '80px 24px',
            backgroundColor: '#f9fafb',
        },
        sectionTitle: {
            fontSize: '36px',
            fontWeight: 'bold',
            marginBottom: '16px',
            textAlign: 'center',
            color: '#111827',
        },
        sectionSubTitle: {
            fontSize: '18px',
            color: '#4b5563',
            textAlign: 'center',
            maxWidth: '800px',
            margin: '0 auto 64px auto',
        },
        grid: {
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
            gap: '32px',
        },
        card: {
            textAlign: 'center',
            backgroundColor: '#ffffff',
            padding: '32px',
            borderRadius: '12px',
            border: '1px solid #e5e7eb',
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.05)',
        },
        iconWrapper: {
            color: '#4f46e5',
            marginBottom: '24px',
        },
        cardTitle: {
            fontSize: '22px',
            fontWeight: 'bold',
            marginBottom: '12px',
            color: '#1f2937',
        },
        cardText: {
            color: '#4b5563',
            fontSize: '16px',
        }
    };

    return (
        <section style={styles.section}>
            <h2 style={styles.sectionTitle}>Upgrade Your Digital Foundation</h2>
            <p style={styles.sectionSubTitle}>
                Our modern, AI-native architecture isn't just an improvement—it's a complete transformation of how you build, operate, and scale your digital business.
            </p>
            <div style={styles.grid}>
                {/* Card 1: TRANSFORM */}
                <div style={styles.card}>
                    <div style={styles.iconWrapper}><TransformIcon /></div>
                    <h3 style={styles.cardTitle}>TRANSFORM User Experiences</h3>
                    <p style={styles.cardText}>
                        Move from rigid clicks to natural conversations. Our AI Co-Pilot understands user intent, not just keywords, enabling a fluid, voice-first experience that builds loyalty and drives conversion.
                    </p>
                </div>
                {/* Card 2: MIGRATE */}
                <div style={styles.card}>
                    <div style={styles.iconWrapper}><MigrateIcon /></div>
                    <h3 style={styles.cardTitle}>MIGRATE With Confidence</h3>
                    <p style={styles.cardText}>
                        Transitioning to a new system is seamless. We partner with you to ensure a smooth migration, and can even leverage our AI to intelligently categorize and tag your legacy product data, saving hundreds of manual hours.
                    </p>
                </div>
                {/* Card 3: DROP LEGACY */}
                <div style={styles.card}>
                    <div style={styles.iconWrapper}><LegacyIcon /></div>
                    <h3 style={styles.cardTitle}>DROP Your Legacy Constraints</h3>
                    <p style={styles.cardText}>
                        Break free from monolithic technical debt. Our composable, microservice architecture gives you the agility to launch new verticals, test ideas, and innovate at a speed your competitors can't match.
                    </p>
                </div>
                {/* Card 4: SAVE */}
                <div style={styles.card}>
                    <div style={styles.iconWrapper}><SaveIcon /></div>
                    <h3 style={styles.cardTitle}>SAVE Money and Resources</h3>
                    <p style={styles.cardText}>
                        Eliminate major recurring costs with our proprietary media and geocoding services. Our efficient, auto-scaling infrastructure lowers your TCO and frees your developers to build value, not manage servers.
                    </p>
                </div>
            </div>
        </section>
    );
};

export default TransformBusinessSection;