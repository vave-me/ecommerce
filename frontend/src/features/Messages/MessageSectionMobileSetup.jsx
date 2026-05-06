// src/components/MobileItemCard/MessageSectionMobileSetup.jsx
import React from 'react';
import PropTypes from 'prop-types';
import MessageSectionMobile from './MessageForm';
import LazyNATSProvider from "../../components/Utils/LazyNATSProvider";
function MessageSectionMobileSetup({
                                       itemId,
                                       onClose,
                                       recipientId,
                                       handleSendMessage = null,
                                       metadata = null,
                                   }) {
    return (
        <LazyNATSProvider>
            <MessageSectionMobile
                itemId={itemId}
                onClose={onClose}
                recipientId={recipientId}
                handleSendMessage={handleSendMessage}
                metadata={metadata}
            />
        </LazyNATSProvider>
    );
}
MessageSectionMobileSetup.propTypes = {
    itemId: PropTypes.string.isRequired,
    recipient: PropTypes.shape({
        id: PropTypes.string.isRequired,
        name: PropTypes.string.isRequired,
        avatar: PropTypes.string,
        online: PropTypes.bool,
    }).isRequired,
    onClose: PropTypes.func.isRequired,
    recipientId: PropTypes.string.isRequired,
    handleSendMessage: PropTypes.func,
    metadata: PropTypes.object,
    isOpen: PropTypes.bool.isRequired,
};
export default MessageSectionMobileSetup;
