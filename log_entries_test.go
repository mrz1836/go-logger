package logger

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestNewLogEntriesClient will test the NewLogEntriesClient() method
func TestNewLogEntriesClient(t *testing.T) {
	token := "token"
	client, err := NewLogEntriesClient(token, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

	if client.port != LogEntriesPort {
		t.Fatalf("[%s] expect [%s] result", LogEntriesPort, client.port)
	}

	if client.endpoint != LogEntriesURL {
		t.Fatalf("[%s] expect [%s] result", LogEntriesURL, client.endpoint)
	}

	if client.token != token {
		t.Fatalf("[%s] expect [%s] result", token, client.token)
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
	client, err = NewLogEntriesClient(token, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}
}

// TestMsgQueue_Enqueue will test the Enqueue() method
func TestMsgQueue_Enqueue(t *testing.T) {

	token := "token"
	client, err := NewLogEntriesClient(token, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

	var buff bytes.Buffer
	buff.WriteString(token)
	buff.WriteByte(' ')
	buff.WriteString("test")
	client.messages.Enqueue(&buff)

	if len(client.messages.messagesToSend) == 0 {
		t.Fatalf("missing messages to send: %v", client.messages.messagesToSend)
	}

	for x := range client.messages.messagesToSend {
		if fmt.Sprintf("%s", x) != token+" test" {
			t.Fatalf("[%s] expect [%s] result", token+" test", fmt.Sprintf("%s", x))
		}
		close(client.messages.messagesToSend)
	}
}

// TestMsgQueue_PushFront will test the PushFront() method
func TestMsgQueue_PushFront(t *testing.T) {

	token := "token"
	client, err := NewLogEntriesClient(token, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

	var buff bytes.Buffer
	buff.WriteString(token)
	buff.WriteByte(' ')
	buff.WriteString("test")
	client.messages.Enqueue(&buff)

	var buff2 bytes.Buffer
	buff2.WriteString(token)
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

	if finalString != token+" first"+token+" test" {
		t.Fatalf("[%s] expect [%s] result", token+" first"+token+" test", finalString)
	}
}

// TestLogClient_ProcessQueue will test the ProcessQueue() method
func TestLogClient_ProcessQueue(t *testing.T) {
	token := "token"
	client, err := NewLogEntriesClient(token, LogEntriesURL, LogEntriesPort)
	if err != nil {
		t.Fatalf("error should have not occurred: %s", err.Error())
	}

	go client.ProcessQueue()

	var buff bytes.Buffer
	buff.WriteString(token)
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
	token := "token"
	client, err := NewLogEntriesClient(token, LogEntriesURL, LogEntriesPort)
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
	token := "token"
	client, err := NewLogEntriesClient(token, LogEntriesURL, LogEntriesPort)
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

	if finalString != token+" test 1" {
		t.Fatalf("[%s] expect [%s] result", token+" test 1", finalString)
	}
}

// TestLogClient_Fatalf will test the Fatalf() method
func TestLogClient_Fatalf(t *testing.T) {

	token := "token"
	client, err := NewLogEntriesClient(token, LogEntriesURL, LogEntriesPort)
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

	token := "token"
	client, err := NewLogEntriesClient(token, LogEntriesURL, LogEntriesPort)
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
