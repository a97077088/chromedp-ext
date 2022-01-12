package actionExt

import (
	"fmt"
	"testing"
)

func TestNodesSelector(t *testing.T) {
	v := int64(111)
	fmt.Println(NodesSelector(v))
}
