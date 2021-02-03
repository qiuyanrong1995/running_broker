package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"running_broker/internal/logger"
	"strings"
)

type Broker struct {
	BreakPointContinue bool `yaml:"break_point_continue" json:"break_point_continue"`
	Md5Check bool `yaml:"md5_check" json:"md5_check"`
	BandWidth string `yaml:"band_width_percent" json:"band_width_percent"`
	ProgramPath string `yaml:"program_path" json:"program_path"`
	TempPath string `yaml:"temp_path" json:"temp_path"`
	Debug bool `yaml:"debug" json:"debug"`
	HeartBeatConfig *HeartBeatConfig `yaml:"heart_beat_config" json:"heart_beat_config"`
	ProgramConfigs *ProgramConfigs `yaml:"program_configs" json:"program_configs"`
	LogConfig *LogConfig `yaml:"log_config" json:"log_config"`
	logger *logrus.Logger
}

type LogConfig struct {
	Level string `yaml:"level" json:"level"`
	RotationCount int `yaml:"rotation_count" json:"rotation_count"`
	Size int `yaml:"size" json:"size"`
	Age int `yaml:"age" json:"age"`
	JsonFormat bool `yaml:"json_format" json:"json_format"`
	Path string `yaml:"path" json:"path"`
	Name string `yaml:"name" json:"name"`
}

type HeartBeatConfig struct {
	Host string `yaml:"host" json:"host"`
	Cycle int `yaml:"cycle" json:"cycle"`
}

type ProgramConfig struct {
	Host string `yaml:"host" json:"host"`
	Md5 string `yaml:"md5" json:"md5"`
	Version string `yaml:"version" json:"version"`
	Cmd string `yaml:"cmd" json:"cmd"`
}

type ProgramConfigs []ProgramConfig

// 打包编译时修改
var (
	BreakPointContinue = true
	Md5Check = true
	BandWidth = "400kb/s"
	ProgramPath = "programs"
	TempPath = "temp"
	Debug = false
	LogLevel = "info"
	LogRotationCount = 7
	LogSize = 100
	LogAge = 180
	LogJsonFormat = true
	LogPath = "logs"
	LogName = "broker.log"
	HeartBeatCycle = 10
)

func NewBroker() *Broker {
	return &Broker{
		BreakPointContinue: BreakPointContinue,
		Md5Check:           Md5Check,
		BandWidth:   		BandWidth,
		ProgramPath:		ProgramPath,
		TempPath:			TempPath,
		HeartBeatConfig:    &HeartBeatConfig{ Cycle: HeartBeatCycle },
		LogConfig:			&LogConfig{
			Level:         LogLevel,
			RotationCount: LogRotationCount,
			Size:          LogSize,
			Age:           LogAge,
			JsonFormat:    LogJsonFormat,
			Path:          LogPath,
			Name:          LogName,
		},
	}
}

func (b *Broker)Run() {

}

func (b *Broker)LoadConfig(reader io.Reader, typ string) error {
	if typ != "json" && typ != "yaml" {
		return errors.New("config type must be json or yaml")
	}

	var err error
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(reader); err != nil {
		return err
	}

	if typ == "json" {
		if err = json.Unmarshal(buf.Bytes(), b); err != nil {
			return err
		}
	} else {
		if err = yaml.Unmarshal(buf.Bytes(), b); err != nil {
			return err
		}
	}

	return b.Validate()
}

func (b *Broker)Validate() error {
	if b.BandWidth != "" && strings.HasSuffix(b.BandWidth, "/s") {
		return errors.New("BandWidth must end with `/s`")
	}
	if b.ProgramPath == "" {
		return errors.New("ProgramPath can`t be null")
	}
	if b.TempPath == "" {
		return errors.New("TempPath can`t be null")
	}
	if b.HeartBeatConfig.Host == "" {
		return errors.New("HeartBeat Host is null")
	}
	if b.HeartBeatConfig.Cycle == 0 {
		return errors.New("HeartBeat Cycle can`t be zero")
	}
	if b.ProgramConfigs == nil {
		return errors.New("ProgramConfigs was not set")
	}
	if b.LogConfig.Level == "" {
		return errors.New("log level can`t be null")
	}
	if b.LogConfig.Path == "" {
		return errors.New("log path can`t be null")
	}
	if b.LogConfig.Name == "" {
		return errors.New("log name can`t be null")
	}
	if b.LogConfig.RotationCount == 0 {
		return errors.New("log rotation count can`t be zero")
	}
	if b.LogConfig.Size == 0 {
		return errors.New("log size can`t be zero")
	}
	if b.LogConfig.Age == 0 {
		return errors.New("log age can`t be zero")
	}
	for _, p := range *b.ProgramConfigs {
		if p.Host == "" {
			return errors.New("program download Host can`t be null")
		}
		if p.Md5 == "" {
			return errors.New("program md5 can`t be null")
		}
		if p.Version == "" {
			return errors.New("program version can`t be null")
		}
	}
	return nil
}

func (b *Broker)GetLogger() (*logrus.Logger, error) {
	if b.logger == nil {
		return logger.NewLogger(b.LogConfig.Level, b.LogConfig.RotationCount, b.LogConfig.Size, b.LogConfig.Age,
								b.LogConfig.Path, b.LogConfig.Name, b.LogConfig.JsonFormat, b.Debug)
	}
	return b.logger, nil
}