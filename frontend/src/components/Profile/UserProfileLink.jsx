"use client";
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import Link from 'next/link';
import { useAuth } from '../../context/AuthContext';
import { getProfileLink } from '../../utils/profileUtils';
/**
 * A component that renders a link to a user's profile
 * 
 * @param {Object} props - Component props
 * @param {string} props.userId - The ID of the user to link to
 * @param {string} props.username - The username to display
 * @param {React.ReactNode} props.children - Optional children to render inside the link
 * @param {string} props.className - Optional CSS class name
 */
const UserProfileLink = memo(({ userId, username, children, className = '' }) => {
    const { user: currentUser } = useAuth();
    const profileUrl = getProfileLink(userId, currentUser?.userId);
    return (
        <Link href={profileUrl} className={className}>
            {children || username}
        </Link>
    );
});
UserProfileLink.displayName = 'UserProfileLink';
UserProfileLink.propTypes = {
    userId: PropTypes.string.isRequired,
    username: PropTypes.string.isRequired,
    children: PropTypes.node,
    className: PropTypes.string,
};
export default UserProfileLink;