package acExt

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"math"
	"strings"
	"time"
)
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

func NodesSelector[T cdp.NodeID | cdp.Node | *cdp.Node | int64](ns ...T) []cdp.NodeID {
	ids := make([]cdp.NodeID, 0)
	for _, n := range ns {
		switch v := (any)(n).(type) {
		case cdp.NodeID:
			ids = append(ids, v)
		case *cdp.NodeID:
			ids = append(ids, *v)
		case int:
			ids = append(ids, cdp.NodeID(v))
		case int64:
			ids = append(ids, cdp.NodeID(v))
		case cdp.Node:
			ids = append(ids, v.NodeID)
		case *cdp.Node:
			ids = append(ids, v.NodeID)
		}
	}
	return ids
}
func WithLog(ctx context.Context,l StdLogger)context.Context{
	return context.WithValue(ctx,"~log~",l)
}
func Panicf(format string, args ...interface{}) chromedp.QueryAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		l,ok:=ctx.Value("~log~").(StdLogger)
		if ok{
			l.Panicf(format,args...)
		}
		return nil
	})
}
func Panic(args ...interface{}) chromedp.QueryAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		l,ok:=ctx.Value("~log~").(StdLogger)
		if ok{
			l.Panic(args...)
		}
		return nil
	})
}
func Panicln(args ...interface{}) chromedp.QueryAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		l,ok:=ctx.Value("~log~").(StdLogger)
		if ok{
			l.Panicln(args...)
		}
		return nil
	})
}
func Fatalf(format string, args ...interface{}) chromedp.QueryAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		l,ok:=ctx.Value("~log~").(StdLogger)
		if ok{
			l.Fatalf(format,args...)
		}
		return nil
	})
}
func Fatalln(args ...interface{}) chromedp.QueryAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		l,ok:=ctx.Value("~log~").(StdLogger)
		if ok{
			l.Fatalln(args...)
		}
		return nil
	})
}
func Fatal(args ...interface{}) chromedp.QueryAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		l,ok:=ctx.Value("~log~").(StdLogger)
		if ok{
			l.Fatal(args...)
		}
		return nil
	})
}
func Printf(format string, args ...interface{}) chromedp.QueryAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		l,ok:=ctx.Value("~log~").(StdLogger)
		if ok{
			l.Printf(format,args...)
		}
		return nil
	})
}
func Println(args ...interface{}) chromedp.QueryAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		l,ok:=ctx.Value("~log~").(StdLogger)
		if ok{
			l.Println(args...)
		}
		return nil
	})
}
func Print(args ...interface{}) chromedp.QueryAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		l,ok:=ctx.Value("~log~").(StdLogger)
		if ok{
			l.Print(args...)
		}
		return nil
	})
}
func WaitForJS(expression string, res interface{}, opts ...chromedp.PollOption) chromedp.ActionFunc {

	return func(ctx context.Context) error {
		tm:=time.NewTicker(time.Second)
		defer tm.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-tm.C:
				err:=chromedp.Poll(expression,res,opts...).Do(ctx)
				if err!=nil{
					if strings.Contains(err.Error(),"Cannot find context with specified id (-32000)")==false&&strings.Contains(err.Error(),"Execution context was destroyed. (-32000)"){
						return err
					}
				}else{
					if res!=nil{
						return nil
					}
				}
			}
		}
		return nil
	}
}
func GetClientSize(sel string, viewPort *page.Viewport, opts ...chromedp.QueryOption)chromedp.QueryAction{
	getClientRectJS:= `function getClientRect() {
		const e = this.getBoundingClientRect(),
		t = this.ownerDocument.documentElement.getBoundingClientRect();
		return {
			x: e.left - t.left,
			y: e.top - t.top,
			width: e.width,
			height: e.height,
		};
	}`
	return chromedp.QueryAfter(sel, func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}
		if err := CallFunctionOnNode(ctx, nodes[0], getClientRectJS, viewPort); err != nil {
			return err
		}
		x, y := math.Round(viewPort.X), math.Round(viewPort.Y)
		viewPort.Width, viewPort.Height = math.Round(viewPort.Width+viewPort.X-x), math.Round(viewPort.Height+viewPort.Y-y)
		viewPort.X, viewPort.Y = x, y
		viewPort.Scale=1

		return nil

	},append(opts, chromedp.NodeVisible)...)
}

func CallFunctionOnNode(ctx context.Context, node *cdp.Node, function string, res interface{}, args ...interface{}) error {
	r, err := dom.ResolveNode().WithNodeID(node.NodeID).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.CallFunctionOn(function, &res,
		func(p *runtime.CallFunctionOnParams) *runtime.CallFunctionOnParams {
			return p.WithObjectID(r.ObjectID)
		},
		args...,
	).Do(ctx)

	if err != nil {
		return err
	}

	// Try to release the remote object.
	// It will fail if the page is navigated or closed,
	// and it's okay to ignore the error in this case.
	_ = runtime.ReleaseObject(r.ObjectID).Do(ctx)

	return nil
}