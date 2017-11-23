package utils

import (
	"fmt"
	log4 "github.com/alecthomas/log4go"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type WebSocketOptions struct {
	HeartbeatInterval time.Duration
	HeartbeatPkg      []byte
}

type WebSocketClient struct {
	host              string
	opts              *WebSocketOptions
	conn              *websocket.Conn
	recvCh            chan []byte
	existCh           chan interface{}
	lastHeartbeatTime time.Time
	lock              sync.RWMutex
	status            bool
}

func NewWebSocketClient(host string, opts ...*WebSocketOptions) *WebSocketClient {
	var options *WebSocketOptions
	if len(opts) == 0 {
		options = &WebSocketOptions{
			HeartbeatInterval: 30 * time.Second,
			HeartbeatPkg:      []byte(`{"Action":"heartbeat"}`),
		}
	} else {
		options = opts[0]
	}
	return &WebSocketClient{
		host:              host,
		opts:              options,
		recvCh:            make(chan []byte, 1),
		existCh:           make(chan interface{}, 0),
		lastHeartbeatTime: time.Now(),
	}
}

func (p *WebSocketClient) Connet() (recvCh chan []byte, err error) {
	p.conn, _, err = websocket.DefaultDialer.Dial(p.host, nil)
	if err != nil {
		return nil, err
	}
	p.status = true
	go p.doRecv()
	go p.heartbeat()
	return p.recvCh, nil
}

func (p *WebSocketClient) updateHearbeatTime() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.lastHeartbeatTime = time.Now()
}

func (p *WebSocketClient) getHeartbeatTime() time.Time {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.lastHeartbeatTime
}

func (p *WebSocketClient) Send(data []byte) error {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if !p.status {
		return fmt.Errorf("WebSocket connect has already closed.")
	}
	return p.conn.WriteMessage(websocket.TextMessage, data)
}

func (p *WebSocketClient) doRecv() {
	defer close(p.recvCh)
	for {
		_, data, err := p.conn.ReadMessage()
		if err != nil {
			if p.Status() {
				log4.Error("WebSocketClient host:%s ReadMessage error:%s", p.host, err)
				p.Close()
			}
			return
		}
		p.updateHearbeatTime()
		p.recvCh <- data
	}
}

func (p *WebSocketClient) heartbeat() {
	timer := time.NewTimer(p.opts.HeartbeatInterval)
	defer timer.Stop()
	for {
		select {
		case <-p.existCh:
			return
		case <-timer.C:
			err := p.Send(p.opts.HeartbeatPkg)
			if err != nil {
				log4.Error("WebSocketClient send heartbeat error:%s", err)
			}
			timer.Reset(p.opts.HeartbeatInterval)
		}
	}
}

func (p *WebSocketClient) Status() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.status
}

func (p *WebSocketClient) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()

	if !p.status {
		return
	}
	p.status = false
	close(p.existCh)
	err := p.conn.Close()
	if err != nil {
		log4.Error("WebSocketClient close error:%s", err)
	}
}
