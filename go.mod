module losh

go 1.18

// replace github.com/osteele/liquid => github.com/aisbergg/go-liquid v1.3.1-0.20220424104007-dfd44a5959fb

replace github.com/osteele/liquid => ./lib/vendor/github.com/aisbergg/go-liquid

require (
	github.com/alecthomas/chroma v0.10.0
	github.com/chigopher/pathlib v0.12.0
	github.com/gofiber/fiber/v2 v2.32.0
	github.com/gofiber/template v1.6.27
	github.com/gofiber/utils v0.1.2
	github.com/golang-module/carbon/v2 v2.1.5
	github.com/gookit/config/v2 v2.1.0
	github.com/gookit/gcli/v3 v3.0.1
	github.com/osteele/liquid v1.3.0
	github.com/rotisserie/eris v0.5.4
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/spf13/viper v1.11.0
	go.uber.org/zap v1.21.0
	golang.org/x/text v0.3.7
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gookit/color v1.5.0 // indirect
	github.com/gookit/goutil v0.5.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/klauspost/compress v1.15.2 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/osteele/tuesday v1.0.3 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.0-beta.8 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.36.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220427172511-eb4f295cb31f // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
