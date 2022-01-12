package actionExt

import "github.com/chromedp/cdproto/cdp"

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
