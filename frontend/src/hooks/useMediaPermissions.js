import { useState, useCallback, useEffect } from 'react';
/**
 * Comprehensive permissions hook for media, location, and notifications
 * Handles microphone, camera, geolocation, and notification permissions
 * with proper error handling and user-friendly messages
 */
const useMediaPermissions = () => {
  // Permission states: 'prompt', 'granted', 'denied', 'checking'
  const [micPermission, setMicPermission] = useState('prompt');
  const [cameraPermission, setCameraPermission] = useState('prompt');
  const [locationPermission, setLocationPermission] = useState('prompt');
  const [notificationPermission, setNotificationPermission] = useState('prompt');
  // Error states
  const [permissionErrors, setPermissionErrors] = useState({});
  /**
   * Check all permission states without requesting them
   */
  const checkAllPermissions = useCallback(async () => {
    try {
      // Check microphone permission
      if (navigator.permissions && navigator.permissions.query) {
        try {
          const micResult = await navigator.permissions.query({ name: 'microphone' });
          setMicPermission(micResult.state);
          // Listen for permission changes
          micResult.onchange = () => setMicPermission(micResult.state);
        } catch (e) {
          // Fallback for browsers that don't support permissions API
          setMicPermission('prompt');
        }
        // Check camera permission
        try {
          const cameraResult = await navigator.permissions.query({ name: 'camera' });
          setCameraPermission(cameraResult.state);
          cameraResult.onchange = () => setCameraPermission(cameraResult.state);
        } catch (e) {
          setCameraPermission('prompt');
        }
        // Check geolocation permission
        try {
          const locationResult = await navigator.permissions.query({ name: 'geolocation' });
          setLocationPermission(locationResult.state);
          locationResult.onchange = () => setLocationPermission(locationResult.state);
        } catch (e) {
          setLocationPermission('prompt');
        }
      }
      // Check notification permission (different API)
      if ('Notification' in window) {
        setNotificationPermission(Notification.permission);
      }
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
  }, []);
  // Check initial permission states on mount
  useEffect(() => {
    checkAllPermissions();
  }, [checkAllPermissions]);
  /**
   * Request microphone permission
   */
  const requestMicPermission = useCallback(async () => {
    setMicPermission('checking');
    setPermissionErrors(prev => ({ ...prev, microphone: null }));
    try {
      if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        throw new Error('Microphone access is not supported by your browser');
      }
      const stream = await navigator.mediaDevices.getUserMedia({ 
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true
        }
      });
      setMicPermission('granted');
      // Release the stream immediately as we only needed permission
      stream.getTracks().forEach(track => track.stop());
      return 'granted';
    } catch (err) {
      let errorMessage = 'Failed to access microphone';
      switch (err.name) {
        case 'NotAllowedError':
        case 'PermissionDeniedError':
          setMicPermission('denied');
          errorMessage = 'Microphone permission denied. Please enable microphone access in your browser settings.';
          break;
        case 'NotFoundError':
        case 'DevicesNotFoundError':
          setMicPermission('denied');
          errorMessage = 'No microphone found. Please ensure a microphone is connected.';
          break;
        case 'NotReadableError':
        case 'TrackStartError':
          setMicPermission('denied');
          errorMessage = 'Microphone is already in use by another application.';
          break;
        case 'OverconstrainedError':
        case 'ConstraintNotSatisfiedError':
          setMicPermission('denied');
          errorMessage = 'Microphone constraints could not be satisfied.';
          break;
        case 'NotSupportedError':
          setMicPermission('denied');
          errorMessage = 'Microphone access is not supported by your browser.';
          break;
        case 'AbortError':
          setMicPermission('denied');
          errorMessage = 'Microphone access was aborted.';
          break;
        default:
          setMicPermission('denied');
          errorMessage = err.message || 'Unknown error accessing microphone';
      }
      setPermissionErrors(prev => ({ ...prev, microphone: errorMessage }));
      return 'denied';
    }
  }, []);
  /**
   * Request camera permission
   */
  const requestCameraPermission = useCallback(async () => {
    setCameraPermission('checking');
    setPermissionErrors(prev => ({ ...prev, camera: null }));
    try {
      if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        throw new Error('Camera access is not supported by your browser');
      }
      const stream = await navigator.mediaDevices.getUserMedia({ 
        video: {
          width: { ideal: 1280 },
          height: { ideal: 720 },
          facingMode: 'user'
        }
      });
      setCameraPermission('granted');
      // Release the stream immediately
      stream.getTracks().forEach(track => track.stop());
      return 'granted';
    } catch (err) {
      let errorMessage = 'Failed to access camera';
      switch (err.name) {
        case 'NotAllowedError':
        case 'PermissionDeniedError':
          setCameraPermission('denied');
          errorMessage = 'Camera permission denied. Please enable camera access in your browser settings.';
          break;
        case 'NotFoundError':
        case 'DevicesNotFoundError':
          setCameraPermission('denied');
          errorMessage = 'No camera found. Please ensure a camera is connected.';
          break;
        case 'NotReadableError':
        case 'TrackStartError':
          setCameraPermission('denied');
          errorMessage = 'Camera is already in use by another application.';
          break;
        case 'OverconstrainedError':
        case 'ConstraintNotSatisfiedError':
          setCameraPermission('denied');
          errorMessage = 'Camera constraints could not be satisfied.';
          break;
        case 'NotSupportedError':
          setCameraPermission('denied');
          errorMessage = 'Camera access is not supported by your browser.';
          break;
        case 'AbortError':
          setCameraPermission('denied');
          errorMessage = 'Camera access was aborted.';
          break;
        default:
          setCameraPermission('denied');
          errorMessage = err.message || 'Unknown error accessing camera';
      }
      setPermissionErrors(prev => ({ ...prev, camera: errorMessage }));
      return 'denied';
    }
  }, []);
  /**
   * Request geolocation permission
   */
  const requestLocationPermission = useCallback(async () => {
    setLocationPermission('checking');
    setPermissionErrors(prev => ({ ...prev, location: null }));
    try {
      if (!navigator.geolocation) {
        throw new Error('Geolocation is not supported by your browser');
      }
      return new Promise((resolve, reject) => {
        navigator.geolocation.getCurrentPosition(
          (position) => {
            setLocationPermission('granted');
            resolve('granted');
          },
          (err) => {
            let errorMessage = 'Failed to access location';
            switch (err.code) {
              case err.PERMISSION_DENIED:
                setLocationPermission('denied');
                errorMessage = 'Location permission denied. Please enable location access in your browser settings.';
                break;
              case err.POSITION_UNAVAILABLE:
                setLocationPermission('denied');
                errorMessage = 'Location information is unavailable. Please check your GPS settings.';
                break;
              case err.TIMEOUT:
                setLocationPermission('denied');
                errorMessage = 'Location request timed out. Please try again.';
                break;
              default:
                setLocationPermission('denied');
                errorMessage = err.message || 'Unknown error accessing location';
            }
            setPermissionErrors(prev => ({ ...prev, location: errorMessage }));
            reject('denied');
          },
          {
            enableHighAccuracy: true,
            timeout: 10000,
            maximumAge: 60000
          }
        );
      });
    } catch (err) {
      setLocationPermission('denied');
      const errorMessage = 'Geolocation is not supported by your browser';
      setPermissionErrors(prev => ({ ...prev, location: errorMessage }));
      return 'denied';
    }
  }, []);
  /**
   * Request notification permission
   */
  const requestNotificationPermission = useCallback(async () => {
    if (!('Notification' in window)) {
      const errorMessage = 'Notifications are not supported by your browser';
      setPermissionErrors(prev => ({ ...prev, notification: errorMessage }));
      return 'denied';
    }
    try {
      const permission = await Notification.requestPermission();
      setNotificationPermission(permission);
      if (permission === 'denied') {
        const errorMessage = 'Notification permission denied. Please enable notifications in your browser settings.';
        setPermissionErrors(prev => ({ ...prev, notification: errorMessage }));
      }
      return permission;
    } catch (err) {
      setNotificationPermission('denied');
      const errorMessage = 'Failed to request notification permission';
      setPermissionErrors(prev => ({ ...prev, notification: errorMessage }));
      return 'denied';
    }
  }, []);
  /**
   * Request all permissions at once
   */
  const requestAllPermissions = useCallback(async () => {
    const results = await Promise.allSettled([
      requestMicPermission(),
      requestCameraPermission(),
      requestLocationPermission(),
      requestNotificationPermission()
    ]);
    return {
      microphone: results[0].status === 'fulfilled' ? results[0].value : 'denied',
      camera: results[1].status === 'fulfilled' ? results[1].value : 'denied',
      location: results[2].status === 'fulfilled' ? results[2].value : 'denied',
      notification: results[3].status === 'fulfilled' ? results[3].value : 'denied'
    };
  }, [requestMicPermission, requestCameraPermission, requestLocationPermission, requestNotificationPermission]);
  /**
   * Check if a permission is granted
   */
  const isPermissionGranted = useCallback((permission) => {
    switch (permission) {
      case 'microphone':
        return micPermission === 'granted';
      case 'camera':
        return cameraPermission === 'granted';
      case 'location':
        return locationPermission === 'granted';
      case 'notification':
        return notificationPermission === 'granted';
      default:
        return false;
    }
  }, [micPermission, cameraPermission, locationPermission, notificationPermission]);
  /**
   * Get permission status
   */
  const getPermissionStatus = useCallback((permission) => {
    switch (permission) {
      case 'microphone':
        return micPermission;
      case 'camera':
        return cameraPermission;
      case 'location':
        return locationPermission;
      case 'notification':
        return notificationPermission;
      default:
        return 'prompt';
    }
  }, [micPermission, cameraPermission, locationPermission, notificationPermission]);
  /**
   * Get permission error message
   */
  const getPermissionError = useCallback((permission) => {
    return permissionErrors[permission] || null;
  }, [permissionErrors]);
  /**
   * Clear permission error
   */
  const clearPermissionError = useCallback((permission) => {
    setPermissionErrors(prev => ({ ...prev, [permission]: null }));
  }, []);
  /**
   * Reset all permissions (useful for testing)
   */
  const resetPermissions = useCallback(() => {
    setMicPermission('prompt');
    setCameraPermission('prompt');
    setLocationPermission('prompt');
    setNotificationPermission('prompt');
    setPermissionErrors({});
  }, []);
  return {
    // Permission states
    micPermission,
    cameraPermission,
    locationPermission,
    notificationPermission,
    // Permission request functions
    requestMicPermission,
    requestCameraPermission,
    requestLocationPermission,
    requestNotificationPermission,
    requestAllPermissions,
    // Utility functions
    isPermissionGranted,
    getPermissionStatus,
    getPermissionError,
    clearPermissionError,
    checkAllPermissions,
    resetPermissions,
    // Error states
    permissionErrors,
    // Legacy compatibility (for existing code)
    requestCameraPermission: requestCameraPermission, // alias
    requestMicPermission: requestMicPermission, // alias
  };
};
export default useMediaPermissions; 