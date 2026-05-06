import React, { memo } from 'react';
import { MapPin, Building2, Home, Map } from '@/icons';
import styles from './LocationSuggestion.module.css';
const LocationSuggestion = memo(({ 
    location, 
    distance, 
    type = 'city',
    onClick,
    isActive = false 
}) => {
    const getIcon = () => {
        switch (type) {
            case 'city':
                return <Building2 size={16} />;
            case 'district':
                return <Home size={16} />;
            case 'region':
                return <Map size={16} />;
            default:
                return <MapPin size={16} />;
        }
    };
    return (
        <div 
            className={`${styles.locationSuggestion} ${isActive ? styles.active : ''}`}
            onClick={onClick}
            role="option"
            aria-selected={isActive}
        >
            <div className={styles.locationIcon}>
                {getIcon()}
            </div>
            <div className={styles.locationDetails}>
                <span className={styles.mainText}>{location}</span>
                {distance && (
                    <span className={styles.subText}>
                        {distance} km away
                    </span>
                )}
            </div>
        </div>
    );
});
LocationSuggestion.displayName = 'LocationSuggestion';
export default LocationSuggestion; 