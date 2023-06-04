package models

type CostLogData struct {
	Id        int    `json:id`
	GroupId   string `json:group_id`
	UserId    string `json:user_id`
	Name      string `json:name`
	Comment   string `json:comment`
	Amount    int    `json:amount`
	CreatedAt string `json:created_at`
	UpdatedAt string `json:updated_at`
}
