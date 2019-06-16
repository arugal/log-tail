package control

import (
	"encoding/json"
	"github.com/Arugal/log-tail/g"
	"github.com/Arugal/log-tail/models/config"
	"github.com/Arugal/log-tail/util/log"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"os"
	"strings"
	"sync/atomic"
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
	Type    Type `json:"type"`
	UiWidth int  `json:"ui_width"`
}

func (t *TailReqProtocol) LineNum() int {
	if t.UiWidth == 0 {
		return 200
	} else {
		return t.UiWidth/9 + 10
	}
}

type ConnManager struct {
	cChan        chan *ConnCarrier
	carrieMap    map[uint64]*ConnCarrier
	idCarrierSeq uint64
	log          log.Logger
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		cChan:     make(chan *ConnCarrier),
		carrieMap: make(map[uint64]*ConnCarrier),
		log:       log.NewPrefixLogger("conn-manager"),
	}
}

func (cc *ConnManager) Run() {
	go func() {
		cc.log.Info("Start ConnManager heartInterval %s maxInterval:%s", time.Second.String(), time.Minute.String())
		heartTimer := time.NewTicker(time.Second)
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

func (cc *ConnManager) GenerateCarrierId() uint64 {
	return atomic.AddUint64(&cc.idCarrierSeq, 1)
}

func (cc *ConnManager) AddConnCarrier(carrier *ConnCarrier) {
	carrier.Add = true
	cc.cChan <- carrier
}

func (cc *ConnManager) delConnCarrier(carrier *ConnCarrier) {
	carrier.Add = false
	cc.cChan <- carrier
}

func (cc *ConnManager) ProcessNewConn() {
	for {
		carrie := <-cc.cChan
		if carrie.Add {
			if _, ok := cc.carrieMap[carrie.Id()]; ok {
				cc.log.Warn("carrie existing %s", carrie.String())
			} else {
				cc.log.Trace("receive new carrie %s", carrie.String())
				go carrie.Handler()
				cc.carrieMap[carrie.Id()] = carrie
			}
		} else {
			delete(cc.carrieMap, carrie.Id())
		}
	}
}

func (cc *ConnManager) CheckHeartTimeout() {
	if len(cc.carrieMap) > 0 {
		currentTime := time.Now().Unix()
		heartInterval := int64(g.GlbServerCfg.HeartInterval/time.Second) * 2 // s
		cc.log.Trace("go check heart time out %d - %d carries:%d", currentTime, heartInterval, len(cc.carrieMap))

		for _, carrier := range cc.carrieMap {
			cc.log.Trace("check %v heart time out", carrier.String())
			if currentTime-heartInterval > carrier.LastHeartTime {
				if carrier.Running {
					carrier.Close()
					cc.log.Info("heart time out auto close %s", carrier.String())
				}
				cc.delConnCarrier(carrier)
			}
		}
	}
}

func (cc *ConnManager) CheckConnMaxTime() {
	if len(cc.carrieMap) > 0 {
		deadline := time.Now().Unix() - int64(g.GlbServerCfg.ConnMaxTime/time.Second) // s
		cc.log.Trace("go check conn max time %d carries:%d", deadline, len(cc.carrieMap))

		for _, carrier := range cc.carrieMap {
			cc.log.Trace("check %v max time %v", carrier.String())
			if deadline > carrier.StartTime {
				if carrier.Running {
					carrier.Close()
					cc.log.Info("to achieve max time auto close %s", carrier.String())
				}
				cc.delConnCarrier(carrier)
			}
		}
	}
}

type ConnCarrier struct {
	id            uint64
	Conn          *websocket.Conn
	Cf            config.CatalogConf
	File          string
	Tail          *tail.Tail
	StartTime     int64
	LastHeartTime int64
	Peer          string
	Add           bool
	handerDone    chan bool
	Running       bool
	log           log.Logger
}

func NewConnCarrier(cm *ConnManager, conn *websocket.Conn, cf config.CatalogConf, file string) ConnCarrier {
	return ConnCarrier{
		id:            cm.GenerateCarrierId(),
		Conn:          conn,
		Cf:            cf,
		File:          file,
		StartTime:     time.Now().Unix(),
		LastHeartTime: time.Now().Unix(),
		handerDone:    make(chan bool, 1),
		Peer:          conn.RemoteAddr().String(),
		Add:           true,
		Running:       true,
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

				go func(cc *ConnCarrier, lineNum, msgType int) {
					for line := range cc.Tail.Lines {
						lenL := strings.Count(line.Text, "") - 1
						if lenL > lineNum {
							for i := 0; i < lenL; i += lineNum {
								max := i + lineNum
								if max > lenL {
									max = lenL
								}
								WriteLine(cc, line.Text[i:max], msgType)
							}
						} else {
							WriteLine(cc, line.Text, msgType)
						}

					}
					cc.log.Debug("Tail Done %s", cc.String())
				}(cc, req.LineNum(), msgType)
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

func WriteLine(cc *ConnCarrier, line string, msgType int) {
	_ = cc.Conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	resp := TailRespProtocol{
		Type: Write,
		Msg:  line,
	}
	buf, _ := json.Marshal(resp)
	err := cc.Conn.WriteMessage(msgType, buf)
	if err != nil {
		cc.log.Error("Tail write message err %s case:%v", cc.String(), err)
	}
	cc.log.Trace("send log line %s", string(buf))
	time.Sleep(time.Microsecond * 200)
}

func (cc *ConnCarrier) Id() uint64 {
	return cc.id
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
	cc.Running = false
}
