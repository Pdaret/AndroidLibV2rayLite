package keys

const APP_KEY = "303a223afec7112541253a81c1043b6107005b3a2856921ca503bd5c7d6c5163"
const SERVER_KEY = "36fb63a6dbb3bef7b6e4470a0bd06e97a187f579d2d6448d45a7bde38d135f8a"

type KeyService interface {
	GetAppKey() string
	GetServerKey() string
}

type keyService struct{}

func Initialize() KeyService {
	return &keyService{}
}

func (k *keyService) GetAppKey() string {
	return APP_KEY
}

func (k *keyService) GetServerKey() string {
	return SERVER_KEY
}
