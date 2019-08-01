package server

import (
	"encoding/json"
	"github.com/arugal/log-tail/server/control"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type GeneralResponse struct {
	Code int
	Msg  string
}

type ErrorInfo struct {
	Err error  `json:"err"`
	Msg string `json:"msg"`
}

type GetCataLogInfo struct {
	Catalog   string   `json:"catalog"`
	ChildFile []string `json:"child_file"`
}

func Write(svr *Service, res *GeneralResponse, w http.ResponseWriter, r *http.Request) {
	svr.log.Info("Http response [%s]: code [%d]", r.URL.Path, res.Code)
	w.WriteHeader(res.Code)
	if len(res.Msg) > 0 {
		_, _ = w.Write([]byte(res.Msg))
	}
}

func (srv *Service) GetCataLogInfo(w http.ResponseWriter, r *http.Request) {
	res := GeneralResponse{Code: 200}
	defer Write(srv, &res, w, r)
	srv.log.Info("Http requst: [%s]", r.URL.Path)
	reps := srv.GetCataLog()
	buf, _ := json.Marshal(reps)
	res.Msg = string(buf)
}

func (srv *Service) GetCataLog() (reps []GetCataLogInfo) {
	catalogs := srv.cm.GetAllCatalogInfo()
	for _, conf := range catalogs {
		info := GetCataLogInfo{
			Catalog:   conf.Name,
			ChildFile: conf.ChildFile,
		}
		reps = append(reps, info)
	}
	return
}

func (srv *Service) GetLogTail(w http.ResponseWriter, r *http.Request) {
	res := GeneralResponse{Code: 200}
	params := mux.Vars(r)
	catalog := params["catalog"]
	file := params["file"]

	srv.log.Info("Http request: [%s]", r.URL.Path)

	cf, ok := srv.cm.GetCatalogInfo(catalog)
	if ok {
		if cf.HasChildFile(file) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				srv.log.Error("Upgrader websocket failing [%s] case:%v", r.URL.Path, err)
				res.Code = 500
				info := ErrorInfo{Err: err, Msg: "Upgrader websocket failing"}
				buf, _ := json.Marshal(info)
				res.Msg = string(buf)
				Write(srv, &res, w, r)
				return
			}
			carrier := control.NewConnCarrier(srv.cm2, conn, cf, file)
			srv.cm2.AddConnCarrier(&carrier)
			return
		} else {
			srv.log.Error("Http request: [%s] child file miss [%s]", r.URL.Path, file)
		}
	} else {
		srv.log.Error("Http request: [%s] catalog miss [%s]", r.URL.Path, catalog)
	}
	res.Code = 404
	Write(srv, &res, w, r)
}
