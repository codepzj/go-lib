package logger

type LogFormat string

const (
	TEXT LogFormat = "text"
	JSON LogFormat = "json"
)

type LogLevel string

const (
	DEBUG LogLevel = "debug"
	INFO  LogLevel = "info"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
)

type Option struct {
	Format LogFormat     // 日志格式
	Level  LogLevel      // 日志级别
	Output *OutputConfig // 输出配置
}

type OutputConfig struct {
	EnableFile bool
	FilePath   string
	MaxAge     int
}

func defaultOption() *Option {
	return &Option{
		Format: "text",
		Level:  "debug",
		Output: &OutputConfig{
			EnableFile: false,
			FilePath:   "./logs/server.log",
			MaxAge:     3,
		},
	}
}

func withDefault(opt *Option) *Option {
	def := defaultOption()

	if opt == nil {
		return def
	}

	if opt.Format != "" {
		def.Format = opt.Format
	}

	if opt.Level != "" {
		def.Level = opt.Level
	}

	if opt.Output != nil {
		def.Output.EnableFile = opt.Output.EnableFile

		if opt.Output.FilePath != "" {
			def.Output.FilePath = opt.Output.FilePath
		}

		if opt.Output.MaxAge > 0 {
			def.Output.MaxAge = opt.Output.MaxAge
		}
	}

	return def
}
