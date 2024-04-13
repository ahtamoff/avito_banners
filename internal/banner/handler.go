package banner

import (
	"avito_banners/internal/handlers"
	"avito_banners/pkg/logging"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//чтобы посмотреть реализует ли handler сущности основной интерфейс Handler
var _ handlers.Handler = &handler{}

const (
	userBannerUrl = "/user_banner"
	bannerUrl = "/banner"
	bannerIdUrl = "/banner/:uuid"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger)handlers.Handler{
	return &handler{
		logger: logger,
	}
}

func(h *handler) Register(router *httprouter.Router){
	router.GET(userBannerUrl, h.GetUserBanner)
	router.GET(bannerUrl, h.GetBanner)
	router.POST(bannerUrl, h.CreateBanner)
	router.PATCH(bannerIdUrl, h.UpdateBanner)
	router.DELETE(bannerIdUrl, h.DeleteBanner)
}

func(h *handler) GetUserBanner(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	w.WriteHeader(200)
	w.Write([]byte("GetUserBanner"))
}

func(h *handler) GetBanner(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	w.WriteHeader(200)
	w.Write([]byte("GetBanner"))
}


func(h *handler) CreateBanner(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	w.WriteHeader(201)
	w.Write([]byte("CreateBanner"))
}

func(h *handler) UpdateBanner(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("create banner")))
}

func(h *handler) DeleteBanner(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	w.WriteHeader(204)
	w.Write([]byte(fmt.Sprintf("delete banner")))
}