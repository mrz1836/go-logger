package logger

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
	"time"
)

const testToken = "token"

// TestNewLogEntriesClient will test the NewLogEntriesClient() method
func TestNewLogEntriesClient(t *testing.T) {
	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

	if client.port != LogEntriesPort {
		t.Fatalf("[%s] expect [%s] result", LogEntriesPort, client.port)
	}

	if client.endpoint != LogEntriesURL {
		t.Fatalf("[%s] expect [%s] result", LogEntriesURL, client.endpoint)
	}

	if client.token != testToken {
		t.Fatalf("[%s] expect [%s] result", testToken, client.token)
	}

	_, err = NewLogEntriesClient("token", LogEntriesURL, "101010")
	if err == nil {
		t.Fatalf("error should have occurred")
	}

	_, err = NewLogEntriesClient("token", "http://badurl.com", LogEntriesPort)
	if err == nil {
		t.Fatalf("error should have occurred")
	}

	// Double open
	_, err = NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}
}

// TestMsgQueue_Enqueue will test the Enqueue() method
func TestMsgQueue_Enqueue(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

	var buff bytes.Buffer
	buff.WriteString(testToken)
	buff.WriteByte(' ')
	buff.WriteString("test")
	client.messages.Enqueue(&buff)

	if len(client.messages.messagesToSend) == 0 {
		t.Fatalf("missing messages to send: %v", client.messages.messagesToSend)
	}

	for x := range client.messages.messagesToSend {
		if x.String() != testToken+" test" {
			t.Fatalf("[%s] expect [%s] result", testToken+" test", x.String())
		}
		close(client.messages.messagesToSend)
	}
}

// TestMsgQueue_PushFront will test the PushFront() method
func TestMsgQueue_PushFront(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

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

	if len(client.messages.messagesToSend) == 0 {
		t.Fatalf("missing messages to send: %v", client.messages.messagesToSend)
	}

	go func() {
		close(client.messages.messagesToSend)
	}()

	finalString := ""
	for x := range client.messages.messagesToSend {
		finalString += x.String()
	}

	if finalString != testToken+" first"+testToken+" test" {
		t.Fatalf("[%s] expect [%s] result", testToken+" first"+testToken+" test", finalString)
	}
}

// TestLogClient_ProcessQueue will test the ProcessQueue() method
func TestLogClient_ProcessQueue(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

	go client.ProcessQueue()

	var buff bytes.Buffer
	buff.WriteString(testToken)
	buff.WriteByte(' ')
	buff.WriteString("test")
	client.messages.Enqueue(&buff)

	time.Sleep(3 * time.Second)

	if len(client.messages.messagesToSend) > 0 {
		t.Fatal("no messages should be in queue", len(client.messages.messagesToSend))
	}
}

// TestLogClient_Println will test the Println() method
func TestLogClient_Println(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

	client.Println("test", "this")

	if len(client.messages.messagesToSend) == 0 {
		t.Fatal("expected message to send")
	}
}

// TestLogClient_Printf will test the Printf() method
func TestLogClient_Printf(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

	client.Printf("test %d", 1)

	if len(client.messages.messagesToSend) == 0 {
		t.Fatal("expected message to send")
	}

	go func() {
		close(client.messages.messagesToSend)
	}()

	finalString := ""
	for x := range client.messages.messagesToSend {
		finalString += x.String()
	}

	if finalString != testToken+" test 1" {
		t.Fatalf("[%s] expect [%s] result", testToken+" test 1", finalString)
	}
}

// TestLogClient_Fatalf will test the Fatalf() method
func TestLogClient_Fatalf(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

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
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

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
