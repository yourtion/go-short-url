package helper

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
)

type ContextKey int

const (
	KeyParsedJsonBody ContextKey = iota
	KeyContext
)

type H = map[string]interface{}

type Context struct {
	Res    http.ResponseWriter
	Req    *http.Request
	goNext bool
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{Res: w, Req: r}
}

func (ctx *Context) Next() {
	ctx.goNext = true
}

func (ctx *Context) GetIp() string {
	ip := ctx.Req.Header.Get("x-real-ip")
	if ip == "" {
		ip = ctx.Req.Header.Get("x-forwarded-for")
	}
	if ip == "" {
		ip = ctx.Req.RemoteAddr
	}
	ip = IpRegexp.FindString(ip)
	return ip
}

func (ctx *Context) GetUserAgent() string {
	return ctx.Req.UserAgent()
}

func (ctx *Context) GetReferer() string {
	return ctx.Req.Referer()
}

func (ctx *Context) SetData(key interface{}, value interface{}) {
	context.Set(ctx.Req, key, value)
}

func (ctx *Context) GetData(key interface{}) (interface{}, bool) {
	return context.GetOk(ctx.Req, key)
}

func (ctx *Context) GetOptionalData(key interface{}, defaultValue interface{}) interface{} {
	value, ok := context.GetOk(ctx.Req, key)
	if !ok {
		return defaultValue
	}
	return value
}

func (ctx *Context) GetCookie(name string) *http.Cookie {
	for _, c := range ctx.Req.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (ctx *Context) SetCookie(c *http.Cookie) {
	http.SetCookie(ctx.Res, c)
}

func (ctx *Context) GetUploadFile(name string) (multipart.File, *multipart.FileHeader, error) {
	if err := ctx.Req.ParseForm(); err != nil {
		return nil, nil, err
	}
	return ctx.Req.FormFile(name)
}

func (ctx *Context) GetForm(name string) (string, error) {
	if err := ctx.Req.ParseForm(); err != nil {
		return "", err
	}
	return ctx.Req.FormValue(name), nil
}

func (ctx *Context) GetParamsString(key string) string {
	vars := mux.Vars(ctx.Req)
	return vars[key]
}

func (ctx *Context) GetQuery(name string) string {
	return ctx.Req.URL.Query().Get(name)
}

func (ctx *Context) ParseJsonBody() (jsoniter.Any, error) {
	// 先检查是否已经 parse 过
	parsed, ok := ctx.GetData(KeyParsedJsonBody)
	if ok {
		body, ok := parsed.(jsoniter.Any)
		if ok {
			log.Debugf("ParseJsonBody: from cache: %+v", body)
			return body, nil
		}
	}

	contentType := ctx.Req.Header.Get("content-type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("ParseJsonBody failed: %s", contentType)
	}
	buf, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		return nil, err
	}
	data := jsoniter.Get(buf)
	log.Debugf("ParseJsonBody: %+v", data.GetInterface())
	ctx.SetData(KeyParsedJsonBody, data)
	return data, nil
}

func (ctx *Context) ResponseJson(data interface{}) {
	ctx.Res.Header().Set("content-type", "application/json")
	buf, err := jsoniter.Marshal(data)
	if err == nil {
		_, _ = ctx.Res.Write(buf)
	} else {
		log.Warnf("ctx.ResponseJson failed: %s", err)
		_, _ = ctx.Res.Write([]byte(fmt.Sprintf(`{"ok":false,"error":"%s"}`, err)))
	}
}

func (ctx *Context) ResponseError(err string) {
	ctx.ResponseJson(map[string]interface{}{
		"ok":    false,
		"error": err,
	})
}

func (ctx *Context) ResponseOk(data interface{}) {
	ctx.ResponseJson(map[string]interface{}{
		"ok":   true,
		"data": data,
	})
}

func (ctx *Context) ResponseText(text string) {
	_, _ = ctx.Res.Write([]byte(text))
}

func (ctx *Context) ResponseHTML(html string) {
	ctx.Res.Header().Set("content-type", "text/html")
	_, _ = ctx.Res.Write([]byte(html))
}

func (ctx *Context) Redirect(url string) {
	ctx.Res.Header().Set("location", url)
	ctx.Res.WriteHeader(302)
	_, _ = ctx.Res.Write([]byte("Redirect to " + url + "..."))
}
