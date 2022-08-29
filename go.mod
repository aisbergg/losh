module losh

go 1.18

require (
	github.com/99designs/gqlgen v0.17.2
	github.com/Yamashou/gqlgenc v0.0.6
	github.com/abadojack/whatlanggo v1.0.1
	github.com/aisbergg/go-copier v0.0.0-20220822180427-d38739757c35
	github.com/aisbergg/go-errors v0.0.0-20220825110559-14f02ed2dc16
	github.com/aisbergg/go-pathlib v0.12.1-0.20220826142213-89265f98d94e
	github.com/aisbergg/go-retry v1.0.0
	github.com/aisbergg/go-unidecode v1.1.1
	github.com/alecthomas/chroma v0.10.0
	github.com/alecthomas/participle/v2 v2.0.0-beta.5
	github.com/alecthomas/repr v0.1.0
	github.com/dgraph-io/dgo/v200 v200.0.0-20210401091508-95bfd74de60e
	github.com/fenos/dqlx v0.2.1-0.20210902154011-e8c319a835d3
	github.com/go-resty/resty/v2 v2.7.0
	github.com/gofiber/fiber/v2 v2.32.0
	github.com/gofiber/template v1.6.27
	github.com/gofiber/utils v0.1.2
	github.com/golang-module/carbon/v2 v2.1.5
	github.com/gookit/color v1.5.1
	github.com/gookit/config/v2 v2.1.0
	github.com/gookit/event v1.0.6
	github.com/gookit/gcli/v3 v3.0.5-0.20220809015211-4fa6d74ed8f2
	github.com/gookit/validate v1.4.2
	github.com/microcosm-cc/bluemonday v1.0.19
	github.com/osteele/liquid v1.3.0
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/spf13/afero v1.8.2
	github.com/wk8/go-ordered-map v1.0.0
	go.uber.org/zap v1.21.0
	golang.org/x/text v0.3.7
	golang.org/x/tools v0.1.10
	google.golang.org/grpc v1.48.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/gookit/filter v1.1.2 // indirect
	github.com/gookit/goutil v0.5.7 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/klauspost/compress v1.15.2 // indirect
	github.com/matryer/moq v0.2.3 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/osteele/tuesday v1.0.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/urfave/cli/v2 v2.3.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.36.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/vektah/gqlparser/v2 v2.4.1 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/net v0.0.0-20220812174116-3211cb980234 // indirect
	golang.org/x/sys v0.0.0-20220818161305-2296e01440c6 // indirect
	golang.org/x/term v0.0.0-20220722155259-a9ba230a4035 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/genproto v0.0.0-20220819174105-e9f053255caa // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/fenos/dqlx => github.com/aisbergg/go-dqlx v0.0.0-20220822230403-6ac098906135
	github.com/osteele/liquid => github.com/aisbergg/go-liquid v1.3.1-0.20220719205958-08f515a337ef
	github.com/wk8/go-ordered-map => github.com/aisbergg/go-orderedmap v1.0.1-0.20220812065708-63d414488359
)
