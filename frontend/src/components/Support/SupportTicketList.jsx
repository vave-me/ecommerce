// File: src/components/SupportTicketList.jsx
import React, { memo } from "react";
import PropTypes from "prop-types";
import SupportTicketItem from "./SupportTicketItem";
// Import the CSS module
import styles from "./SupportTicketList.module.css";
/**
 * SupportTicketList Component
 * Displays a list of support tickets.
 */
const SupportTicketList = memo(({ tickets, onUpdate, onDelete }) => {
    return (
        <div className={styles.listContainer}>
            {tickets.map((ticket) => (
                <SupportTicketItem
                    key={ticket.id}
                    ticket={ticket}
                    onUpdate={onUpdate}
                    onDelete={onDelete}
                />
            ))}
        </div>
    );
});
SupportTicketList.displayName = 'SupportTicketList';
SupportTicketList.propTypes = {
    tickets: PropTypes.arrayOf(PropTypes.object).isRequired,
    onUpdate: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
export default SupportTicketList;
