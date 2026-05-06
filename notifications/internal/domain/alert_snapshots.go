package domain

type AlertV1 struct {
	UserID    string
	AlertType AlertType
	Message   string
	Payload   map[string]interface{}
	IsRead    bool
}

func (AlertV1) SnapshotName() string { return "notifications.AlertV1" }
