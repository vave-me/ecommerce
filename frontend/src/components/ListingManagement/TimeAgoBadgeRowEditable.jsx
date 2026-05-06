// File: src/components/MobileItemCard/TimeAgoBadgesRowEditable.jsx
"use client"
import React, {useEffect, useState, useRef} from 'react';
import ReactDOM from 'react-dom';
import PropTypes from 'prop-types';
import styled, {keyframes} from 'styled-components';
import { MapPin, Edit } from '@/icons';
import MapView from "../Location/MapView";
/**
 * TimeAgoBadgesRowEditable:
 * - Displays `timeAgo` text with optional inline editing.
 * - Shows a "Nearby" button that opens a map modal, also with optional
 *   inline editing for lat/lon if you choose.
 * - Locks body scroll when the map modal is open; closes on Escape.
 */
function TimeAgoBadgesRowEditable({
                                      timeAgo: initialTimeAgo,
                                      approximateLat: initialLat,
                                      approximateLon: initialLon,
                                  }) {
    const [timeAgo, setTimeAgo] = useState(initialTimeAgo);
    const [editingTime, setEditingTime] = useState(false);
    const [latString, setLatString] = useState(
        initialLat != null ? String(initialLat) : ''
    );
    const [lonString, setLonString] = useState(
        initialLon != null ? String(initialLon) : ''
    );
    const [editingLocation, setEditingLocation] = useState(false);
    const [showMapModal, setShowMapModal] = useState(false);
    // Parse lat/lon to floats
    const latNum = parseFloat(latString);
    const lonNum = parseFloat(lonString);
    const lat = Number.isFinite(latNum) ? latNum : null;
    const lon = Number.isFinite(lonNum) ? lonNum : null;
    /* ---------------- Inline Edit: timeAgo text ---------------- */
    const handleStartEditTime = () => setEditingTime(true);
    const handleCancelEditTime = () => {
        setTimeAgo(initialTimeAgo); // revert
        setEditingTime(false);
    };
    const handleSaveEditTime = () => {
        // For a real app, you might call an API here
        // then confirm the new `timeAgo` is stored
        setEditingTime(false);
    };
    /* ---------------- Inline Edit: location (lat, lon) ---------------- */
    const handleStartEditLocation = () => setEditingLocation(true);
    const handleCancelEditLocation = () => {
        setLatString(initialLat != null ? String(initialLat) : '');
        setLonString(initialLon != null ? String(initialLon) : '');
        setEditingLocation(false);
    };
    const handleSaveEditLocation = () => {
        // Example: call an API or confirm the new lat/lon
        setEditingLocation(false);
    };
    /* ---------------- Map Modal logic ---------------- */
    const handleOpenMap = () => {
        if (!lat || !lon) {
            return;
        }
        setShowMapModal(true);
    };
    const handleCloseMap = () => setShowMapModal(false);
    // Body Scroll Lock
    useEffect(() => {
        if (showMapModal) {
            document.body.style.overflow = 'hidden';
        } else {
            document.body.style.overflow = '';
        }
        return () => {
            document.body.style.overflow = '';
        };
    }, [showMapModal]);
    // Close on Escape
    useEffect(() => {
        const handleEsc = (e) => {
            if (e.key === 'Escape') handleCloseMap();
        };
        if (showMapModal) {
            window.addEventListener('keydown', handleEsc);
        }
        return () => window.removeEventListener('keydown', handleEsc);
    }, [showMapModal]);
    const modalPortalContent =
        showMapModal && lat != null && lon != null ? (
            <Overlay onClick={handleCloseMap}>
                <Modal onClick={(e) => e.stopPropagation()}>
                    <CloseBtn onClick={handleCloseMap} aria-label="Close Modal">
                        &times;
                    </CloseBtn>
                    <MapContainerStyle>
                        <MapView
                            lat={lat}
                            lon={lon}
                            show={showMapModal}
                            zoom={13}
                            radiusMeters={2000}
                        />
                    </MapContainerStyle>
                </Modal>
            </Overlay>
        ) : null;
    return (
        <>
            <Row>
                {/* 1) TimeAgo editable display */}
                <TimeSection>
                    {editingTime ? (
                        <InlineEditRow>
                            <EditInput
                                value={timeAgo}
                                onChange={(e) => setTimeAgo(e.target.value)}
                            />
                            <ButtonRow>
                                <SmallButton onClick={handleSaveEditTime}>Save</SmallButton>
                                <SmallButtonGray onClick={handleCancelEditTime}>
                                    Cancel
                                </SmallButtonGray>
                            </ButtonRow>
                        </InlineEditRow>
                    ) : (
                        <TimeAgoDisplay onClick={handleStartEditTime}>
                            {timeAgo} <Edit className="icon"/>
                        </TimeAgoDisplay>
                    )}
                </TimeSection>
                {/* 2) Location editable display */}
                <LocationSection>
                    {editingLocation ? (
                        <InlineEditRow>
                            <EditInput
                                style={{width: '70px'}}
                                value={latString}
                                onChange={(e) => setLatString(e.target.value)}
                            />
                            <EditInput
                                style={{width: '70px'}}
                                value={lonString}
                                onChange={(e) => setLonString(e.target.value)}
                            />
                            <ButtonRow>
                                <SmallButton onClick={handleSaveEditLocation}>Save</SmallButton>
                                <SmallButtonGray onClick={handleCancelEditLocation}>
                                    Cancel
                                </SmallButtonGray>
                            </ButtonRow>
                        </InlineEditRow>
                    ) : (
                        <LocButton onClick={handleOpenMap}>
                            <MapPin/>
                            <span>Nearby</span>
                        </LocButton>
                    )}
                    {!editingLocation && (
                        <EditLocationIcon onClick={handleStartEditLocation}>
                            <Edit className="icon"/>
                        </EditLocationIcon>
                    )}
                </LocationSection>
            </Row>
            {modalPortalContent &&
                ReactDOM.createPortal(modalPortalContent, document.body)}
        </>
    );
}
TimeAgoBadgesRowEditable.propTypes = {
    /** The displayed "time ago" text (e.g. "2 days ago") */
    timeAgo: PropTypes.string.isRequired,
    /** The approximate latitude or lat string the user can edit (optional) */
    approximateLat: PropTypes.oneOfType([PropTypes.number, PropTypes.string]),
    /** The approximate longitude or lon string the user can edit (optional) */
    approximateLon: PropTypes.oneOfType([PropTypes.number, PropTypes.string]),
};
TimeAgoBadgesRowEditable.defaultProps = {
    approximateLat: null,
    approximateLon: null,
};
export default TimeAgoBadgesRowEditable;
/* ----------------- STYLED COMPONENTS & ANIMATIONS ----------------- */
const fadeIn = keyframes`
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
`;
const slideUp = keyframes`
    from {
        transform: translate(-50%, -60%);
        opacity: 0;
    }
    to {
        transform: translate(-50%, -50%);
        opacity: 1;
    }
`;
const Row = styled.div`
    display: flex;
    align-items: center;
    gap: 16px;
    background: #fff;
    padding: 8px 0;
`;
/* ---------------- Time Section (timeAgo) ---------------- */
const TimeSection = styled.div`
    display: flex;
    align-items: center;
    gap: 8px;
`;
const TimeAgoDisplay = styled.span`
    color: #666;
    font-size: 0.85rem;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    gap: 4px;
    .icon {
        font-size: 0.8rem;
        color: #999;
    }
    &:hover {
        color: #0077cc;
        .icon {
            color: #0077cc;
        }
    }
`;
/* ---------------- Location Section (lat/lon + map) ---------------- */
const LocationSection = styled.div`
    display: flex;
    align-items: center;
    gap: 6px;
    margin-left: auto;
    position: relative;
`;
const LocButton = styled.button`
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: #f8f9fa;
    border: 1px solid #ccc;
    padding: 6px 12px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 500;
    color: #333;
    svg {
        color: #ff6000;
        font-size: 1.1rem;
    }
    &:hover {
        background: #ebecee;
    }
`;
const EditLocationIcon = styled.button`
    background: none;
    border: none;
    padding: 0;
    margin-left: 4px;
    cursor: pointer;
    .icon {
        font-size: 0.85rem;
        color: #999;
    }
    &:hover {
        .icon {
            color: #0077cc;
        }
    }
`;
/* ---------------- Inline Editing (common UI) ---------------- */
const InlineEditRow = styled.div`
    display: flex;
    align-items: center;
    gap: 6px;
`;
const EditInput = styled.input`
    border: 1px solid #aaa;
    border-radius: 4px;
    padding: 4px 8px;
    font-size: 0.85rem;
    &:focus {
        outline: none;
        border-color: #0077cc;
        background-color: #f0f8ff;
    }
`;
const ButtonRow = styled.div`
    display: flex;
    gap: 6px;
`;
const SmallButton = styled.button`
    background-color: #0061bb;
    color: #fff;
    border: none;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
    &:hover {
        background-color: #004c92;
    }
`;
const SmallButtonGray = styled.button`
    background-color: #999;
    color: #fff;
    border: none;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
    &:hover {
        background-color: #777;
    }
`;
/* ---------------- Map Modal & Overlay ---------------- */
const Overlay = styled.div`
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 9999;
    animation: ${fadeIn} 0.3s forwards;
    display: flex;
    justify-content: center;
    align-items: center;
`;
const Modal = styled.div`
    position: fixed;
    top: 50%;
    left: 50%;
    width: 90%;
    max-width: 600px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 5px 18px rgba(0, 0, 0, 0.3);
    transform: translate(-50%, -50%);
    animation: ${slideUp} 0.3s forwards;
    padding: 20px;
    overflow: hidden;
    @media (max-width: 480px) {
        max-width: 90%;
        padding: 16px;
    }
`;
const CloseBtn = styled.button`
    position: absolute;
    top: 10px;
    right: 14px;
    border: none;
    background: none;
    color: #888;
    font-size: 1.5rem;
    cursor: pointer;
    &:hover {
        color: #e63946;
    }
`;
const MapContainerStyle = styled.div`
    width: 100%;
    height: 400px;
    border-radius: 8px;
    overflow: hidden;
    margin-top: 30px;
    @media (max-width: 480px) {
        height: 300px;
    }
`;
