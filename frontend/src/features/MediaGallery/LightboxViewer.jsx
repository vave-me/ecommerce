// File: src/components/LightboxViewer.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import Lightbox from 'yet-another-react-lightbox';
import 'yet-another-react-lightbox/styles.css';
/**
 * LightboxViewer:
 *  A small wrapper around 'yet-another-react-lightbox'.
 */
function LightboxViewer({slides, currentIndex, onClose, onNext, onPrev}) {
    return (
        <Lightbox
            open
            close={onClose}
            slides={slides}
            index={currentIndex}
            on={{
                clickNext: onNext,
                clickPrev: onPrev,
                close: onClose
            }}
            styles={{
                container: {
                    backgroundColor: 'rgba(0, 0, 0, 0.9)'
                }
            }}
        />
    );
}
LightboxViewer.propTypes = {
    slides: PropTypes.arrayOf(
        PropTypes.shape({
            src: PropTypes.string.isRequired,
            type: PropTypes.oneOf(['image', 'video']).isRequired,
            alt: PropTypes.string
        })
    ).isRequired,
    currentIndex: PropTypes.number.isRequired,
    onClose: PropTypes.func.isRequired,
    onNext: PropTypes.func.isRequired,
    onPrev: PropTypes.func.isRequired
};
export default memo(LightboxViewer);
