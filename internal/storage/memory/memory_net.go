package memory

import (
	"time"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// SaveTalkersNet Сохраняем текущую статистику по трафику сети
func (s *Store) SaveTalkersNet(data *model.TalkersNet) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.dbTalkersNet = append(s.dbTalkersNet, *data)
	return nil
}

// GetAvgTalkersNet Возврат среднего значения трафика по сетевым интерфейсам за period
func (s *Store) GetAvgTalkersNet(period int32) (*model.TalkersNet, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	return avgTalkersNet(s.dbTalkersNet, period)
}

//avgTalkersNet получить среднее значение показателей за период
func avgTalkersNet(s []model.TalkersNet, period int32) (*model.TalkersNet, error) {
	res := model.TalkersNet{}
	sumDevNet := map[string]*model.DeviceNet{}

	now := time.Now().Local()
	timeStart := now.Add(-time.Second * time.Duration(period))

	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		if timeStart.Before(s[i].QueryTime) {
			for _, dv := range s[i].DevNet {
				_, ok := sumDevNet[dv.NetInterface]
				if !ok {
					sumDevNet[dv.NetInterface] = &model.DeviceNet{}
				}
				sumDevNet[dv.NetInterface].NetInterface = dv.NetInterface
				sumDevNet[dv.NetInterface].Receive.Bytes += dv.Receive.Bytes
				sumDevNet[dv.NetInterface].Receive.Packets += dv.Receive.Packets
				sumDevNet[dv.NetInterface].Receive.Errs += dv.Receive.Errs
				sumDevNet[dv.NetInterface].Transmit.Bytes += dv.Transmit.Bytes
				sumDevNet[dv.NetInterface].Transmit.Packets += dv.Transmit.Packets
				sumDevNet[dv.NetInterface].Receive.Errs = dv.Transmit.Errs
			}
			count++
		} else {
			break
		}
	}
	if count > 1 {
		// расчитываем средние значения по каждому интерфейсу
		for _, v := range sumDevNet {
			sumDevNet[v.NetInterface].Receive.Bytes = sumDevNet[v.NetInterface].Receive.Bytes / count
			sumDevNet[v.NetInterface].Receive.Packets = sumDevNet[v.NetInterface].Receive.Packets / count
			sumDevNet[v.NetInterface].Receive.Errs = sumDevNet[v.NetInterface].Receive.Errs / count
			sumDevNet[v.NetInterface].Transmit.Bytes = sumDevNet[v.NetInterface].Transmit.Bytes / count
			sumDevNet[v.NetInterface].Transmit.Packets = sumDevNet[v.NetInterface].Transmit.Packets / count
			sumDevNet[v.NetInterface].Transmit.Errs = sumDevNet[v.NetInterface].Transmit.Errs / count
		}
		for _, v := range sumDevNet {
			res.DevNet = append(res.DevNet, *v)
		}
	} else {
		// берём последнее значение
		if len(s) > 0 {
			res = s[len(s)-1]
		}
	}
	res.QueryTime = now

	return &res, nil
}

// SaveTalkersNet Сохраняем текущую статистику по трафику сети
func (s *Store) SaveNetworkStatistics(data *model.NetworkStatistics) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.dbNetStat = append(s.dbNetStat, *data)
	return nil
}

// GetAvgTalkersNet Возврат среднего значения трафика по сетевым интерфейсам за period
func (s *Store) GetAvgNetworkStatistics(period int32) (*model.NetworkStatistics, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	return avgNetworkStatistics(s.dbNetStat, period)
}

// avgNetworkStatistics получить среднее значение показателей за период
func avgNetworkStatistics(s []model.NetworkStatistics, period int32) (*model.NetworkStatistics, error) {
	res := model.NetworkStatistics{}
	sumNet := map[string]*model.NetStatDetail{}

	now := time.Now().Local()
	timeStart := now.Add(-time.Second * time.Duration(period))

	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		if timeStart.Before(s[i].QueryTime) {
			for _, dv := range s[i].StatNet {
				_, ok := sumNet[dv.LocalAddress]
				if !ok {
					sumNet[dv.LocalAddress] = &model.NetStatDetail{}
				}
				sumNet[dv.LocalAddress].State = dv.State
				sumNet[dv.LocalAddress].LocalAddress = dv.LocalAddress
				sumNet[dv.LocalAddress].PeerAddress = dv.PeerAddress

				sumNet[dv.LocalAddress].Recv += dv.Recv
				sumNet[dv.LocalAddress].Send += dv.Send
			}
			count++
		} else {
			break
		}
	}
	if count > 1 {
		// расчитываем средние значения по локальному адресу
		for _, v := range sumNet {
			sumNet[v.LocalAddress].Recv = sumNet[v.LocalAddress].Recv / count
			sumNet[v.LocalAddress].Send = sumNet[v.LocalAddress].Send / count
		}
		for _, v := range sumNet {
			res.StatNet = append(res.StatNet, *v)
		}
	} else {
		// берём последнее значение
		if len(s) > 0 {
			res = s[len(s)-1]
		}
	}
	res.QueryTime = now

	return &res, nil
}
