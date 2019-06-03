package control

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"log-tail/g"
	"log-tail/models/config"
	"log-tail/util/log"
	"os"
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
	Msg  string `json:"msg"`
}

type TailReqProtocol struct {
	Type Type `json:"type"`
}

type ConnManager struct {
	CChan     chan ConnCarrier
	carrieMap map[string]ConnCarrier
	log       log.Logger
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		CChan:     make(chan ConnCarrier),
		carrieMap: make(map[string]ConnCarrier),
		log:       log.NewPrefixLogger("ConnManager"),
	}
}

func (cc *ConnManager) Run() {
	go func() {
		heartTimer := time.NewTicker(g.GlbServerCfg.HeartInterval)
		maxTimer := time.NewTicker(time.Minute)
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

func (cc *ConnManager) AddConnCarrier(carrier ConnCarrier) {
	carrier.Add = true
	cc.CChan <- carrier
}

func (cc *ConnManager) delConnCarrier(carrier ConnCarrier) {
	carrier.Add = false
	cc.CChan <- carrier
}

func (cc *ConnManager) ProcessNewConn() {
	for {
		carrie := <-cc.CChan
		if carrie.Add {
			if _, ok := cc.carrieMap[carrie.String()]; ok {
				cc.log.Warn("carrie existing %s", carrie.String())
			} else {
				cc.log.Trace("receive new carrie %s", carrie.String())
				go carrie.Handler()
				cc.carrieMap[carrie.String()] = carrie
			}
		} else {
			delete(cc.carrieMap, carrie.String())
		}
	}
}

func (cc *ConnManager) CheckHeartTimeout() {
	if len(cc.carrieMap) > 0 {
		currentTime := time.Now().Unix()
		heartInterval := int64(g.GlbServerCfg.HeartInterval) * 2
		cc.log.Trace("go check heart time out %d - %d carries:%d", currentTime, heartInterval, len(cc.carrieMap))

		for _, carrier := range cc.carrieMap {
			cc.log.Trace("check %v heart time out", carrier.String())
			if currentTime-heartInterval > carrier.LastHeartTime {
				carrier.Close()
				cc.log.Info("heart time out auto close %s", carrier.String())
				cc.delConnCarrier(carrier)
			}
		}
	}
}

func (cc *ConnManager) CheckConnMaxTime() {
	if len(cc.carrieMap) > 0 {
		deadline := time.Now().Unix() - int64(g.GlbServerCfg.ConnMaxTime)
		cc.log.Trace("go check conn max time %d carries:%d", deadline, len(cc.carrieMap))

		for _, carrier := range cc.carrieMap {
			cc.log.Trace("check %v max time %v", carrier.String())
			if deadline > carrier.StartTime {
				carrier.Close()
				cc.log.Info("to achieve max time auto close %s", carrier.String())
				cc.delConnCarrier(carrier)
			}
		}
	}
}

type ConnCarrier struct {
	Conn          *websocket.Conn
	Cf            config.CatalogConf
	File          string
	Tail          *tail.Tail
	StartTime     int64
	LastHeartTime int64
	Peer          string
	Add           bool
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
		handerDone:    make(chan bool, 1),
		Peer:          conn.RemoteAddr().String(),
		Add:           true,
		log:           log.NewPrefixLogger(cf.Name + ":" + file),
	}
}

func (cc *ConnCarrier) Handler() {
	for {
		select {
		case <-cc.handerDone:
			cc.log.Debug("Hander done %s", cc.String())
			return
		default:
			_ = cc.Conn.SetReadDeadline(time.Now().Add(time.Duration(g.GlbServerCfg.HeartInterval) * 2))
			msgType, msg, err := cc.Conn.ReadMessage()
			if err != nil {
				cc.log.Error("Read message err case:%v", err)
				cc.Close()
				continue
			}

			var req TailReqProtocol
			err = json.Unmarshal(msg, &req)
			if err != nil {
				cc.log.Error("Unmarshal %s err case:%v", msg, err)
				continue
			}

			cc.log.Trace("received msg %v", req)

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

				t, err := tail.TailFile(cc.Cf.FullFilePath(cc.File), tail.Config{Location: &tail.SeekInfo{Offset: int64(offset), Whence: 0}, Follow: true})
				if err != nil {
					cc.log.Error("tail file err case:%v", err)
					cc.Close()
					return
				}
				cc.Tail = t

				go func(cc *ConnCarrier, msgType int) {
					for line := range cc.Tail.Lines {
						_ = cc.Conn.SetWriteDeadline(time.Now().Add(time.Second))
						resp := TailRespProtocol{
							Type: Write,
							Msg:  line.Text,
						}
						buf, _ := json.Marshal(resp)
						err := cc.Conn.WriteMessage(msgType, buf)
						if err != nil {
							cc.log.Error("Tail write message err %s case:%v", cc.String(), err)
							return
						}
						cc.log.Trace("send log line %s", string(buf))
					}
					cc.log.Debug("Tail Done %s", cc.String())
				}(cc, msgType)
				resp := TailRespProtocol{
					Type: Heart,
					Msg:  g.GlbServerCfg.HeartInterval.String(),
				}
				buf, _ := json.Marshal(resp)
				_ = cc.Conn.SetWriteDeadline(time.Now().Add(time.Second))
				err = cc.Conn.WriteMessage(msgType, buf)
				if err != nil {
					cc.log.Error("Tail write message err case:%v", err)
					cc.Close()
					return
				}
				cc.log.Trace("send heart interval %s", string(buf))
			case Heart:
				cc.LastHeartTime = time.Now().Unix()
			case Close:
				cc.Close()
			}
		}
	}
}

func (cc *ConnCarrier) String() string {
	return cc.Peer + "-" + cc.Cf.Name + "-" + cc.File
}

func (cc *ConnCarrier) Close() {
	cc.handerDone <- true
	_ = cc.Conn.Close()
	if cc.Tail != nil {
		cc.Tail.Cleanup()
		_ = cc.Tail.Stop()
	}
}
