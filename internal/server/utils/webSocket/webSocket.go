package websocket

import (
	"context"
	"net/http"

	"github.com/coder/websocket"
)

type SafeConn struct {
	Conn        *websocket.Conn
	SessionName string
	IsClosed    bool
}

func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*SafeConn, error) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return nil, err
	}
	conn.SetReadLimit(1 << 30)
	return &SafeConn{Conn: conn}, nil
}

func NewWebSocketConnect(ctx context.Context, urlStr string, requestHeader http.Header) (*SafeConn, *http.Response, error) {
	conn, response, err := websocket.Dial(ctx, urlStr, nil)
	if err != nil {
		return nil, nil, err
	}
	conn.SetReadLimit(1 << 30)
	socketConn := &SafeConn{Conn: conn}
	return socketConn, response, err
}

func (s *SafeConn) WriteMessage(ctx context.Context, data []byte) error {
	return s.Conn.Write(ctx, websocket.MessageText, data)
}

func (s *SafeConn) ReadMessage(ctx context.Context) (websocket.MessageType, []byte, error) {
	return s.Conn.Read(ctx)
}

func (s *SafeConn) SetConnectName(name string) *SafeConn {
	s.SessionName = name
	return s
}

func (s *SafeConn) Close(ctx context.Context) error {
	s.IsClosed = true
	return s.Conn.Close(websocket.StatusNormalClosure, "")
}
