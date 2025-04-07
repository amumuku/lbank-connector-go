package sve

import "fmt"

type SpotService struct {
	c  *Client
	hs *HttpService
}

// GetHttpService 返回 HttpService 实例
// func (s *SpotService) GetHttpService() *HttpService {
// 	return s.hs
// }

func (s *SpotService) CreateOrder(data map[string]string) {
	url := s.c.Host + PathSupplementCreatOrder
	params := s.hs.BuildSignBody(data)
	s.hs.Post(url, params)
}

// CancelClientOrders 批量取消订单
func (s *SpotService) CancelClientOrders(orders []map[string]string) ([]string, error) {
	var responses []string
	for _, order := range orders {
		url := s.c.Host + PathSupplementCancelOrder
		params := s.hs.BuildSignBody(order)
		s.hs.Post(url, params)
		if s.hs.Error != nil {
			return responses, fmt.Errorf("取消订单 %s 失败: %v", order["origClientOrderId"], s.hs.Error)
		}
		responses = append(responses, s.hs.Text)
	}
	return responses, nil
}

func (s *SpotService) CancelOrder(data map[string]string) {
	url := s.c.Host + PathSupplementCancelOrder
	params := s.hs.BuildSignBody(data)
	s.hs.Post(url, params)
}

func (s *SpotService) CancelOrderBySymbol(data map[string]string) {
	url := s.c.Host + PathSupplementCancelOrderBySymbol
	params := s.hs.BuildSignBody(data)
	s.hs.Post(url, params)
}

func (s *SpotService) OrdersInfo(data map[string]string) {
	url := s.c.Host + PathSupplementOrdersInfo
	params := s.hs.BuildSignBody(data)
	s.hs.Post(url, params)
}

// func (s *SpotService) OrdersInfoNoDeal(data map[string]string) {
// 	url := s.c.Host + PathSupplementOrdersInfoNoDeal
// 	params := s.hs.BuildSignBody(data)
// 	s.hs.Post(url, params)
// }

func (s *SpotService) OrdersInfoNoDeal(data map[string]string) ([]byte, error) {
	url := s.c.Host + PathSupplementOrdersInfoNoDeal
	params := s.hs.BuildSignBody(data)
	s.hs.Post(url, params)
	if s.hs.Error != nil {
		return nil, s.hs.Error
	}
	return []byte(s.hs.Text), nil
}

func (s *SpotService) OrdersInfoHistory(data map[string]string) {
	url := s.c.Host + PathSupplementOrdersInfoHistory
	params := s.hs.BuildSignBody(data)
	s.hs.Post(url, params)
}

func (s *SpotService) UserInfoAccount(data map[string]string) {
	url := s.c.Host + PathSupplementUserInfoAccount
	params := s.hs.BuildSignBody(data)
	s.hs.Post(url, params)
}
func (s *SpotService) TransactionHistory(data map[string]string) {
	url := s.c.Host + PathSupplementTransactionHistory
	params := s.hs.BuildSignBody(data)
	s.hs.Post(url, params)
}
