package zlog

type LogConf struct {
	Director      string `default:"log" yaml:"director"`
	Level         string `default:"debug" yaml:"level"`
	ShowLine      bool   `default:"true" yaml:"showLine"`
	StacktraceKey string `default:"stacktrace" yaml:"stacktraceKey"`
	EncodeLevel   string `default:"LowercaseColorLevelEncoder" yaml:"encodeLevel"`
	Format        string `default:"console" yaml:"format"`
	Prefix        string `default:"" yaml:"prefix"`
	LogInConsole  bool   `default:"true" yaml:"logInConsole"`
	LinkName      string `default:"latest_log" yaml:"linkName"`
	Trimmed       bool   `default:"true" yaml:"trimmed"`
}
