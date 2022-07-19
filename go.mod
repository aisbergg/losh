module losh

go 1.18

require (
	github.com/99designs/gqlgen v0.17.2
	github.com/Yamashou/gqlgenc v0.0.6
	github.com/abadojack/whatlanggo v1.0.1
	github.com/aisbergg/go-copier v0.0.0-20220719150757-87f748b479c0
	github.com/aisbergg/go-errors v0.0.0-20220713173946-6fef60b496f0
	github.com/aisbergg/go-pathlib v0.12.1-0.20220717203945-9b2c4f6e078d
	github.com/alecthomas/chroma v0.10.0
	github.com/go-resty/resty/v2 v2.7.0
	github.com/gofiber/fiber/v2 v2.32.0
	github.com/gofiber/template v1.6.27
	github.com/gofiber/utils v0.1.2
	github.com/golang-module/carbon/v2 v2.1.5
	github.com/gookit/color v1.5.1
	github.com/gookit/config/v2 v2.1.0
	github.com/gookit/event v1.0.6
	github.com/gookit/gcli/v3 v3.0.1
	github.com/gookit/validate v1.4.2
	github.com/microcosm-cc/bluemonday v1.0.18
	github.com/osteele/liquid v1.3.0
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/sethvargo/go-retry v0.2.3
	github.com/spf13/afero v1.8.2
	github.com/wk8/go-ordered-map v1.0.0
	go.uber.org/zap v1.21.0
	golang.org/x/text v0.3.7
	golang.org/x/tools v0.1.10
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/gookit/filter v1.1.2 // indirect
	github.com/gookit/goutil v0.5.6 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/klauspost/compress v1.15.2 // indirect
	github.com/matryer/moq v0.2.3 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/osteele/tuesday v1.0.3 // indirect
	github.com/urfave/cli/v2 v2.3.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.36.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/vektah/gqlparser/v2 v2.4.1 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220427172511-eb4f295cb31f // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/sethvargo/go-retry => github.com/aisbergg/go-retry v0.2.4-0.20220608144822-87d55adc1c0c
	github.com/osteele/liquid => github.com/aisbergg/go-liquid v1.3.1-0.20220719205958-08f515a337ef
	github.com/wk8/go-ordered-map => github.com/aisbergg/go-orderedmap v1.0.1-0.20220718132943-bb550a985f23
)
