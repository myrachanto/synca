package load

import httperors "github.com/myrachanto/erroring"

var (
	LoadService LoadServiceInterface = &loadService{}
)

type LoadServiceInterface interface {
	Synca()
	GetAll() ([]*Synca, httperors.HttpErr)
}
type loadService struct {
	repo LoadRepoInterface
}

func NewloadService(repository LoadRepoInterface) LoadServiceInterface {
	return &loadService{
		repository,
	}
}
func (service *loadService) Synca() {
	service.repo.Synca()
}
func (service *loadService) GetAll() ([]*Synca, httperors.HttpErr) {
	return service.repo.GetAll()
}
