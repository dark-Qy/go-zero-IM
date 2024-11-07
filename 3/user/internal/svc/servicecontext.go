package svc

import (
	"3v3/user/internal/config"
	"3v3/user/models"
)

type ServiceContext struct {
	Config config.Config

	UserModel models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	//sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,

		//UserModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}
