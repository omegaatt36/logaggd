package appcontext

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPError defines error under HTTP handler.
type HTTPError interface {
	error
	Code() int
}

// AppContext is a gin context wrapper that includes session information.
type AppContext struct {
	GinCtx *gin.Context
	ReqID  string

	Resp     Response
	HTTPCode int
	Done     bool
}

// fromGinContext generates AppContext from gin.Context.
func fromGinContext(c *gin.Context) *AppContext {
	ctx := &AppContext{
		GinCtx: c,
		ReqID:  randomString(20),
	}
	return ctx
}

// Apply returns gin HTTP handler that executes handler accepts AppContext.
func Apply(handler func(c *AppContext)) func(*gin.Context) {

	return func(x *gin.Context) {
		// install ctx.
		ctx := fromGinContext(x)
		x.Set("AppContext", ctx)

		raw, ok := x.Get("CachedResult")
		if ok {
			value, ok := raw.([]byte)
			if ok {
				x.Writer.Header().Set("Backend-Source", "redis")
				x.Data(http.StatusOK, gin.MIMEJSON, value)
				x.Abort()
				ctx.Done = true
				return
			}
		}

		handler(ctx)
	}
}

// AbortFromError returns error message with 500.
func (c *AppContext) AbortFromError(err HTTPError) {
	type SetBaseInfoer interface {
		SetBaseInfo(email, method, URI string)
	}
	y, ok := err.(SetBaseInfoer)
	if ok {
		var (
			email  = "invalid"
			method = "invalid"
			URI    = "invalid"
		)
		if c.GinCtx != nil {
			if c.GinCtx.Request != nil {
				method = c.GinCtx.Request.Method

				if c.GinCtx.Request.URL != nil {
					URI = c.GinCtx.Request.URL.Path
				}
			}
		}

		y.SetBaseInfo(email, method, URI)
	}
	c.HTTPCode = err.Code()
	c.Resp = Response{
		Err: Error{Success: false, Msg: err.Error(), ReqID: c.ReqID},
	}
}

// Error defines HTTP's response when any errors happened.
type Error struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	ReqID   string `json:"reqid"`
}

// Response is universal response in HTTP.
type Response struct {
	Err    Error       `json:"error"`
	Result interface{} `json:"result"`
}

// OK returns json response with 200.
func (c *AppContext) OK(x interface{}) {
	c.HTTPCode = http.StatusOK
	c.Resp = Response{Err: Error{Success: true}, Result: x}
}

const charset = "abcdef0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return stringWithCharset(length, charset)
}
