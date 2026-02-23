package api

// User represents a Backlog user.
type User struct {
	ID          int    `json:"id"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	MailAddress string `json:"mailAddress"`
}

// Project represents a Backlog project.
type Project struct {
	ID          int    `json:"id"`
	ProjectKey  string `json:"projectKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Issue represents a Backlog issue.
type Issue struct {
	ID          int        `json:"id"`
	ProjectID   int        `json:"projectId"`
	IssueKey    string     `json:"issueKey"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Status      *Status    `json:"status"`
	Assignee    *User      `json:"assignee"`
	Priority    *Priority  `json:"priority"`
	IssueType   *IssueType `json:"issueType"`
	DueDate     string     `json:"dueDate"`
	StartDate   string     `json:"startDate"`
	CreatedUser *User      `json:"createdUser"`
	Created     string     `json:"created"`
	Updated     string     `json:"updated"`
	Milestone   []Milestone `json:"milestone"`
	Category    []Category  `json:"category"`
}

// IssueType represents a Backlog issue type.
type IssueType struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"projectId"`
	Name      string `json:"name"`
	Color     string `json:"color"`
}

// Status represents a Backlog status.
type Status struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"projectId"`
	Name      string `json:"name"`
	Color     string `json:"color"`
}

// Priority represents a Backlog priority.
type Priority struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Comment represents a Backlog comment.
type Comment struct {
	ID          int    `json:"id"`
	Content     string `json:"content"`
	CreatedUser *User  `json:"createdUser"`
	Created     string `json:"created"`
}

// Milestone represents a Backlog milestone/version.
type Milestone struct {
	ID             int    `json:"id"`
	ProjectID      int    `json:"projectId"`
	Name           string `json:"name"`
	ReleaseDueDate string `json:"releaseDueDate"`
}

// Category represents a Backlog category.
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// BacklogError represents an error response from the Backlog API.
type BacklogError struct {
	Message  string `json:"message"`
	Code     int    `json:"code"`
	MoreInfo string `json:"moreInfo"`
}

// BacklogErrorResponse wraps the errors array in the API response.
type BacklogErrorResponse struct {
	Errors []BacklogError `json:"errors"`
}
