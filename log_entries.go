package logger

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// msgQueue is message queue channel
type msgQueue struct {
	messagesToSend chan *bytes.Buffer
}

// Enqueue will enqueue the message
func (m *msgQueue) Enqueue(msg *bytes.Buffer) {
	m.messagesToSend <- msg
}

// PushFront push the message to front
func (m *msgQueue) PushFront(msg *bytes.Buffer) {
	messages := []*bytes.Buffer{msg}
	for {
		select {
		case msg := <-m.messagesToSend:
			messages = append(messages, msg)
		default:
			for _, msg := range messages {
				m.messagesToSend <- msg
			}
			return
		}
	}
}

// LogClient configuration
type LogClient struct {
	conn       *net.TCPConn
	endpoint   string
	messages   msgQueue
	port       string
	retryDelay time.Duration
	token      string
}

// NewLogEntriesClient new client
func NewLogEntriesClient(token, endpoint, port string) (*LogClient, error) {
	l := &LogClient{
		endpoint:   endpoint,
		port:       port,
		retryDelay: RetryDelay,
		token:      token,
	}
	l.messages.messagesToSend = make(chan *bytes.Buffer, 1000)

	if err := l.Connect(); err != nil {
		return l, err
	}

	return l, nil
}

// Connect will connect to Log Entries
func (l *LogClient) Connect() error {
	if l.conn != nil {
		_ = l.conn.Close() // close the connection, don't care about the error
	}
	l.conn = nil

	addr, err := net.ResolveTCPAddr("tcp", l.endpoint+":"+l.port)
	if err != nil {
		return err
	}

	var conn *net.TCPConn
	if conn, err = net.DialTCP("tcp", nil, addr); err != nil {
		l.retryDelay *= 2
		if l.retryDelay > MaxRetryDelay {
			l.retryDelay = MaxRetryDelay
		}
		return err
	}

	_ = conn.SetNoDelay(true)
	// if err != nil {
	//	return err
	// }
	_ = conn.SetKeepAlive(true)
	// if err != nil {
	//	return err
	// }

	l.conn = conn
	l.retryDelay = RetryDelay

	return nil
}

// ProcessQueue process the queue
func (l *LogClient) ProcessQueue() {
	for msg := range l.messages.messagesToSend {
		if l.conn == nil {
			l.messages.PushFront(msg)
			time.Sleep(l.retryDelay)
			if err := l.Connect(); err != nil {
				log.Println("failed reconnecting to log provider", err)
				continue
			}
		}
		if _, err := l.conn.Write(msg.Bytes()); err != nil {
			l.messages.PushFront(msg)
			log.Println("failed to write to log provider", err)
			time.Sleep(l.retryDelay)
			if err = l.Connect(); err != nil {
				log.Println("failed reconnecting to log provider after failing to write", err)
				continue
			}
		}
	}
}

// Panic overloads built-in method
func (l *LogClient) Panic(v ...interface{}) {
	var buff bytes.Buffer
	buff.WriteString(l.token)
	buff.WriteByte(' ')
	fmt.Fprintln(&buff, v...)
	_ = l.sendOne(&buff)
	time.Sleep(2 * time.Second)
	os.Exit(1)
}

// Panicln overloads built-in method
func (l *LogClient) Panicln(v ...interface{}) {
	var buff bytes.Buffer
	buff.WriteString(l.token)
	buff.WriteByte(' ')
	fmt.Fprintln(&buff, v...)
	_ = l.sendOne(&buff)
	time.Sleep(2 * time.Second)
	os.Exit(1)
}

// Panicf overloads built-in method
func (l *LogClient) Panicf(format string, v ...interface{}) {
	var buff bytes.Buffer
	buff.WriteString(l.token)
	buff.WriteByte(' ')
	fmt.Fprintf(&buff, format, v...)
	_ = l.sendOne(&buff)
	time.Sleep(2 * time.Second)
	os.Exit(1)
}

// Print overloads built-in method
func (l *LogClient) Print(v ...interface{}) {
	l.write(fmt.Sprintln(v...))
}

// Println overloads built-in method
func (l *LogClient) Println(v ...interface{}) {
	l.write(fmt.Sprintln(v...))
}

// Printf overloads built-in method
func (l *LogClient) Printf(format string, v ...interface{}) {
	l.write(fmt.Sprintf(format, v...))
}

// Fatal overloads built-in method
func (l *LogClient) Fatal(v ...interface{}) {
	var buff bytes.Buffer
	buff.WriteString(l.token)
	buff.WriteByte(' ')
	fmt.Fprintln(&buff, v...)
	_ = l.sendOne(&buff)
	time.Sleep(2 * time.Second)
	os.Exit(1)
}

// Fatalln overloads built-in method
func (l *LogClient) Fatalln(v ...interface{}) {
	var buff bytes.Buffer
	buff.WriteString(l.token)
	buff.WriteByte(' ')
	fmt.Fprintln(&buff, v...)
	_ = l.sendOne(&buff)
	time.Sleep(2 * time.Second)
	os.Exit(1)
}

// Fatalf overloads built-in method
func (l *LogClient) Fatalf(format string, v ...interface{}) {
	var buff bytes.Buffer
	buff.WriteString(l.token)
	buff.WriteByte(' ')
	fmt.Fprintf(&buff, format, v...)
	_ = l.sendOne(&buff)
	time.Sleep(2 * time.Second)
	os.Exit(1)
}

// write will write the data to the que
func (l *LogClient) write(data string) {
	var buff bytes.Buffer
	buff.WriteString(l.token)
	buff.WriteByte(' ')
	buff.WriteString(data)
	l.messages.Enqueue(&buff)
}

// sendOne sends one log
func (l *LogClient) sendOne(msg *bytes.Buffer) (err error) {
	if l.conn == nil {
		if err = l.Connect(); err != nil {
			log.Println(msg.String())
			log.Println("failed reconnecting to log provider", err)

			return err
		}
	}
	if _, err = l.conn.Write(msg.Bytes()); err != nil {
		log.Println(msg.String())
		log.Println("failed to write to log provider", err)
		return err
	}
	return nil
}
