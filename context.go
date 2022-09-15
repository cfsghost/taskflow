package taskflow

import "sync"

var ctxPool = sync.Pool{
	New: func() interface{} {
		return &Context{}
	},
}

type Context struct {
	meta     sync.Map
	privData interface{}
}

func NewContext() *Context {
	return ctxPool.Get().(*Context)
}

func (ctx *Context) SetPrivData(privData interface{}) {
	ctx.privData = privData
}

func (ctx *Context) GetPrivData() interface{} {
	return ctx.privData
}

func (ctx *Context) SetMeta(key interface{}, value interface{}) {
	ctx.meta.Store(key, value)
}

func (ctx *Context) RemoveMeta(key interface{}) {
	ctx.meta.Delete(key)
}

func (ctx *Context) GetMeta(key interface{}) (interface{}, bool) {
	return ctx.meta.Load(key)
}

func (ctx *Context) Reset() {
	ctx.privData = nil
	ctx.meta.Range(func(key interface{}, value interface{}) bool {
		ctx.meta.Delete(key)
		return true
	})

}
