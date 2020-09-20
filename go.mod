module github.com/TiboStev/hugo-wrapper

go 1.14

require (
	bou.ke/monkey v1.0.2
	github.com/golang/mock v1.4.3
	github.com/google/go-github/v31 v31.0.0
	github.com/mholt/archiver/v3 v3.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.8.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.2.2
)

replace github.com/spf13/pflag => github.com/TiboStev/pflag v1.0.6-0.20200918204434-33dec6aac494
