// src/components/MapView.jsx
"use client"
import React, {useRef, useEffect, memo} from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
import {MapContainer, TileLayer, Circle} from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
/**
 * A refined map component to display a ~2km radius around a location.
 * Only shows if `show` is true, calling invalidateSize() to avoid tile "chop."
 * Memoized for performance optimization
 *
 * If you want to fix Leaflet icons, do that once in your top-level with initLeafletIcons().
 */
const MapView = memo(function MapView({lat, lon, show, zoom, radiusMeters}) {
    const mapRef = useRef(null);
    // Invalidate size whenever `show` transitions to true
    useEffect(() => {
        if (show && mapRef.current) {
            mapRef.current.invalidateSize();
        }
    }, [show]);
    // If the map is hidden (modal closed), skip rendering
    if (!show) return null;
    // If lat/lon are invalid, show fallback
    if (lat === null || lat === undefined || lon === null || lon === undefined) {
        return <div>Invalid location</div>;
    }
    return (
        <StyledMapContainer
            center={[lat, lon]}
            zoom={zoom}
            scrollWheelZoom
            whenCreated={(mapInstance) => {
                mapRef.current = mapInstance;
            }}
        >
            <TileLayer
                attribution='&copy; OpenStreetMap contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            {/* Circle to show the approximate region */}
            {radiusMeters > 0 && (
                <Circle
                    center={[lat, lon]}
                    radius={radiusMeters}
                    pathOptions={{
                        color: '#ff6000', // stroke
                        weight: 2,
                        fillColor: '#ffa552', // fill
                        fillOpacity: 0.25,
                    }}
                />
            )}
            {/*
        If you want a marker in the center:
        <Marker position={[lat, lon]}>
          <Popup>
            Approx. center
          </Popup>
        </Marker>
      */}
        </StyledMapContainer>
    );
});
MapView.propTypes = {
    lat: PropTypes.number,
    lon: PropTypes.number,
    show: PropTypes.bool,
    zoom: PropTypes.number,
    radiusMeters: PropTypes.number,
};
MapView.defaultProps = {
    lat: null,
    lon: null,
    show: true,
    zoom: 13,
    radiusMeters: 2000, // ~2km
};
export default MapView;
/* --------------- STYLES --------------- */
const StyledMapContainer = styled(MapContainer)`
    width: 100%;
    height: 100%;
    border-radius: 8px;
    overflow: hidden;
`;
