package utils

import "time"

var (
	// SessMap 存sessionid
	SessMap = make(map[string]*Session, 64)
)

func init() {
	go checkExpire()
}

// Session session
type Session struct {
	// ID  sessionId
	ID string
	// 最长保存时间
	maxAge time.Duration
	latest int64
	data   interface{}
}

// NewSession 构造session
func NewSession(id string, maxAge time.Duration, latest int64, data interface{}) *Session {
	return &Session{id, maxAge, latest, data}
}

// Update  更新最近访问时间
func (s *Session) Update() {
	s.latest = time.Now().UnixNano()
}

// Delete  删除Session
func (s *Session) Delete() {
	delete(SessMap, s.ID)
}

func gc() {
	nowTime := time.Now().UnixNano()
	for k, v := range SessMap {
		if (v.latest + int64(v.maxAge)) <= nowTime {
			delete(SessMap, k)
		}
	}
}

// CheckSess 检查session是否存在
func CheckSess(sessionID string) bool {
	if _, ok := SessMap[sessionID]; ok {
		return true
	}
	return false
}

func checkExpire() {
	for {
		select {
		case <-time.Tick(time.Second * 10):
			gc()
		}
	}
}
