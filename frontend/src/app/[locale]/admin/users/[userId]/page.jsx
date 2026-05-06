import UserDetail from './UserDetail.client';

export const metadata = {
  title: 'User Details - Admin',
  description: 'View user details and information',
};

export default function UserDetailPage({ params: { locale, userId } }) {
  return <UserDetail locale={locale} userId={userId} />;
}