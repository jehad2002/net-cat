// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// 	"strings"
// 	"time"
// )

// type Server struct {
// 	Server         net.Listener
// 	Connections    map[net.Conn]string
// 	UsedNames      map[string]bool
// 	MaxConnections int
// 	AllMessages    []string
// 	mutex          chan struct{}
// 	ShutdownChan   chan bool
// }

// const (
// 	WelcomeMessage  = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    .       | ' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     -'       --'\n[ENTER YOUR NAME]: "
// 	PatternSending  = "[%v][%v]:"
// 	PatternMessage  = "[%v][%s]: %s"
// 	PatternJoinChat = "%s has joined our chat...\n"
// 	PatternLeftChat = "%s has left our chat...\n"
// )

// const (
// 	ModeJoinChat = iota
// 	ModeSendMessage
// 	ModeLeftChat
// )

// const (
// 	TimeDefault = "2006-01-02 15:04:05"
// )

// func getFormattedMessage(serv *Server, conn net.Conn, message string, mode int) string {
// 	serv.mutex <- struct{}{}
// 	name := serv.Connections[conn]
// 	<-serv.mutex
// 	switch mode {
// 	case ModeSendMessage:
// 		if message == "\n" {
// 			return ""
// 		}
// 		currentTime := time.Now().Format(TimeDefault)
// 		message = fmt.Sprintf(PatternMessage, currentTime, name, message)
// 	case ModeJoinChat:
// 		message = fmt.Sprintf(PatternJoinChat, name)
// 	case ModeLeftChat:
// 		message = fmt.Sprintf(PatternLeftChat, name)
// 	}
// 	return message
// }

// func (s *Server) Constructor(port string, maxConn int) {
// 	serv, _ := net.Listen("tcp", port)
// 	s.Server = serv
// 	s.MaxConnections = maxConn
// 	s.Connections = make(map[net.Conn]string, maxConn)
// 	s.UsedNames = make(map[string]bool, maxConn)
// 	s.mutex = make(chan struct{}, 1)
// 	s.ShutdownChan = make(chan bool)
// }

// func (s *Server) CanConnect(conn net.Conn) bool {
// 	s.mutex <- struct{}{}
// 	defer func() { <-s.mutex }()
// 	return !(s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections)
// }

// func (s *Server) ConnectMessenger(conn net.Conn) {
// 	if !s.CanConnect(conn) {
// 		fmt.Fprint(conn, "The room is full, please try again later...")
// 		conn.Close()
// 		return
// 	}
// 	fmt.Fprint(conn, WelcomeMessage)
// 	name, _ := bufio.NewReader(conn).ReadString('\n')
// 	name = strings.TrimSpace(name)
// 	if !s.addConnection(conn, name) {
// 		conn.Close()
// 		return
// 	}
// 	s.startChatting(conn)
// 	s.removeConnection(conn)
// }

// func (s *Server) CloseServer() {
// 	log.Println("Closing Server")
// 	s.mutex <- struct{}{}
// 	for conn := range s.Connections {
// 		fmt.Fprint(conn, "\nServer Was Closed!\n")
// 		conn.Close()
// 	}
// 	<-s.mutex
// 	s.Server.Close()
// 	log.Println("Server Closed")
// }

// func (s *Server) startChatting(conn net.Conn) {
// 	s.loadMessages(conn)
// 	message := getFormattedMessage(s, conn, "", ModeJoinChat)
// 	s.sendMessage(conn, message)
// 	for {
// 		message, err := bufio.NewReader(conn).ReadString('\n')
// 		if err != nil {
// 			break
// 		}

// 		if strings.HasPrefix(message, "/name ") {
// 			newName := strings.TrimSpace(strings.TrimPrefix(message, "/name "))
// 			s.changeUserName(conn, newName)
// 			continue
// 		}

// 		message = getFormattedMessage(s, conn, message, ModeSendMessage)
// 		s.sendMessage(conn, message)
// 		s.saveMessage(message)
// 	}
// 	message = getFormattedMessage(s, conn, "", ModeLeftChat)
// 	s.sendMessage(conn, message)
// }

// func (s *Server) sendMessage(conn net.Conn, message string) {
// 	if message == "" {
// 		fmt.Fprintf(conn, PatternSending, time.Now().Format(TimeDefault), s.Connections[conn])
// 		return
// 	}
// 	currentTime := time.Now().Format(TimeDefault)
// 	s.mutex <- struct{}{}
// 	for con := range s.Connections {
// 		if con != conn {
// 			fmt.Fprint(con, message)
// 		}
// 		fmt.Fprintf(con, PatternSending, currentTime, s.Connections[con])
// 	}
// 	<-s.mutex
// }

// func (s *Server) loadMessages(conn net.Conn) {
// 	for _, message := range s.AllMessages {
// 		fmt.Fprint(conn, message)
// 	}
// }

// func (s *Server) saveMessage(message string) {
// 	s.mutex <- struct{}{}
// 	s.AllMessages = append(s.AllMessages, message)
// 	<-s.mutex
// }

// func (s *Server) addConnection(conn net.Conn, name string) bool {
// 	s.mutex <- struct{}{}
// 	defer func() { <-s.mutex }()
// 	if name == "" || s.UsedNames[name] || (s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections) {
// 		return false
// 	}
// 	s.UsedNames[name] = true
// 	s.Connections[conn] = name
// 	log.Printf("%s has joined our chat...\n", name)
// 	return true
// }

// func (s *Server) removeConnection(conn net.Conn) {
// 	s.mutex <- struct{}{}
// 	name := s.Connections[conn]
// 	delete(s.UsedNames, name)
// 	delete(s.Connections, conn)
// 	log.Printf("%s has left our chat...\n", name)
// 	<-s.mutex
// }

// func (s *Server) changeUserName(conn net.Conn, newName string) {
// 	s.mutex <- struct{}{}
// 	defer func() { <-s.mutex }()
// 	if newName == "" || s.UsedNames[newName] {
// 		fmt.Fprint(conn, "Name is either empty or already in use.\n")
// 		return
// 	}
// 	oldName := s.Connections[conn]
// 	delete(s.UsedNames, oldName)
// 	s.UsedNames[newName] = true
// 	s.Connections[conn] = newName
// 	log.Printf("%s changed their name to %s\n", oldName, newName)
// 	message := fmt.Sprintf("%s changed their name to %s\n", oldName, newName)
// 	s.sendMessage(conn, message)
// }

// func (s *Server) WaitForExitCommand() {
// 	reader := bufio.NewReader(os.Stdin)
// 	for {
// 		command, _ := reader.ReadString('\n')
// 		if strings.TrimSpace(command) == "exit" {
// 			close(s.ShutdownChan)
// 			return
// 		}
// 	}
// }

// func main() {
// 	port := GetPort()
// 	server := &Server{}
// 	server.Constructor(port, 0)
// 	fmt.Printf("Listening on the port %v\n", port)

// 	go server.WaitForExitCommand()

// 	go func() {
// 		for {
// 			conn, err := server.Server.Accept()
// 			if err != nil {
// 				break
// 			}
// 			go server.ConnectMessenger(conn)
// 		}
// 	}()

// 	<-server.ShutdownChan
// 	server.CloseServer()
// }

// func GetPort() string {
// 	args := os.Args
// 	if len(args) < 2 {
// 		fmt.Println("[USAGE]: ./TCPChat $port")
// 		os.Exit(0)
// 	}
// 	return ":" + args[1]
// }

//================================change name====================================

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Server struct {
	Server         net.Listener
	Connections    map[net.Conn]string
	UsedNames      map[string]bool
	MaxConnections int
	AllMessages    []string
	mutex          chan struct{}
	ShutdownChan   chan bool
}

const (
	WelcomeMessage  = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    .       | ' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     -'       --'\n[ENTER YOUR NAME]: "
	PatternSending  = "[%v][%v]:"
	PatternMessage  = "[%v][%s]: %s"
	PatternJoinChat = "%s has joined our chat...\n"
	PatternLeftChat = "%s has left our chat...\n"
)

const (
	ModeJoinChat = iota
	ModeSendMessage
	ModeLeftChat
)

const (
	TimeDefault = "2006-01-02 15:04:05"
)

func getFormattedMessage(serv *Server, conn net.Conn, message string, mode int) string {
	serv.mutex <- struct{}{}
	name := serv.Connections[conn]
	<-serv.mutex
	switch mode {
	case ModeSendMessage:
		if message == "\n" {
			return ""
		}
		currentTime := time.Now().Format(TimeDefault)
		message = fmt.Sprintf(PatternMessage, currentTime, name, message)
	case ModeJoinChat:
		message = fmt.Sprintf(PatternJoinChat, name)
	case ModeLeftChat:
		message = fmt.Sprintf(PatternLeftChat, name)
	}
	return message
}

func (s *Server) Constructor(port string, maxConn int) {
	serv, _ := net.Listen("tcp", port)
	s.Server = serv
	s.MaxConnections = maxConn
	s.Connections = make(map[net.Conn]string, maxConn)
	s.UsedNames = make(map[string]bool, maxConn)
	s.mutex = make(chan struct{}, 1)
	s.ShutdownChan = make(chan bool)
}

func (s *Server) CanConnect(conn net.Conn) bool {
	s.mutex <- struct{}{}
	defer func() { <-s.mutex }()
	return !(s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections)
}

func (s *Server) ConnectMessenger(conn net.Conn) {
	if !s.CanConnect(conn) {
		fmt.Fprint(conn, "The room is full, please try again later...")
		conn.Close()
		return
	}
	fmt.Fprint(conn, WelcomeMessage)
	name, _ := bufio.NewReader(conn).ReadString('\n')
	name = strings.TrimSpace(name)
	if !s.addConnection(conn, name) {
		conn.Close()
		return
	}
	s.startChatting(conn)
	s.removeConnection(conn)
}

func (s *Server) CloseServer() {
	log.Println("Closing Server")
	s.mutex <- struct{}{}
	for conn := range s.Connections {
		fmt.Fprint(conn, "\nServer Was Closed!\n")
		conn.Close()
	}
	<-s.mutex
	s.Server.Close()
	log.Println("Server Closed")
}

func (s *Server) startChatting(conn net.Conn) {
	s.loadMessages(conn)
	message := getFormattedMessage(s, conn, "", ModeJoinChat)
	s.sendMessage(conn, message)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}

		if strings.HasPrefix(message, "/name ") {
			newName := strings.TrimSpace(strings.TrimPrefix(message, "/name "))
			s.changeUserName(conn, newName)
			continue
		}

		message = getFormattedMessage(s, conn, message, ModeSendMessage)
		s.sendMessage(conn, message)
		s.saveMessage(message)
	}
	message = getFormattedMessage(s, conn, "", ModeLeftChat)
	s.sendMessage(conn, message)
}

func (s *Server) sendMessage(conn net.Conn, message string) {
	if message == "" {
		fmt.Fprintf(conn, PatternSending, time.Now().Format(TimeDefault), s.Connections[conn])
		return
	}
	currentTime := time.Now().Format(TimeDefault)
	s.mutex <- struct{}{}
	for con := range s.Connections {
		if con != conn {
			fmt.Fprint(con, message)
		}
		fmt.Fprintf(con, PatternSending, currentTime, s.Connections[con])
	}
	<-s.mutex
}

func (s *Server) loadMessages(conn net.Conn) {
	s.mutex <- struct{}{}
	defer func() { <-s.mutex }()
	for _, message := range s.AllMessages {
		fmt.Fprint(conn, message)
	}
}

func (s *Server) saveMessage(message string) {
	s.mutex <- struct{}{}
	s.AllMessages = append(s.AllMessages, message)
	<-s.mutex
}

func (s *Server) addConnection(conn net.Conn, name string) bool {
	s.mutex <- struct{}{}
	defer func() { <-s.mutex }()
	if name == "" || s.UsedNames[name] || (s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections) {
		return false
	}
	s.UsedNames[name] = true
	s.Connections[conn] = name
	log.Printf("%s has joined our chat...\n", name)
	return true
}

func (s *Server) removeConnection(conn net.Conn) {
	s.mutex <- struct{}{}
	name := s.Connections[conn]
	delete(s.UsedNames, name)
	delete(s.Connections, conn)
	log.Printf("%s has left our chat...\n", name)
	<-s.mutex
}

func (s *Server) changeUserName(conn net.Conn, newName string) {
	s.mutex <- struct{}{}
	if newName == "" || s.UsedNames[newName] {
		<-s.mutex
		fmt.Fprint(conn, "Name is either empty or already in use.\n")
		return
	}
	oldName := s.Connections[conn]
	s.UsedNames[newName] = true
	s.Connections[conn] = newName
	delete(s.UsedNames, oldName)
	<-s.mutex
	log.Printf("%s changed their name to %s\n", oldName, newName)
	message := fmt.Sprintf("%s changed their name to %s\n", oldName, newName)
	s.sendMessage(conn, message)
}

func (s *Server) WaitForExitCommand() {
	reader := bufio.NewReader(os.Stdin)
	for {
		command, _ := reader.ReadString('\n')
		if strings.TrimSpace(command) == "exit" {
			close(s.ShutdownChan)
			return
		}
	}
}

func main() {
	port := GetPort()
	server := &Server{}
	server.Constructor(port, 0)
	fmt.Printf("Listening on the port %v\n", port)

	go server.WaitForExitCommand()

	go func() {
		for {
			conn, err := server.Server.Accept()
			if err != nil {
				break
			}
			go server.ConnectMessenger(conn)
		}
	}()

	<-server.ShutdownChan
	server.CloseServer()
}

func GetPort() string {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	}
	return ":" + args[1]
}
