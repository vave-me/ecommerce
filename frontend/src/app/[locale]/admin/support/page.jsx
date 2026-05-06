import SupportTickets from './SupportTickets.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Support Tickets - Admin',
  description: 'Manage customer support tickets and inquiries',
};

export default function SupportTicketsPage() {
  return <SupportTickets />;
}