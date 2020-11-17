package logger

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testToken = "token"

// TestNewLogEntriesClient will test the NewLogEntriesClient() method
func TestNewLogEntriesClient(t *testing.T) {
	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, LogEntriesPort, client.port)
	assert.Equal(t, LogEntriesURL, client.endpoint)
	assert.Equal(t, testToken, client.token)

	client, err = NewLogEntriesClient(testToken, LogEntriesURL, "101010")
	assert.Error(t, err)
	assert.NotNil(t, client)

	client, err = NewLogEntriesClient(testToken, "http://badurl.com", LogEntriesPort)
	assert.Error(t, err)
	assert.NotNil(t, client)

	// Double open
	client, err = NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

// TestMsgQueue_Enqueue will test the Enqueue() method
func TestMsgQueue_Enqueue(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	var buff bytes.Buffer
	buff.WriteString(testToken)
	buff.WriteByte(' ')
	buff.WriteString("test")
	client.messages.Enqueue(&buff)

	assert.Equal(t, 1, len(client.messages.messagesToSend))

	for x := range client.messages.messagesToSend {
		assert.Equal(t, testToken+" test", x.String())
		close(client.messages.messagesToSend)
	}
}

// TestMsgQueue_PushFront will test the PushFront() method
func TestMsgQueue_PushFront(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	assert.NoError(t, err)
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

	assert.Equal(t, 2, len(client.messages.messagesToSend))

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

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	go client.ProcessQueue()

	var buff bytes.Buffer
	buff.WriteString(testToken)
	buff.WriteByte(' ')
	buff.WriteString("test")
	client.messages.Enqueue(&buff)

	time.Sleep(3 * time.Second)

	assert.Equal(t, 0, len(client.messages.messagesToSend))
}

// TestLogClient_Println will test the Println() method
func TestLogClient_Println(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client.Println("test", "this")

	assert.Equal(t, 1, len(client.messages.messagesToSend))
}

// TestLogClient_Printf will test the Printf() method
func TestLogClient_Printf(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client.Printf("test %d", 1)

	assert.Equal(t, 1, len(client.messages.messagesToSend))

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

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		client.Fatalf("test %d", 1)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLogClient_Fatalf")
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestLogClient_Fatalln will test the Fatalln() method
func TestLogClient_Fatalln(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		client.Fatalln("test exit")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLogClient_Fatalln")
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
