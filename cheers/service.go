package cheers

import "time"

type Cheer struct {
	Value    string
	DateTime time.Time
}

type Service struct {
	cheers []*Cheer
}

func (s *Service) getCheers() []*Cheer {
	return s.cheers
}

func (s *Service) addCheer(cheer *Cheer) {
	s.cheers = append(s.cheers, cheer)
}
