// File: src/components/MapView.client.jsx
"use client";
import React, { memo } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
const MapViewClient = memo(function MapViewClient({lat, lon, zoom = 13}) {
    if (!lat || !lon) {
        return <div>Invalid coordinates</div>;
    }
    return (
        <MapContainer
            center={[lat, lon]}
            zoom={zoom}
            style={{ height: '400px', width: '100%' }}
        >
            <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            />
            <Marker position={[lat, lon]}>
                <Popup>
                    Location: {lat}, {lon}
                </Popup>
            </Marker>
        </MapContainer>
    );
});
export default MapViewClient;
