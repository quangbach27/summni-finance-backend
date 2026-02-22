package ports

import (
	"sumni-finance-backend/internal/finance/app"
)

type HttpServer struct {
	application app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		application: application,
	}
}
