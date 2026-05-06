package domain

type UserRole string

const (
	UserRoleUndefined  UserRole = ""
	UserRoleUser       UserRole = "user"       // Default role for regular users
	UserRoleCustomer   UserRole = "customer"   // Legacy - same as user
	UserRoleAdmin      UserRole = "admin"      // System administrators
	UserRoleSuperadmin UserRole = "superadmin" // Super administrators
	UserRolePublicist  UserRole = "publicist"  // Content creators
	UserRoleBusiness   UserRole = "business"   // Business users
	UserRoleVendor     UserRole = "vendor"     // Marketplace vendors
	UserRoleSupport    UserRole = "support"    // Support staff
	UserRoleSystem     UserRole = "system"     // System/automated users
)

func (s UserRole) String() string {
	switch s {
	case UserRoleUser, UserRoleCustomer, UserRoleAdmin, UserRoleSuperadmin, 
	     UserRolePublicist, UserRoleBusiness, UserRoleVendor, UserRoleSupport, UserRoleSystem:
		return string(s)
	default:
		return string(UserRoleUser) // Default to user role
	}
}

func ToUserRole(status string) UserRole {
	switch status {
	case UserRoleUser.String():
		return UserRoleUser
	case UserRoleCustomer.String():
		return UserRoleCustomer
	case UserRoleAdmin.String():
		return UserRoleAdmin
	case UserRoleSuperadmin.String():
		return UserRoleSuperadmin
	case UserRolePublicist.String():
		return UserRolePublicist
	case UserRoleBusiness.String():
		return UserRoleBusiness
	case UserRoleVendor.String():
		return UserRoleVendor
	case UserRoleSupport.String():
		return UserRoleSupport
	case UserRoleSystem.String():
		return UserRoleSystem
	default:
		return UserRoleUser // Default to user role
	}
}
