import React, { useState, useEffect, useRef, memo } from 'react';
import { MapPin, Globe, ExternalLink, Building2, User, Briefcase } from 'lucide-react';
import sharedStyles from './CardShared.module.css';
// Simple Leaflet Map Modal Component
const LeafletMapModal = memo(({ lat, lng, onClose }) => {
    const [leafletLoaded, setLeafletLoaded] = useState(false);
    const [mapInstance, setMapInstance] = useState(null);
    const mapContainerRef = useRef(null);
    const mapIdRef = useRef(`leaflet-map-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`);
    useEffect(() => {
        let mounted = true;
        let currentMapInstance = null;
        const loadLeaflet = async () => {
            if (typeof window !== 'undefined' && mounted) {
                try {
                    // Load Leaflet CSS
                    const link = document.createElement('link');
                    link.rel = 'stylesheet';
                    link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css';
                    if (!document.querySelector(`link[href="${link.href}"]`)) {
                        document.head.appendChild(link);
                    }
                    // Load Leaflet JS
                    const L = await import('leaflet');
                    // Fix default markers
                    delete L.Icon.Default.prototype._getIconUrl;
                    L.Icon.Default.mergeOptions({
                        iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
                        iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
                        shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
                    });
                    if (mounted) {
                        setLeafletLoaded(true);
                        // Initialize map with a small delay to ensure DOM is ready
                        setTimeout(() => {
                            if (!mounted) return;
                            const mapElement = document.getElementById(mapIdRef.current);
                            if (mapElement && !currentMapInstance) {
                                try {
                                    // Clear any existing map container content
                                    mapElement.innerHTML = '';
                                    // Create new map instance
                                    currentMapInstance = L.map(mapIdRef.current, {
                                        zoomControl: true,
                                        scrollWheelZoom: true,
                                        doubleClickZoom: true,
                                        touchZoom: true
                                    }).setView([lat, lng], 15);
                                    // Add tile layer
                                    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                                        attribution: '© OpenStreetMap contributors',
                                        maxZoom: 19
                                    }).addTo(currentMapInstance);
                                    // Add marker
                                    L.marker([lat, lng]).addTo(currentMapInstance);
                                    // Add blue circle area around the location
                                    L.circle([lat, lng], {
                                        color: '#4a90e2',
                                        fillColor: '#4a90e2',
                                        fillOpacity: 0.2,
                                        radius: 500 // 500 meters radius
                                    }).addTo(currentMapInstance);
                                    if (mounted) {
                                        setMapInstance(currentMapInstance);
                                    }
                                } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
                            }
                        }, 150);
                    }
                } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
            }
        };
        loadLeaflet();
        return () => {
            mounted = false;
            if (currentMapInstance) {
                try {
                    currentMapInstance.remove();
                    currentMapInstance = null;
                } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
            }
            setMapInstance(null);
        };
    }, [lat, lng]);
    const handleBackdropClick = (e) => {
        if (e.target === e.currentTarget) {
            onClose();
        }
    };
    return (
        <div 
            style={{
                position: 'fixed',
                top: 0,
                left: 0,
                width: '100%',
                height: '100%',
                backgroundColor: 'rgba(0, 0, 0, 0.5)',
                zIndex: 9999,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
            }}
            onClick={handleBackdropClick}
        >
            <div 
                style={{
                    backgroundColor: 'white',
                    borderRadius: '12px',
                    padding: window.innerWidth <= 768 ? '15px' : '20px',
                    width: window.innerWidth <= 768 ? '95%' : '90%',
                    maxWidth: window.innerWidth <= 768 ? '600px' : '800px',
                    maxHeight: window.innerWidth <= 768 ? '80vh' : '80%',
                    position: 'relative',
                    boxShadow: '0 10px 25px rgba(0, 0, 0, 0.3)'
                }}
                onClick={(e) => e.stopPropagation()}
            >
                <button 
                    onClick={onClose}
                    style={{
                        position: 'absolute',
                        top: '15px',
                        right: '15px',
                        background: 'rgba(255, 255, 255, 0.9)',
                        border: 'none',
                        borderRadius: '50%',
                        fontSize: '20px',
                        cursor: 'pointer',
                        width: '32px',
                        height: '32px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: '#666',
                        zIndex: 1000,
                        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.2)'
                    }}
                >
                    ×
                </button>
                <div 
                    ref={mapContainerRef}
                    id={mapIdRef.current}
                    style={{ 
                        width: '100%', 
                        height: window.innerWidth <= 768 ? '400px' : '450px',
                        minHeight: '350px',
                        borderRadius: '12px',
                        backgroundColor: '#f5f5f5'
                    }}
                >
                    {!leafletLoaded && (
                        <div style={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            height: '100%',
                            color: '#666'
                        }}>
                            Loading map...
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
});
LeafletMapModal.displayName = 'LeafletMapModal';
/**
 * PriceLocationRow - Smart card type-aware price and location display
 * 
 * Automatically handles different display patterns based on card type:
 * - jobs: salary/undisclosed + company name + location/remote
 * - deals: deal price + regular price + merchant link + location/online  
 * - classified: price + offer type/user type + user name + location/online
 * - properties: price + merchant name + location/online
 * - vehicles: price + merchant name + location/online
 * 
 * @param {Object} props - Component props
 * @param {string} props.cardType - Type of card: 'jobs', 'deals', 'classified', 'properties', 'vehicles'
 * @param {string} props.price - Main price/salary (required)
 * @param {string} props.originalPrice - Original/regular price (for deals)
 * @param {string} props.priceLabel - Accessible label for price
 * @param {string} props.location - Location string
 * @param {number} props.lat - Latitude coordinate
 * @param {number} props.lng - Longitude coordinate
 * @param {string} props.companyName - Company/merchant/business name
 * @param {string} props.userName - User name (for classified)
 * @param {string} props.userType - 'private' or 'business' (for classified)
 * @param {string} props.offerType - 'sell', 'buy', 'rent' etc (for classified)
 * @param {string} props.dealUrl - External URL (for deals)
 * @param {string} props.serviceUrl - External URL (for services)
 * @param {boolean} props.isRemote - Whether job is remote (for jobs)
 * @param {boolean} props.salaryUndisclosed - Whether salary is undisclosed (for jobs)
 * @param {string} props.onlineText - Fallback text when no location
 * @param {string} props.className - Additional CSS classes
 * @param {React.ReactNode} props.additionalContent - Extra content (warranties, shipping, etc.)
 */
const PriceLocationRow = memo(({
    // Card type identification
    cardType = 'classified', // 'jobs', 'deals', 'classified', 'properties', 'vehicles'
    // Price information
    price = '',
    originalPrice = '', // For deals with regular price
    priceLabel = '',
    // Location information
    location = null,
    lat = null,
    lng = null,
    // Entity information (company/merchant/user)
    companyName = '', // For jobs, deals, properties, vehicles
    userName = '', // For classified private listings
    userType = 'private', // 'private' or 'business' for classified
    // Card type specific fields
    offerType = '', // For classified: 'sell', 'buy', 'rent', etc.
    dealUrl = '', // For deals external link
    serviceUrl = '', // For services external link
    isRemote = false, // For jobs
    salaryUndisclosed = false, // For jobs
    // Generic options
    onlineText = 'Online',
    className = '',
    additionalContent = null,
    showLocation = true, // New prop to control location display
    showMerchant = true, // New prop to control merchant/user display
}) => {
    const [showMapModal, setShowMapModal] = useState(false);
    // Parse coordinates from location string if it looks like "lat, lng"
    const parseCoordinates = (locationStr) => {
        if (!locationStr || typeof locationStr !== 'string') return null;
        const parts = locationStr.trim().split(',').map(part => parseFloat(part.trim()));
        if (parts.length === 2 && !isNaN(parts[0]) && !isNaN(parts[1])) {
            return { lat: parts[0], lng: parts[1] };
        }
        return null;
    };
    // Use provided lat/lng props first, then try parsing from location string
    const coordinates = (lat !== null && lng !== null) 
        ? { lat, lng } 
        : parseCoordinates(location);
    const isCoordinateLocation = coordinates !== null;
    const hasLocation = location && location.trim();
    // Get URL reference
    const urlReference = dealUrl || serviceUrl;
    // Card type-specific business logic
    const getCardTypeConfig = () => {
        switch (cardType) {
            case 'jobs':
                return {
                    // Price: salary or "undisclosed"
                    displayPrice: salaryUndisclosed ? 'Undisclosed' : price,
                    priceAriaLabel: salaryUndisclosed ? 'Salary: Undisclosed' : (priceLabel || `Salary: ${price}`),
                    // Middle: company name
                    middleContent: companyName || 'Company',
                    middleIcon: <Building2 size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    middleClickable: false,
                    middleAriaLabel: `Company: ${companyName || 'Company'}`,
                    // Right: location or "Remote"
                    rightContent: isRemote ? 'Remote' : (hasLocation ? location.trim() : 'Remote'),
                    rightIcon: isRemote ? <Globe size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" /> : 
                              <MapPin size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    rightClickable: !isRemote && isCoordinateLocation,
                    rightAriaLabel: isRemote ? 'Work arrangement: Remote' : 
                                   (isCoordinateLocation ? `View location on map: ${location}` : `Location: ${location || 'Remote'}`)
                };
            case 'deals':
                return {
                    // Price: deal price
                    displayPrice: price,
                    priceAriaLabel: priceLabel || `Deal price: ${price}`,
                    // Middle: merchant name with link
                    middleContent: companyName || 'Merchant',
                    middleIcon: <ExternalLink size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    middleClickable: !!urlReference,
                    middleAriaLabel: urlReference ? `Visit ${companyName}` : `Merchant: ${companyName || 'Merchant'}`,
                    // Right: "Show on Map" button or location name or "Online"
                    rightContent: hasLocation ? (isCoordinateLocation ? 'Show on Map' : location.trim()) : onlineText,
                    rightIcon: hasLocation ? <MapPin size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" /> :
                              <Globe size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    rightClickable: hasLocation && isCoordinateLocation,
                    rightAriaLabel: hasLocation && isCoordinateLocation ? `Show ${location} on map` : 
                                   (hasLocation ? `Location: ${location}` : `Available ${onlineText.toLowerCase()}`)
                };
            case 'classified':
                return {
                    // Price: listing price
                    displayPrice: price,
                    priceAriaLabel: priceLabel || `Price: ${price}`,
                    // Middle: offer type / user type + user name
                    middleContent: userType === 'business' ? (companyName || userName || 'Business') : 
                                  (offerType ? `${offerType} • ${userName || 'Private'}` : (userName || 'Private')),
                    middleIcon: userType === 'business' ? 
                               <Building2 size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" /> :
                               <User size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    middleClickable: false,
                    middleAriaLabel: userType === 'business' ? `Business: ${companyName || userName || 'Business'}` :
                                   `Private seller: ${userName || 'Private'}`,
                    // Right: "Show on Map" button or location name or "Online"
                    rightContent: hasLocation ? (isCoordinateLocation ? 'Show on Map' : location.trim()) : onlineText,
                    rightIcon: hasLocation ? <MapPin size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" /> :
                              <Globe size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    rightClickable: hasLocation && isCoordinateLocation,
                    rightAriaLabel: hasLocation && isCoordinateLocation ? `Show ${location} on map` : 
                                   (hasLocation ? `Location: ${location}` : `Available ${onlineText.toLowerCase()}`)
                };
            case 'properties':
                return {
                    // Price: listing price / offer type
                    displayPrice: offerType ? `${price} / ${offerType}` : price,
                    priceAriaLabel: offerType ? `${price} for ${offerType}` : (priceLabel || `Price: ${price}`),
                    // Middle: merchant/agent name with link
                    middleContent: companyName || userName || 'Agent',
                    middleIcon: companyName ? 
                               <Building2 size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" /> :
                               <User size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    middleClickable: !!urlReference,
                    middleAriaLabel: urlReference ? `Visit ${companyName || userName || 'agent'}` : 
                                   `Agent: ${companyName || userName || 'Agent'}`,
                    // Right: "Show on Map" button or location name or "Online"
                    rightContent: hasLocation ? (isCoordinateLocation ? 'Show on Map' : location.trim()) : onlineText,
                    rightIcon: hasLocation ? <MapPin size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" /> :
                              <Globe size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    rightClickable: hasLocation && isCoordinateLocation,
                    rightAriaLabel: hasLocation && isCoordinateLocation ? `Show ${location} on map` : 
                                   (hasLocation ? `Location: ${location}` : `Available ${onlineText.toLowerCase()}`)
                };
            case 'vehicles':
            default:
                return {
                    // Price: vehicle price
                    displayPrice: price,
                    priceAriaLabel: priceLabel || `Price: ${price}`,
                    // Middle: dealer/seller name
                    middleContent: companyName || userName || 'Seller',
                    middleIcon: companyName ? 
                               <Building2 size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" /> :
                               <User size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    middleClickable: !!urlReference,
                    middleAriaLabel: urlReference ? `Visit ${companyName || userName || 'seller'}` : 
                                   `Seller: ${companyName || userName || 'Seller'}`,
                    // Right: "Show on Map" button or location name or "Online"
                    rightContent: hasLocation ? (isCoordinateLocation ? 'Show on Map' : location.trim()) : onlineText,
                    rightIcon: hasLocation ? <MapPin size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" /> :
                              <Globe size={14} className={sharedStyles.sharedLocationIcon} aria-hidden="true" />,
                    rightClickable: hasLocation && isCoordinateLocation,
                    rightAriaLabel: hasLocation && isCoordinateLocation ? `Show ${location} on map` : 
                                   (hasLocation ? `Location: ${location}` : `Available ${onlineText.toLowerCase()}`)
                };
        }
    };
    const config = getCardTypeConfig();
    // Event handlers
    const handleMiddleClick = () => {
        if (config.middleClickable && urlReference) {
            window.open(urlReference, '_blank', 'noopener,noreferrer');
        }
    };
    const handleRightClick = () => {
        if (config.rightClickable && isCoordinateLocation) {
            setShowMapModal(true);
                }
    };
    return (
        <>
            <div className={`${sharedStyles.sharedPriceLocationRow} ${className}`}>
                {/* Price Section */}
                {config.displayPrice && (
                    <div className={sharedStyles.priceSection}>
                        <span className={sharedStyles.sharedPrice} aria-label={config.priceAriaLabel}>
                            {config.displayPrice}
                        </span>
                        {/* Original Price for Deals */}
                        {cardType === 'deals' && originalPrice && (
                            <span className={sharedStyles.originalPrice} aria-label={`Original price: ${originalPrice}`}>
                                {originalPrice}
                    </span>
                )}
                        {/* Additional content (warranties, shipping, etc.) */}
                {additionalContent}
                    </div>
                )}
                {/* Middle Section - Entity (Company/User/Merchant) */}
                {showMerchant && (
                    <div 
                        className={`${sharedStyles.sharedLocationSection} ${config.middleClickable ? sharedStyles.clickable : ''}`}
                        aria-label={config.middleAriaLabel}
                        onClick={config.middleClickable ? handleMiddleClick : undefined}
                        style={{ cursor: config.middleClickable ? 'pointer' : 'default' }}
                    >
                        {config.middleIcon}
                        <span className={sharedStyles.sharedLocationText}>{config.middleContent}</span>
                    </div>
                )}
                {/* Right Section - Location */}
                {showLocation && (
                    <div 
                        className={`${sharedStyles.sharedLocationSection} ${config.rightClickable ? sharedStyles.clickable : ''}`}
                        aria-label={config.rightAriaLabel}
                        onClick={config.rightClickable ? handleRightClick : undefined}
                        style={{ cursor: config.rightClickable ? 'pointer' : 'default' }}
                    >
                        {config.rightIcon}
                        <span className={sharedStyles.sharedLocationText}>{config.rightContent}</span>
                        </div>
                )}
            </div>
            {/* Leaflet Map Modal */}
            {showMapModal && coordinates && (
                <LeafletMapModal 
                    lat={coordinates.lat}
                    lng={coordinates.lng}
                    onClose={() => setShowMapModal(false)}
                />
            )}
        </>
    );
});
PriceLocationRow.displayName = 'PriceLocationRow';
export default PriceLocationRow;