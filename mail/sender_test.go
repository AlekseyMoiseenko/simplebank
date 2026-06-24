package mail

import (
	"testing"

	"github.com/AlekseyMoiseenko/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	// todo: add bad scenario tests
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("../.env")
	require.NoError(t, err)

	sender, err := NewGmailSender(config.SmtpName, config.SmtpUser, config.SmtpPass)
	require.NoError(t, err)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from Aleksey</p>
	`
	to := []string{"aleksey.moiseenko79+simplebanktest@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail("test_template", subject, content, to, attachFiles)
	require.NoError(t, err)
}
