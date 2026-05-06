// src/components/NotificationSettings.jsx
import React, { useState, memo } from 'react';
import PropTypes from 'prop-types';
import styles from './NotificationSettings.module.css';
const NotificationSettings = memo(({ settings, onUpdateSettings }) => {
    const [localSettings, setLocalSettings] = useState(settings);
    const handleToggle = (type) => {
        const updatedSettings = {
            ...localSettings,
            [type]: !localSettings[type],
        };
        setLocalSettings(updatedSettings);
    };
    const handleSave = () => {
        onUpdateSettings(localSettings);
    };
    return (
        <div className={styles.settingsContainer}>
            <h3 className={styles.settingsTitle}>Notification Preferences</h3>
            <ul className={styles.settingsList}>
                {Object.keys(localSettings).map((type) => (
                    <li className={styles.settingsItem} key={type}>
                        <label className={styles.label}>
                            <input
                                className={styles.checkbox}
                                type="checkbox"
                                checked={localSettings[type]}
                                onChange={() => handleToggle(type)}
                            />
                            {type.replace('_', ' ')}
                        </label>
                    </li>
                ))}
            </ul>
            <button className={styles.saveButton} onClick={handleSave}>Save Preferences</button>
        </div>
    );
});
NotificationSettings.displayName = 'NotificationSettings';
NotificationSettings.propTypes = {
    settings: PropTypes.object.isRequired,
    onUpdateSettings: PropTypes.func.isRequired,
};
export default NotificationSettings;