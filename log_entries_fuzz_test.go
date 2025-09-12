package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func FuzzNewLogEntriesClient(f *testing.F) {
	f.Add("token123", "example.com", "443")
	f.Add("", "localhost", "80")
	f.Add("tk", "", "10000")
	f.Add("very-long-token-"+strings.Repeat("x", 1000), "host.com", "8080")
	f.Add("token", "host", "")
	f.Add("token:with:colons", "127.0.0.1", "22")
	f.Add("unicode-token-test", "test.com", "80")

	f.Fuzz(func(t *testing.T, token, endpoint, port string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("NewLogEntriesClient panicked: token=%q, endpoint=%q, port=%q, panic: %v",
					token, endpoint, port, r)
			}
		}()

		client, err := NewLogEntriesClient(token, endpoint, port)

		if client == nil {
			t.Error("NewLogEntriesClient should return non-nil client even on error")
			return
		}

		if client.token != token {
			t.Errorf("Client token mismatch: expected %q, got %q", token, client.token)
		}
		if client.endpoint != endpoint {
			t.Errorf("Client endpoint mismatch: expected %q, got %q", endpoint, client.endpoint)
		}
		if client.port != port {
			t.Errorf("Client port mismatch: expected %q, got %q", port, client.port)
		}

		if err == nil && client.retryDelay != RetryDelay {
			t.Errorf("Client retryDelay should be initialized to RetryDelay (%v), got %v",
				RetryDelay, client.retryDelay)
		}

		if client.messages.messagesToSend == nil {
			t.Error("Client message queue should be initialized")
		}

		if err == nil && (endpoint == "" || port == "") {
			t.Error("Expected error for empty endpoint or port, but got nil")
		}
	})
}

func FuzzLogClient_write(f *testing.F) {
	f.Add("test message")
	f.Add("")
	f.Add("message with\nnewlines\tand\ttabs")
	f.Add("unicode: test rocket cafe")
	f.Add(strings.Repeat("long message ", 1000))
	f.Add("special chars: !@#$%^&*(){}[]|\\:;\"'<>?,./")
	f.Add("null byte: \x00 embedded")

	f.Fuzz(func(t *testing.T, data string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("LogClient.write panicked with data %q: %v", data, r)
			}
		}()

		client := &LogClient{
			token:      "test-token",
			retryDelay: time.Millisecond,
		}
		client.messages.messagesToSend = make(chan *bytes.Buffer, 1000)

		client.write(data)

		select {
		case msg := <-client.messages.messagesToSend:
			msgStr := msg.String()
			if !strings.HasPrefix(msgStr, "test-token ") {
				t.Errorf("Message should start with token, got: %q", msgStr)
			}

			expectedData := "test-token " + data
			if msgStr != expectedData {
				t.Errorf("Message mismatch: expected %q, got %q", expectedData, msgStr)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Expected message to be queued within timeout")
		}
	})
}

func FuzzMsgQueue_PushFront(f *testing.F) {
	f.Add("test1", "test2", "test3")
	f.Add("", "msg", "")
	f.Add("a", "", "c")
	f.Add("long"+strings.Repeat("x", 1000), "msg2", "msg3")

	f.Fuzz(func(t *testing.T, msg1, msg2, msg3 string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("msgQueue.PushFront panicked: msg1=%q, msg2=%q, msg3=%q, panic: %v",
					msg1, msg2, msg3, r)
			}
		}()

		queue := &msgQueue{
			messagesToSend: make(chan *bytes.Buffer, 100),
		}

		buf2 := bytes.NewBufferString(msg2)
		buf3 := bytes.NewBufferString(msg3)
		queue.messagesToSend <- buf2
		queue.messagesToSend <- buf3

		buf1 := bytes.NewBufferString(msg1)
		queue.PushFront(buf1)

		select {
		case receivedMsg := <-queue.messagesToSend:
			if receivedMsg.String() != msg1 {
				t.Errorf("First message should be %q, got %q", msg1, receivedMsg.String())
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Expected first message within timeout")
		}

		select {
		case receivedMsg := <-queue.messagesToSend:
			if receivedMsg.String() != msg2 {
				t.Errorf("Second message should be %q, got %q", msg2, receivedMsg.String())
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Expected second message within timeout")
		}

		select {
		case receivedMsg := <-queue.messagesToSend:
			if receivedMsg.String() != msg3 {
				t.Errorf("Third message should be %q, got %q", msg3, receivedMsg.String())
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Expected third message within timeout")
		}
	})
}

func FuzzMsgQueue_Enqueue(f *testing.F) {
	f.Add("test message")
	f.Add("")
	f.Add("unicode: test")
	f.Add(strings.Repeat("x", 10000))
	f.Add("special\nchars\ttabs")

	f.Fuzz(func(t *testing.T, message string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("msgQueue.Enqueue panicked with message %q: %v", message, r)
			}
		}()

		queue := &msgQueue{
			messagesToSend: make(chan *bytes.Buffer, 100),
		}

		buf := bytes.NewBufferString(message)
		queue.Enqueue(buf)

		select {
		case receivedMsg := <-queue.messagesToSend:
			if receivedMsg.String() != message {
				t.Errorf("Message mismatch: expected %q, got %q", message, receivedMsg.String())
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Expected message to be queued within timeout")
		}
	})
}

func FuzzLogClient_sendOne(f *testing.F) {
	f.Add("test message")
	f.Add("")
	f.Add("message\nwith\nnewlines")
	f.Add(strings.Repeat("x", 10000))

	f.Fuzz(func(t *testing.T, message string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("LogClient.sendOne panicked with message %q: %v", message, r)
			}
		}()

		client := &LogClient{
			endpoint:   "nonexistent.invalid",
			port:       "12345",
			token:      "test-token",
			retryDelay: time.Millisecond,
		}
		client.messages.messagesToSend = make(chan *bytes.Buffer, 100)

		buf := bytes.NewBufferString(message)
		err := client.sendOne(buf)

		if err == nil && client.conn == nil {
			t.Error("sendOne should return error when connection fails and conn is nil")
		}
	})
}
