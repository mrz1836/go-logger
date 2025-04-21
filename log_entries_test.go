package logger

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testToken = "token"

// TestNewLogEntriesClient will test the NewLogEntriesClient() method
func TestNewLogEntriesClient(t *testing.T) {
	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	require.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, LogEntriesPort, client.port)
	assert.Equal(t, LogEntriesTestEndpoint, client.endpoint)
	assert.Equal(t, testToken, client.token)

	client, err = NewLogEntriesClient(testToken, LogEntriesTestEndpoint, "101010")
	require.Error(t, err)
	assert.NotNil(t, client)

	client, err = NewLogEntriesClient(testToken, "http://badurl.com", LogEntriesPort)
	require.Error(t, err)
	assert.NotNil(t, client)

	// Double open
	client, err = NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	require.NoError(t, err)
	assert.NotNil(t, client)
}

// TestMsgQueue_Enqueue will test the Enqueue() method
func TestMsgQueue_Enqueue(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	require.NoError(t, err)
	assert.NotNil(t, client)

	var buff bytes.Buffer
	buff.WriteString(testToken)
	buff.WriteByte(' ')
	buff.WriteString("test")
	client.messages.Enqueue(&buff)

	assert.Len(t, client.messages.messagesToSend, 1)

	for x := range client.messages.messagesToSend {
		assert.Equal(t, testToken+" test", x.String())
		close(client.messages.messagesToSend)
	}
}

// TestMsgQueue_PushFront will test the PushFront() method
func TestMsgQueue_PushFront(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	require.NoError(t, err)
	assert.NotNil(t, client)

	var buff bytes.Buffer
	buff.WriteString(testToken)
	buff.WriteByte(' ')
	buff.WriteString("test")
	client.messages.Enqueue(&buff)

	var buff2 bytes.Buffer
	buff2.WriteString(testToken)
	buff2.WriteByte(' ')
	buff2.WriteString("first")
	client.messages.PushFront(&buff2)

	assert.Len(t, client.messages.messagesToSend, 2)

	go func() {
		close(client.messages.messagesToSend)
	}()

	finalString := ""
	for x := range client.messages.messagesToSend {
		finalString += x.String()
	}

	assert.Equal(t, testToken+" first"+testToken+" test", finalString)
}

// TestLogClient_ProcessQueue will test the ProcessQueue() method
func TestLogClient_ProcessQueue(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	require.NoError(t, err)
	assert.NotNil(t, client)

	go client.ProcessQueue()

	var buff bytes.Buffer
	buff.WriteString(testToken)
	buff.WriteByte(' ')
	buff.WriteString("test")
	client.messages.Enqueue(&buff)

	time.Sleep(3 * time.Second)

	assert.Empty(t, client.messages.messagesToSend)
}

// TestLogClient_Println will test the Println() method
func TestLogClient_Println(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	require.NoError(t, err)
	assert.NotNil(t, client)

	client.Println("test", "this")

	assert.Len(t, client.messages.messagesToSend, 1)
}

// TestLogClient_Printf will test the Printf() method
func TestLogClient_Printf(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	require.NoError(t, err)
	assert.NotNil(t, client)

	client.Printf("test %d", 1)

	assert.Len(t, client.messages.messagesToSend, 1)

	go func() {
		close(client.messages.messagesToSend)
	}()

	finalString := ""
	for x := range client.messages.messagesToSend {
		finalString += x.String()
	}

	assert.Equal(t, testToken+" test 1", finalString)
}

// TestLogClient_Fatalf will test the Fatalf() method
func TestLogClient_Fatalf(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	require.NoError(t, err)
	assert.NotNil(t, client)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		client.Fatalf("test %d", 1)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLogClient_Fatalf") //nolint:gosec // G204
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestLogClient_Fatalln will test the Fatalln() method
func TestLogClient_Fatalln(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	require.NoError(t, err)
	assert.NotNil(t, client)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		client.Fatalln("test exit")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLogClient_Fatalln") //nolint:gosec // G204
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
