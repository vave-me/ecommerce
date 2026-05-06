// src/components/Location/LocationMap.jsx
"use client"
import React, { useEffect, useState, memo } from 'react';
import PropTypes from 'prop-types';
import { MapPin, X } from '@/icons';
import dynamic from 'next/dynamic';
import styles from './LocationMap.module.css';
// Dynamically import map components to avoid SSR issues
const MapContainer = dynamic(() => import('react-leaflet').then(mod => mod.MapContainer), { ssr: false });
const TileLayer = dynamic(() => import('react-leaflet').then(mod => mod.TileLayer), { ssr: false });
const Marker = dynamic(() => import('react-leaflet').then(mod => mod.Marker), { ssr: false });
const Circle = dynamic(() => import('react-leaflet').then(mod => mod.Circle), { ssr: false });
/**
 * LocationMap Component
 * Shows a modal with Leaflet map for displaying location coordinates
 */
const ChangeView = memo(({ center, zoom }) => {
    const { useMap } = require('react-leaflet');
    const map = useMap();
    React.useEffect(() => {
        if (map) {
            map.setView(center, zoom);
        }
    }, [center, zoom, map]);
    return null;
});
ChangeView.displayName = 'ChangeView';
const LocationMap = memo(({ 
    lat = 52.5145, 
    lng = 13.4791, 
    show = false, 
    onClose = () => {}, 
    zoom = 14,
    radius = 2000,
    title = "Location"
}) => {
    const [isClosing, setIsClosing] = useState(false);
    const [isClient, setIsClient] = useState(false);
    const [leafletLoaded, setLeafletLoaded] = useState(false);
    // Initialize Leaflet on client side
    useEffect(() => {
        setIsClient(true);
        // Load Leaflet CSS and initialize icons
        const loadLeaflet = async () => {
            if (typeof window !== 'undefined') {
                // Import Leaflet CSS
                await import('leaflet/dist/leaflet.css');
                // Fix Leaflet default markers
                const L = await import('leaflet');
                delete L.Icon.Default.prototype._getIconUrl;
                L.Icon.Default.mergeOptions({
                    iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
                    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
                    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
                });
                setLeafletLoaded(true);
            }
        };
        loadLeaflet();
    }, []);
    // Handle close with animation
    const handleClose = () => {
        setIsClosing(true);
        setTimeout(() => {
            setIsClosing(false);
            onClose();
        }, 300); // Match animation duration
    };
    // Handle backdrop click
    const handleBackdropClick = (e) => {
        if (e.target === e.currentTarget) {
            handleClose();
        }
    };
    // Handle escape key
    useEffect(() => {
        const handleEscape = (e) => {
            if (e.key === 'Escape' && show) {
                handleClose();
            }
        };
        if (show) {
            document.addEventListener('keydown', handleEscape);
            // Prevent body scroll when modal is open
            document.body.style.overflow = 'hidden';
        }
        return () => {
            document.removeEventListener('keydown', handleEscape);
            document.body.style.overflow = 'unset';
        };
    }, [show]);
    // Don't render if not shown
    if (!show || !isClient) {
        return null;
    }
    return (
        <>
            {/* Backdrop */}
            <div 
                className={`${styles.backdrop} ${show ? styles.visible : styles.hidden}`}
                onClick={handleBackdropClick}
            />
            {/* Modal */}
            <div className={`${styles.modal} ${show ? styles.visible : styles.hidden} ${isClosing ? styles.closing : ''}`}>
                {/* Header */}
                <div className={styles.header}>
                    <h2 className={styles.title}>
                        {title} ({lat.toFixed(4)}, {lng.toFixed(4)})
                    </h2>
                    <button 
                        className={styles.closeButton} 
                        onClick={handleClose} 
                        aria-label="Close Map"
                        type="button"
                    >
                        <X size={24} />
                    </button>
                </div>
                {/* Map */}
                <div className={styles.mapWrapper}>
                    {leafletLoaded ? (
                        <MapContainer
                            center={[lat, lng]}
                            zoom={zoom}
                            scrollWheelZoom={true}
                            style={{ width: '100%', height: '400px', minHeight: '400px' }}
                            key={`${lat}-${lng}`} // Force re-render when coordinates change
                        >
                            <ChangeView
                                center={[lat, lng]}
                                zoom={zoom}
                            />
                            <TileLayer
                                attribution='&copy; <a href="https://osm.org/copyright">OpenStreetMap</a>'
                                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                            />
                            <Marker position={[lat, lng]} />
                            {radius > 0 && (
                                <Circle
                                    center={[lat, lng]}
                                    radius={radius}
                                    pathOptions={{
                                        color: 'rgba(74, 144, 226, 0.8)',
                                        fillColor: 'rgba(74, 144, 226, 0.4)',
                                        fillOpacity: 0.4,
                                    }}
                                />
                            )}
                        </MapContainer>
                    ) : (
                        <div style={{ 
                            width: '100%', 
                            height: '400px', 
                            display: 'flex', 
                            alignItems: 'center', 
                            justifyContent: 'center',
                            background: '#f5f5f5',
                            color: '#666'
                        }}>
                            Loading map...
                        </div>
                    )}
                </div>
                {/* Footer */}
                <div className={styles.footer}>
                    <button 
                        className={styles.actionButton} 
                        onClick={handleClose}
                        type="button"
                    >
                        Close Map
                    </button>
                </div>
            </div>
        </>
    );
});
LocationMap.displayName = 'LocationMap';
LocationMap.propTypes = {
    lat: PropTypes.number,
    lng: PropTypes.number,
    show: PropTypes.bool,
    onClose: PropTypes.func,
    zoom: PropTypes.number,
    radius: PropTypes.number,
    title: PropTypes.string,
};
LocationMap.defaultProps = {
    lat: 52.5145,
    lng: 13.4791,
    show: false,
    onClose: () => {},
    zoom: 14,
    radius: 2000,
    title: "Location",
};
export default LocationMap;
