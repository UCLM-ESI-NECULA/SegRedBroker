package service

type BrokerServiceImpl struct {
}

func NewBrokerService() *BrokerServiceImpl {
	return &BrokerServiceImpl{}
}

type BrokerService interface {
	GetVersion() map[string]string
}

func (svc *BrokerServiceImpl) GetVersion() map[string]string {
	return map[string]string{
		"version": "1.0",
	}
}
