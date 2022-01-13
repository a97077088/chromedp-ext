package acExt

import (
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"context"
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