package load

var (
	LoadService LoadServiceInterface = &loadService{}
)

type LoadServiceInterface interface {
	Synca()
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
