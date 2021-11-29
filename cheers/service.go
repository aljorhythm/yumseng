package cheers

import "time"

type Cheer struct {
	Value           string    `json:"value"`
	ClientCreatedAt time.Time `json:"client_created_at"`
}

type Servicer interface {
	GetCheers() []*Cheer
	AddCheer(cheer *Cheer)
}

type service struct {
	cheers []*Cheer
}

func (s *service) GetCheers() []*Cheer {
	return s.cheers
}

func (s *service) AddCheer(cheer *Cheer) {
	s.cheers = append(s.cheers, cheer)
}

func NewService() Servicer {
	service := &service{}
	service.cheers = []*Cheer{}
	return service
}
