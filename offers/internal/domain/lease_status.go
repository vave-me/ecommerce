package domain // LeaseStatus indicates the stage of a lease

type LeaseStatus string

const (
	LeaseStatusPending     LeaseStatus = "pending"     // e.g. created, not active yet
	LeaseStatusNegotiating LeaseStatus = "negotiating" // new: negotiation in progress
	LeaseStatusActive      LeaseStatus = "active"
	LeaseStatusDeclined    LeaseStatus = "declined" // new: user-seller or user-customer declined the lease
	LeaseStatusRejected    LeaseStatus = "rejected" // new: lessor (or system) rejects lease
	LeaseStatusCompleted   LeaseStatus = "completed"
	LeaseStatusDefaulted   LeaseStatus = "defaulted"
)
