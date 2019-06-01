package control

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"log-tail/g"
	"log-tail/models/config"
	"log-tail/util/log"
	"os"
	"sync"
	"time"
)

type Type int8

const (
	Read Type = iota
	Write
	Heart
	Close
)

type TailRespProtocol struct {
	Type Type   `json:"type"`
	Line string `json:"line"`
}

type TailReqProtocol struct {
	Type Type `json:"type"`
}

type ConnManager struct {
	Cc       chan ConnCarrier
	carries  []ConnCarrier
	carrieMu sync.Mutex
	log      log.Logger
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Cc:       make(chan ConnCarrier, 10),
		carries:  []ConnCarrier{},
		carrieMu: sync.Mutex{},
		log:      log.NewPrefixLogger("ConnManager"),
	}
}

func (cc *ConnManager) Run() {

	go func() {
		heartTimer := time.NewTimer(time.Second)
		maxTimer := time.NewTimer(time.Minute)
		for {
			select {
			case <-heartTimer.C:
				cc.CheckHeartTimeout()
			case <-maxTimer.C:
				cc.CheckConnMaxTime()
			}
		}
	}()

	go cc.ProcessNewConn()
}

func (cc *ConnManager) refreshCarrier(newCarries []ConnCarrier) {
	cc.carrieMu.Lock()
	defer cc.carrieMu.Unlock()
	cc.carries = newCarries
}

func (cc *ConnManager) ProcessNewConn() {
	for {
		carrie := <-cc.Cc

		go carrie.Handler()

		cc.carrieMu.Lock()
		defer cc.carrieMu.Unlock()
		cc.carries = append(cc.carries, carrie)
	}
}

func (cc *ConnManager) CheckHeartTimeout() {
	currentTime := time.Now().Unix()
	newCarries := []ConnCarrier{}
	heartInterval := int64(g.GlbServerCfg.HeartInterval)

	for _, carrier := range cc.carries {
		if currentTime-heartInterval <= carrier.LastHeartTime {
			newCarries = append(newCarries, carrier)
		} else {
			carrier.Close()
			cc.log.Info("heart time out auto close %v", carrier)
		}
	}
	cc.refreshCarrier(newCarries)
}

func (cc *ConnManager) CheckConnMaxTime() {
	deadline := time.Now().Unix() - int64(g.GlbServerCfg.ConnMaxTime)
	newCarries := []ConnCarrier{}
	for _, carrier := range cc.carries {
		if deadline <= carrier.StartTime {
			newCarries = append(newCarries, carrier)
		} else {
			carrier.Close()
			cc.log.Info("to achieve max time auto close %v", carrier)
		}
	}
	cc.refreshCarrier(newCarries)
}

type ConnCarrier struct {
	Conn          *websocket.Conn
	Cf            config.CatalogConf
	File          string
	Tail          *tail.Tail
	StartTime     int64
	LastHeartTime int64
	tailDone      chan bool
	handerDone    chan bool
	log           log.Logger
}

func NewConnCarrier(conn *websocket.Conn, cf config.CatalogConf, file string) ConnCarrier {
	return ConnCarrier{
		Conn:          conn,
		Cf:            cf,
		File:          file,
		StartTime:     time.Now().Unix(),
		LastHeartTime: time.Now().Unix(),
		tailDone:      make(chan bool),
		handerDone:    make(chan bool),
		log:           log.NewPrefixLogger(cf.Name + ":" + file),
	}
}

func (cc *ConnCarrier) Handler() {
	for {
		select {
		case <-cc.handerDone:
			return
		default:
			_ = cc.Conn.SetReadDeadline(time.Now().Add(time.Duration(g.GlbServerCfg.HeartInterval) * 2))
			msgType, msg, err := cc.Conn.ReadMessage()
			if err != nil {
				cc.log.Error("Read message err case:%v", err)
				continue
			}

			var req TailReqProtocol
			err = json.Unmarshal([]byte(msg), &req)
			if err != nil {
				cc.log.Error("Unmarshal %s err case:%v", msg, err)
				continue
			}

			switch req.Type {
			case Read:
				fileInfo, err := os.Stat(cc.Cf.FullFilePath(cc.File))
				if err != nil {
					cc.log.Error("os stat err case:%v", err)
					cc.Close()
					return
				}
				var offset int64
				if fileInfo.Size() > g.GlbServerCfg.LastReadOffset {
					offset = fileInfo.Size() - g.GlbServerCfg.LastReadOffset
				} else {
					offset = 0
				}

				t, err := tail.TailFile(cc.Cf.FullFilePath(cc.File), tail.Config{Location: &tail.SeekInfo{Offset: int64(offset), Whence: 10}, Follow: true})
				if err != nil {
					cc.log.Error("tail file err case:%v", err)
					cc.Close()
					return
				}
				cc.Tail = t

				go func(cc *ConnCarrier, msgType int) {
					for {
						select {
						case line := <-cc.Tail.Lines:
							resp := TailRespProtocol{
								Type: Write,
								Line: line.Text,
							}
							buf, _ := json.Marshal(resp)
							err = cc.Conn.WriteMessage(msgType, buf)
						case <-cc.tailDone:
							return
						}
					}
				}(cc, msgType)
			case Heart:
				cc.LastHeartTime = time.Now().Unix()
			case Close:
				cc.Close()
			}
		}
	}
}

func (cc *ConnCarrier) Close() {
	cc.tailDone <- true
	cc.handerDone <- true
	_ = cc.Conn.Close()
	cc.Tail.Cleanup()
	_ = cc.Tail.Stop()
}
