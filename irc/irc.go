package irc

import (
  "fmt"
  "net"
  "strings"
  "strconv"
  "os"
  "bufio"
  "time"
)

type IRCClient interface {
  JOIN(s string);
  PART();
  PRIVMSG(s string);
  MainLoop(messagech chan<- TwitchMessage);
  LoopInput();
}

type ircClient struct {
  conn net.Conn;
  channelName string;
  sentTimes []time.Time;
  privmsgSemaphore chan bool;
}

//Paramters for PRIVMSG to keep Twitch happy
var (
  maxMessagesPerInterval = 3;
  timePerInterval = time.Duration(1)*time.Second;
  maxMessageSize = 450;
)

type TwitchMessage struct {
  Sender string;
  Body string;
}

func NewIRCClient(address string, port int, nick string, pass string) IRCClient {
  //Create a TCP connection to the host on desired port
  conn, error := net.Dial("tcp", address + ":" + strconv.Itoa(port));
  //If there was an error then exit
  if error != nil {
    fmt.Printf("It fucked up!\n");
    os.Exit(1);
  }
  //Create a new ircClient from the connection
  ircClient := &ircClient {
    conn : conn,
  }
  //login to the IRC host with the given nickname and password
  conn.Write([]byte("PASS " + pass + "\n"));
  conn.Write([]byte("NICK " + nick + "\n"));
  //Initialize the sentTimes to make sure we're not going over Twitch limit
  ircClient.sentTimes = make([]time.Time, 0, maxMessagesPerInterval);
  //As well initialize a channel to use as a semaphore to make privmsg thread safe
  ircClient.privmsgSemaphore = make(chan bool, 1);
  //Return the usable ircClient object
  return ircClient;
}

func (c ircClient) read() string {
  //create a buffer of 1000 bytes and read the connection
  buffer := make([]byte, 1000);
  n,_ := c.conn.Read(buffer);
  //stringify the buffer
  data := string(buffer[:n]);
  //while the buffer was filled with the last read, keep reading and appending data
  for n >= 1000 {
    n,_ := c.conn.Read(buffer);
    data += string(buffer[:n]);
  }
  //return a string of everything available on the connection
  return data;
}

func (c ircClient) write(s string) {
  fmt.Fprintf(c.conn, s + "\r\n");
}

func (c *ircClient) JOIN(s string) {
  c.channelName = s;
  c.write("JOIN #" + s);
}

func (c *ircClient) PART() {
  c.write("PART #" + c.channelName);
  c.channelName = "";
}

/*
Twitch has a policy of only accepting 20 messages per 30 seconds (100 for moderators),
so we need to make sure we're not getting ourselves muted by keeping track of sent amount
*/
func (c *ircClient) PRIVMSG(s string) {
  //Truncate messages which are longer than allowed size, otherwize Twitch drops them
  if len(s) > maxMessageSize {
    s = s[:maxMessageSize];
  }
  //Get exlusive access to prevent multiple from threads slipping through while loop
  c.privmsgSemaphore <- true;
  defer func() {<-c.privmsgSemaphore}();
  //While sentTimes is full and the timer on the last isn't up, wait it out
  for len(c.sentTimes) == maxMessagesPerInterval && time.Now().Sub(c.sentTimes[maxMessagesPerInterval - 1]) < timePerInterval {
    //busy wait
  }
  //Remove everything from sentTimes which is no longer counting towards the interval
  for len(c.sentTimes) > 0 && time.Now().Sub(c.sentTimes[len(c.sentTimes) - 1]) > timePerInterval {
    c.sentTimes = c.sentTimes[:len(c.sentTimes)-1];
  }
  //Append the current time and send the message
  c.sentTimes = append(c.sentTimes, time.Now());
  c.write("PRIVMSG #" + c.channelName + " :" + s);
}

func (c ircClient) getMessages() []string {
  //read all lines and split by irc separator
  lines := c.read();
  messages := strings.Split(lines, "\r\n");
  //return the array of strings
  return messages;
}

func (c ircClient) MainLoop(messagech chan<- TwitchMessage) {
  for {
    //read the messages from the IRC connection
    messages := c.getMessages();
    //Loop through the messages and handle themS
    for _, m := range messages {
      if len(m) == 0 {
        //do nothing, just a remnant of strings.Split
      } else if len(m) >= 4 && m[:4] == "PING" {
        //respond to ping messages with similar pong
        c.write("PONG" + m[4:]);
      } else if strings.Index(m,"!") == -1 || strings.Index(m, "GLHF") != -1 || strings.Index(m, "JOIN") != -1 {
        //print messages not coming from other uses
        fmt.Println(m);
      } else {
        //otherwise get the username and body and send it on the message channel
        messagech <- TwitchMessage {
          Sender : m[1:strings.Index(m,"!")],
          Body : m[strings.Index(m[1:],":")+2:],
        }
      }
    }
  }
}

func (c ircClient) LoopInput() {
  //create reader for stdin
  reader := bufio.NewReader(os.Stdin);
  for {
    //On hitting enter, send the message to the channel
    text, _ := reader.ReadString('\n');
    c.PRIVMSG(text);
  }
}