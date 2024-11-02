package handler

import (
	whttp "github.com/SyaibanAhmadRamadhan/http-wrapper"
	"github.com/go-chi/chi/v5"
	"order-service/generated/proto/secret_proto"
	"order-service/internal/services/order"
)

type handler struct {
	r                  *chi.Mux
	httpOtel           *whttp.Opentelemetry
	serv               serv
	jwtAccessTokenConf *secret_proto.JwtAccessToken
}

type serv struct {
	orderService order.Service
}

type Opt struct {
	JwtAccessTokenConf *secret_proto.JwtAccessToken
	OrderService       order.Service
}

func Init(r *chi.Mux, opt Opt) {
	h := &handler{
		r: r,
		httpOtel: whttp.NewOtel(
			whttp.WithRecoverMode(true),
			whttp.WithPropagator(),
			whttp.WithValidator(nil, nil),
		),
		jwtAccessTokenConf: opt.JwtAccessTokenConf,
		serv: serv{
			orderService: opt.OrderService,
		},
	}
	h.route()
}

func (h *handler) route() {
	h.r.Post("/v1/order", h.httpOtel.Trace(
		h.OrderPost,
	))
}
