package services

import (
	"github.com/gin-gonic/gin"
	"github.com/saifwork/price-tracker-bot.git/app/configs"
	"github.com/saifwork/price-tracker-bot.git/app/services/domains"
)

type Initializer struct {
	gin  *gin.Engine
	conf *configs.Config
}

func NewInitializer(gin *gin.Engine, conf *configs.Config) *Initializer {
	s := &Initializer{
		gin:  gin,
		conf: conf,
	}
	return s
}

func (s *Initializer) RegisterDomains(domains []domains.IDomain) {
	for _, domain := range domains {
		domain.SetupRoutes()
	}
}
