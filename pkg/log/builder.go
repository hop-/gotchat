package log

import "os"

type LogBuilder struct {
	logInstance *logger
	level       int
}

func Configure() *LogBuilder {
	return &LogBuilder{
		&logger{formatLogMessageFn: formatLogMessageWithTime},
		INFO, // Default log level
	}
}

func (b *LogBuilder) Init() error {
	err := b.logInstance.init()
	if err != nil {
		return err
	}

	logInstance = b.logInstance
	level = b.level
	isInitialized = true

	return nil
}

func (b *LogBuilder) InMemory() *LogBuilder {
	b.logInstance.inMemory = true

	return b
}

func (b *LogBuilder) StdOut() *LogBuilder {
	b.logInstance.stdOut = true

	return b
}

func (b *LogBuilder) File(filePath string) *LogBuilder {
	var err error
	b.logInstance.logFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}

	return b
}

func (b *LogBuilder) Level(l int) *LogBuilder {
	if l < FATAL || l > DEBUG {
		return nil // Invalid log level
	}
	b.level = l

	return b
}

func (b *LogBuilder) WithTimestamps() *LogBuilder {
	b.logInstance.formatLogMessageFn = formatLogMessageWithTime

	return b
}

func (b *LogBuilder) WithoutTimestamps() *LogBuilder {
	b.logInstance.formatLogMessageFn = formatLogMessageWithoutTime

	return b
}
