package systemStatus

const (
	wssConnCount = "ws_conn_count"
)

func (s *Service) AddWebSocketConnecionCount() {
	value := s.App.Cache.Get(wssConnCount)
	s.App.Cache.Set(wssConnCount, s.getInt(value)+1)
}

func (s *Service) RemoveWebSocketConnecionCount() {
	value := s.App.Cache.Get(wssConnCount)
	s.App.Cache.Set(wssConnCount, s.getInt(value)-1)
}

func (s *Service) GetWebSocketConnecionCount() int {
	value := s.App.Cache.Get(wssConnCount)
	return s.getInt(value)
}

func (s *Service) getInt(value any) int {
	if value == nil {
		return 0
	}
	return value.(int)
}
