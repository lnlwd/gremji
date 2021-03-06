package gremji

import (
    "net/url"
    "github.com/gorilla/websocket"
    "net/http"
    "errors"
    "log"
)

type Client struct {
    Remote *url.URL
    Ws *websocket.Conn
}

func NewClient(urlStr string) (*Client, error) {
    r, err := url.Parse(urlStr)
    if err != nil {
        return nil, err
    }
    dialer := websocket.Dialer{}
    ws, _, err := dialer.Dial(urlStr, http.Header{})

    if err != nil {
        log.Println(err)
        return nil, err
    }

    return &Client{r,ws}, nil
}

func (c *Client) ExecQuery(query QueryArgs) (*Response, error) {
    req := Query(query)
    return c.Exec(req)
}

func (c *Client) Exec(req *Request) (*Response, error) {
	requestMessage, err := GraphSONSerializer(req)

	if err = c.Ws.WriteMessage(websocket.BinaryMessage, requestMessage); err != nil {
		return nil, err
	}

    return c.ReadResponse()
}

func (c *Client) ReadResponse() (*Response, error) {
    res := &Response{}

    var err error

    if err = c.Ws.ReadJSON(res); err != nil {
        return nil, err
    }

    switch res.Status.Code {
    case StatusNoContent:
        return nil, nil
    case StatusSuccess:
        return res, nil
    default:
        if msg, exists := ErrorMsg[res.Status.Code]; exists {
            err := errors.New(msg)
            return nil, err
        }
        err := errors.New("An unknown error ocurred")
        return nil, err

    }
}