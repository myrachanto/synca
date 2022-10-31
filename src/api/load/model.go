package load

import (
	"time"

	httperrors "github.com/myrachanto/erroring"
	"gorm.io/gorm"
)

type Synca struct {
	Name      string    `json:"name,omitempty"`
	DatabaseA string    `json:"databaseA,omitempty"`
	DatabaseB string    `json:"databaseB,omitempty"`
	Dated     time.Time `json:"dated,omitempty"`
	Message   string    `json:"message,omitempty"`
	Status    bool      `json:"status,omitempty"`
	Items     int       `json:"items,omitempty"`
	gorm.Model
}
type Result struct {
	Name  string    `json:"name,omitempty"`
	Dated time.Time `json:"dated,omitempty"`
}

func (l Synca) Validate() httperrors.HttpErr {
	if l.Name == "" {
		return httperrors.NewBadRequestError("Url must not be empty")
	}
	return nil
}
