package service

type HealthyRepo interface {
	Ping() error 
}

type HealthyService struct {
	repo HealthyRepo
}

func NewHealthyService(repo HealthyRepo) (*HealthyService) {
	return &HealthyService{repo: repo}
}

func (hs *HealthyService) IsHealthy() (bool, error) {
	err := hs.repo.Ping()
	if (err != nil) {
		return false, err
	}
	return true, nil
}
