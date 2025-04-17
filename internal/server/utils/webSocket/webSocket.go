package websocket

import (
	"ai-calls/internal/config"
	"context"
	"net/http"

	"github.com/coder/websocket"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type SafeConn struct {
	Conn        *websocket.Conn
	SessionName string
	IsClosed    bool
	Config      *config.Config
}

func newSocket(conn *websocket.Conn, config *config.Config) *SafeConn {
	conn.SetReadLimit(config.WebSocket.ReadLimit)
	return &SafeConn{Conn: conn, Config: config}
}

func NewServer(config *config.Config, w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*SafeConn, error) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return nil, err
	}
	return newSocket(conn, config), nil
}

func NewClient(config *config.Config, ctx context.Context, urlStr string, requestHeader http.Header) (*SafeConn, *http.Response, error) {
	conn, response, err := websocket.Dial(ctx, urlStr, nil)
	if err != nil {
		return nil, nil, err
	}
	socketConn := newSocket(conn, config)
	return socketConn, response, err
}

func (s *SafeConn) WriteMessage(ctx context.Context, data []byte) error {
	return s.Conn.Write(ctx, websocket.MessageText, data)
}

func (s *SafeConn) WriteJsonMessage(ctx context.Context, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.Conn.Write(ctx, websocket.MessageText, jsonData)
}

func (s *SafeConn) ReadMessage(ctx context.Context) (websocket.MessageType, []byte, error) {
	return s.Conn.Read(ctx)
}

func (s *SafeConn) ReadJsonMessage(ctx context.Context, data any) error {
	_, msg, err := s.Conn.Read(ctx)
	if err != nil {
		return err
	}
	err = json.Unmarshal(msg, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *SafeConn) SetConnectName(name string) *SafeConn {
	s.SessionName = name
	return s
}

func (s *SafeConn) Close(ctx context.Context) error {
	s.IsClosed = true
	return s.Conn.Close(websocket.StatusNormalClosure, "")
}
