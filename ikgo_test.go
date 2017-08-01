package ikgo

import (
	"fmt"
	"testing"
)

func TestIkgo(t *testing.T) {
	InitDict("/data/baixing/dict/ik/analysis-ik", true)
	text := "测试分词实例"

	segmenter := NewIKSegmenter(text, true)
	for {
		lexme := segmenter.next()
		if lexme == nil {
			break
		}
		fmt.Printf("%+v\n", lexme)
	}
}
