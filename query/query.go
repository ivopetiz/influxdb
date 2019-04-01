package query // import "github.com/ivopetiz/influxdb/query"

//go:generate tmpl -data=@tmpldata iterator.gen.go.tmpl
//go:generate tmpl -data=@tmpldata point.gen.go.tmpl
//go:generate tmpl -data=@tmpldata functions.gen.go.tmpl

//go:generate protoc --gogo_out=. internal/internal.proto
