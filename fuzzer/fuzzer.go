package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/siemens/ZapSmtp/_test"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/mail"
	"strings"
	"time"
)

type Fuzzer struct {
	logger  *zap.SugaredLogger
	cores   []zapcore.Core
	closeFn func() error
}

func NewFuzzer(signatureCertPath string, signatureKeyPath string, encryptionCertPaths []string) (*Fuzzer, error) {
	f := &Fuzzer{
		cores: make([]zapcore.Core, 0, 2),
	}

	if err := f.initialize(signatureCertPath,
		signatureKeyPath,
		encryptionCertPaths); err != nil {
		return nil, fmt.Errorf("failed to initialize fuzzer: %w", err)
	}

	return f, nil
}

func (f *Fuzzer) initialize(signatureCertPath string, signatureKeyPath string, encryptionCertPaths []string) error {
	// Initialize console core
	coreConsole, errConsole := initConsoleCore(zapcore.DebugLevel)
	if errConsole != nil {
		return fmt.Errorf("console core initialization failed: %w", errConsole)
	}
	f.cores = append(f.cores, coreConsole)

	// Initialize SMTP core
	coreSmtp, coreCloseFn, errSmtp := initSmtpCore(
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		time.Minute,
		time.Second,
		_test.Server,
		_test.Port,
		_test.Username,
		_test.Password,
		"Example Logger",
		_test.Sender,
		[]mail.Address{_test.Recipient},
		_test.OpensslPath,
		signatureCertPath,
		signatureKeyPath,
		encryptionCertPaths,
		"",
	)

	if errSmtp != nil {
		return fmt.Errorf("SMTP core initialization failed: %w", errSmtp)
	}

	f.cores = append(f.cores, coreSmtp)
	f.closeFn = coreCloseFn

	// Initialize logger
	tee := zapcore.NewTee(f.cores...)
	f.logger = zap.New(tee).Sugar()

	return nil
}

func (f *Fuzzer) Close() error {
	if err := f.logger.Sync(); err != nil {
		return fmt.Errorf("error while syncing logger: %w", err)
	}
	if f.closeFn != nil {
		return f.closeFn()
	}
	return nil
}

func (f *Fuzzer) FuzzLog(data string) {
	f.logger.Errorf(data)
}

func GenerateRandomLog() string {
	var result bytes.Buffer

	result.WriteString("This is a mixed content message with various elements:\n\n")
	result.WriteString("<html><body><h1>HTML Section</h1><p>This is HTML content</p></body></html>\n\n")
	result.WriteString("{\n  \"name\": \"Test Message\",\n  \"type\": \"mixed\",\n  \"items\": [1, 2, 3]\n}\n\n")
	result.WriteString("\n\nUnicode section: 你好, こんにちは, 안녕하세요, Привет\n")
	result.WriteString("Emoji section: 😀 🙌 👍 🎉 🚀 💯 🔥\n")

	return result.String()
}

func GenerateMalformedLog() string {
	return strings.Repeat("\b", 1000) + "Test message"
}

func GenerateRandomLongMessage(length int) string {
	data := make([]byte, length)
	_, err := rand.Read(data)
	if err != nil {
		return ""
	}
	return string(data)
}
