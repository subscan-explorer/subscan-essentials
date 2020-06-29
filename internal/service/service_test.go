package service

var testSrv Service

func init() {
	testSrv = Service{
		dao: nil,
	}
}
