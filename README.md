# gluayaml

Yaml parser for [gopher-lua](https://github.com/yuin/gopher-lua)

```go
package main

import (
	"github.com/yuin/gopher-lua"
	"github.com/kohkimakimoto/gluayaml"
)


func main() {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("yaml", gluayaml.Loader)
	if err := L.DoString(`
local yaml = require("yaml")
local str = [==[
key1: value1
key2:
  - value2
  - value3
]==]

local tb = yaml.parse(str)
print(tb.key1)    -- value1
print(tb.key2[1]) -- value2
print(tb.key2[2]) -- value3
`); err != nil {
		panic(err)
	}
}
```
