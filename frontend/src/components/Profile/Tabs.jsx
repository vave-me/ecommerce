// src/components/UserProfile/Tabs.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import styles from './Tabs.module.css';
const Tabs = memo(({ activeTab, setActiveTab }) => {
    const tabs = ['Items', 'Posts', 'Comments', 'Reviews', 'Contact', 'Gallery'];
    return (
        <div className={styles.tabsContainer} role="tablist">
            {tabs.map((tab) => {
                const isActive = activeTab === tab;
                return (
                    <button
                        key={tab}
                        role="tab"
                        aria-selected={isActive}
                        className={`${styles.tab} ${isActive ? styles.active : ''}`}
                        onClick={() => setActiveTab(tab)}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter' || e.key === ' ') {
                                setActiveTab(tab);
                            }
                        }}
                        tabIndex={0}
                    >
                        {tab}
                    </button>
                );
            })}
        </div>
    );
});
Tabs.displayName = 'Tabs';
Tabs.propTypes = {
    activeTab: PropTypes.string.isRequired,
    setActiveTab: PropTypes.func.isRequired,
};
export default Tabs;
