import React, { useState, useCallback, useRef, useEffect, useMemo } from 'react';
import PropTypes from 'prop-types';
import { MapPin, X, Loader2, Target } from '@/icons';
import { suggestCity, getCurrentLocation, reverseGeocode } from '../../../../api/geocodingApi';
import { debounce } from '../../../../utils/debounce';
import useMediaPermissions from '../../../../hooks/useMediaPermissions';
import styles from './LocationInput.module.css';
/**
 * Shared LocationInput Component
 * Used across all Create modals for consistent location input functionality
 * Supports zipcode input with city suggestions and current location detection
 */
export function LocationInput({
    value = '',
    latitude = null,
    longitude = null,
    onLocationChange,
    onCoordinatesChange,
    placeholder = 'Enter zipcode or city name',
    label = 'Location (Optional)',
    disabled = false,
    error = '',
    className = '',
    showCurrentLocationButton = true,
    maxSuggestions = 5,
    debounceMs = 300,
    required = false,
    'aria-label': ariaLabel
}) {
    // State management
    const [inputValue, setInputValue] = useState(value);
    const [suggestions, setSuggestions] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [showSuggestions, setShowSuggestions] = useState(false);
    const [currentLocationLoading, setCurrentLocationLoading] = useState(false);
    const [activeIndex, setActiveIndex] = useState(-1);
    const [localError, setLocalError] = useState('');
    // Use enhanced permissions hook
    const { 
        locationPermission, 
        requestLocationPermission, 
        getPermissionError, 
        clearPermissionError 
    } = useMediaPermissions();
    // Refs for DOM manipulation
    const inputRef = useRef(null);
    const suggestionsRef = useRef(null);
    const containerRef = useRef(null);
    // Sync input value with prop
    useEffect(() => {
        setInputValue(value);
    }, [value]);
    // Memoized debounced fetch function
    const fetchSuggestions = useMemo(() => 
        debounce(async (query) => {
            if (!query || query.length < 2) {
                setSuggestions([]);
                setShowSuggestions(false);
                return;
            }
            setIsLoading(true);
            setLocalError('');
            try {
                const response = await suggestCity(query);
                const cities = response.suggestedCities || [];
                setSuggestions(cities.slice(0, maxSuggestions));
                setShowSuggestions(cities.length > 0);
            } catch (err) {
                setLocalError('Error fetching suggestions');
                setSuggestions([]);
                setShowSuggestions(false);
            } finally {
                setIsLoading(false);
            }
        }, debounceMs), 
    [maxSuggestions, debounceMs]);
    // Handle input change
    const handleInputChange = useCallback((e) => {
        const newValue = e.target.value;
        setInputValue(newValue);
        setActiveIndex(-1);
        setLocalError('');
        if (newValue) {
            fetchSuggestions(newValue);
        } else {
            setSuggestions([]);
            setShowSuggestions(false);
            onLocationChange?.('');
            onCoordinatesChange?.(null, null);
        }
    }, [fetchSuggestions, onLocationChange, onCoordinatesChange]);
    // Handle suggestion selection
    const handleSuggestionSelect = useCallback((suggestion) => {
        const cityName = suggestion.suggestedCity;
        const lat = suggestion.latitude;
        const lng = suggestion.longitude;
        setInputValue(cityName);
        setShowSuggestions(false);
        setActiveIndex(-1);
        onLocationChange?.(cityName);
        onCoordinatesChange?.(lat, lng);
    }, [onLocationChange, onCoordinatesChange]);
    // Handle current location
    const handleCurrentLocation = useCallback(async () => {
        setCurrentLocationLoading(true);
        setLocalError('');
        clearPermissionError('location');
        // Check and request location permission if needed
        if (locationPermission !== 'granted') {
            try {
                const permission = await requestLocationPermission();
                if (permission !== 'granted') {
                    const permissionError = getPermissionError('location');
                    setLocalError(permissionError || 'Location permission is required');
                    setCurrentLocationLoading(false);
                    return;
                }
            } catch (err) {
                const permissionError = getPermissionError('location');
                setLocalError(permissionError || 'Failed to request location permission');
                setCurrentLocationLoading(false);
                return;
            }
        }
        try {
            const position = await getCurrentLocation({
                enableHighAccuracy: true,
                timeout: 10000,
                maximumAge: 60000
            });
            const { latitude: lat, longitude: lng } = position.coords;
            // Reverse geocode to get city name
            const address = await reverseGeocode(lat, lng);
            const cityName = address.city || address.town || address.village || 
                           `${lat.toFixed(4)}, ${lng.toFixed(4)}`;
            setInputValue(cityName);
            onLocationChange?.(cityName);
            onCoordinatesChange?.(lat, lng);
            setShowSuggestions(false);
        } catch (err) {
            let errorMessage = 'Error getting location';
            switch (err.code) {
                case 1: // PERMISSION_DENIED
                    errorMessage = 'Location permission denied. Please enable location access in your browser settings.';
                    break;
                case 2: // POSITION_UNAVAILABLE
                    errorMessage = 'Location unavailable. Please check your GPS settings.';
                    break;
                case 3: // TIMEOUT
                    errorMessage = 'Location request timed out. Please try again.';
                    break;
                default:
                    errorMessage = err.message || 'Unknown error getting location';
            }
            setLocalError(errorMessage);
        } finally {
            setCurrentLocationLoading(false);
        }
    }, [onLocationChange, onCoordinatesChange, locationPermission, requestLocationPermission, getPermissionError, clearPermissionError]);
    // Handle clear location
    const handleClearLocation = useCallback(() => {
        setInputValue('');
        setLocalError('');
        onLocationChange?.('');
        onCoordinatesChange?.(null, null);
        setShowSuggestions(false);
        setActiveIndex(-1);
    }, [onLocationChange, onCoordinatesChange]);
    // Handle keyboard navigation
    const handleKeyDown = useCallback((e) => {
        if (!showSuggestions || suggestions.length === 0) return;
        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                setActiveIndex(prev => 
                    prev < suggestions.length - 1 ? prev + 1 : prev
                );
                break;
            case 'ArrowUp':
                e.preventDefault();
                setActiveIndex(prev => prev > 0 ? prev - 1 : -1);
                break;
            case 'Enter':
                e.preventDefault();
                if (activeIndex >= 0 && activeIndex < suggestions.length) {
                    handleSuggestionSelect(suggestions[activeIndex]);
                }
                break;
            case 'Escape':
                setShowSuggestions(false);
                setActiveIndex(-1);
                break;
        }
    }, [showSuggestions, suggestions, activeIndex, handleSuggestionSelect]);
    // Handle click outside to close suggestions
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (containerRef.current && !containerRef.current.contains(event.target)) {
                setShowSuggestions(false);
                setActiveIndex(-1);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);
    // Determine if location is set
    const hasLocation = latitude !== null && longitude !== null;
    const displayError = error || localError;
    return (
        <div className={`${styles.locationInputContainer} ${className}`} ref={containerRef}>
            {/* Label */}
            {label && (
                <label 
                    className={`${styles.label} ${required ? styles.required : ''}`}
                    htmlFor="location-input"
                >
                    {label}
                </label>
            )}
            {/* Input Container */}
            <div className={styles.inputContainer}>
                <div className={styles.inputWrapper}>
                    <MapPin className={styles.inputIcon} size={16} />
                    <input
                        ref={inputRef}
                        id="location-input"
                        type="text"
                        className={`${styles.input} ${displayError ? styles.inputError : ''} ${hasLocation ? styles.inputSuccess : ''}`}
                        placeholder={placeholder}
                        value={inputValue}
                        onChange={handleInputChange}
                        onKeyDown={handleKeyDown}
                        disabled={disabled}
                        aria-label={ariaLabel || label}
                        aria-expanded={showSuggestions}
                        aria-haspopup="listbox"
                        aria-autocomplete="list"
                        autoComplete="off"
                    />
                    {isLoading && (
                        <Loader2 className={styles.loadingIcon} size={16} />
                    )}
                    {inputValue && !isLoading && (
                        <button
                            type="button"
                            className={styles.clearButton}
                            onClick={handleClearLocation}
                            aria-label="Clear location"
                        >
                            <X size={14} />
                        </button>
                    )}
                </div>
                {/* Action Buttons */}
                <div className={styles.actionButtons}>
                    {showCurrentLocationButton && (
                        <button
                            type="button"
                            className={`${styles.locationButton} ${hasLocation ? styles.activeLocationButton : ''}`}
                            onClick={handleCurrentLocation}
                            disabled={currentLocationLoading || disabled}
                            aria-label="Use current location"
                        >
                            <Target size={14} />
                            {currentLocationLoading ? 'Getting...' : 'Current'}
                        </button>
                    )}
                    {hasLocation && (
                        <button
                            type="button"
                            className={styles.clearLocationButton}
                            onClick={handleClearLocation}
                            disabled={disabled}
                            aria-label="Clear location"
                        >
                            Clear
                        </button>
                    )}
                </div>
            </div>
            {/* Error Message */}
            {displayError && (
                <div className={styles.errorMessage} role="alert">
                    {displayError}
                </div>
            )}
            {/* Location Confirmation */}
            {hasLocation && !displayError && (
                <div className={styles.locationConfirmation}>
                    <MapPin size={14} />
                    <span>Location set: {inputValue}</span>
                </div>
            )}
            {/* Suggestions Dropdown */}
            {showSuggestions && suggestions.length > 0 && (
                <div 
                    ref={suggestionsRef}
                    className={styles.suggestionsList}
                    role="listbox"
                    aria-label="Location suggestions"
                >
                    {suggestions.map((suggestion, index) => (
                        <div
                            key={`${suggestion.suggestedCity}-${index}`}
                            className={`${styles.suggestionItem} ${index === activeIndex ? styles.suggestionActive : ''}`}
                            onClick={() => handleSuggestionSelect(suggestion)}
                            role="option"
                            aria-selected={index === activeIndex}
                            onMouseEnter={() => setActiveIndex(index)}
                        >
                            <MapPin size={12} className={styles.suggestionIcon} />
                            <div className={styles.suggestionContent}>
                                <span className={styles.suggestionCity}>
                                    {suggestion.suggestedCity}
                                </span>
                                <span className={styles.suggestionCoords}>
                                    {suggestion.latitude.toFixed(4)}, {suggestion.longitude.toFixed(4)}
                                </span>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
LocationInput.propTypes = {
    value: PropTypes.string,
    latitude: PropTypes.number,
    longitude: PropTypes.number,
    onLocationChange: PropTypes.func,
    onCoordinatesChange: PropTypes.func,
    placeholder: PropTypes.string,
    label: PropTypes.string,
    disabled: PropTypes.bool,
    error: PropTypes.string,
    className: PropTypes.string,
    showCurrentLocationButton: PropTypes.bool,
    maxSuggestions: PropTypes.number,
    debounceMs: PropTypes.number,
    required: PropTypes.bool,
    'aria-label': PropTypes.string,
};
export default LocationInput; 