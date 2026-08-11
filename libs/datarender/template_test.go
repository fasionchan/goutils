/*
 * Author: wxhuangjuguan
 * Created time: 2025-10-29 14:30:00
 * Last Modified by: fasion
 * Last Modified time: 2025-11-05 22:45:27
 */

package datarender

import (
	"fmt"
	"testing"

	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/types"
)

func TestStringsRender(t *testing.T) {
	strs, err := ParseAndRenderStringsTemplate("{{ .Value }},{{ now.Unix }},,", templatex.TemplateFuncs, StringSeparator(",").SplitTpe, nil, "Body", types.JsonMap{"Value": "v,V"})
	if err != nil {
		t.Errorf("case failed: %s", err)
		return
	}

	fmt.Println(strs)
}

func TestSmartJsonMapRender(t *testing.T) {
	jsonMap, err := ParseAndRenderSmartJsonMapTemplate(`{"User": {"Name": "Test", "Age": 11}}`, templatex.TemplateFuncs, nil)
	if err != nil {
		t.Errorf("case failed: %s", err)
		return
	}

	fmt.Println(jsonMap)
}
