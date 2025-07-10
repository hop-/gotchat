package log

import "os"

type LogBuilder struct {
	logInstance *logger
}

func Configure() *LogBuilder {
	return &LogBuilder{
		&logger{},
	}
}

func (b *LogBuilder) Init() error {
	err := b.logInstance.init()
	if err != nil {
		return err
	}

	logInstance = b.logInstance
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
