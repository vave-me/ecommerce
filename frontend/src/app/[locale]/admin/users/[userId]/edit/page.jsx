import UserEdit from './UserEdit.client';

export const metadata = {
  title: 'Edit User - Admin',
  description: 'Edit user information and settings',
};

export default function UserEditPage({ params: { locale, userId } }) {
  return <UserEdit locale={locale} userId={userId} />;
}