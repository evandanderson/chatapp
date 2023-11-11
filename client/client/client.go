package client

import (
	"bufio"
	"bytes"
	"chatapp/client/consts"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	conn net.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Dial(network, address string) error {
	log.Printf("Attempting to create %s connection to %s\n", network, address)
	conn, err := net.DialTimeout(network, address, consts.Timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Printf("Created %s connection to %s\n", network, address)
	return nil
}

func (c *Client) HandleInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Println(err)
		} else {
			go c.Send(scanner.Bytes())
		}
	}
}

func (c *Client) Receive() {
	defer c.conn.Close()

	buffer := bytes.NewBuffer(make([]byte, 4096))
	for {
		n, err := c.conn.Read(buffer.Bytes())
		if err != nil {
			if err == io.EOF {
				log.Fatalf("Connection from %s closed (EOF)\n", c.conn.RemoteAddr().String())
				break
			}
			fmt.Printf("Failed to read from buffer: %s\n", err)
		}
		fmt.Printf("%s\n", string(buffer.Bytes()[:n]))
	}
}

func (c *Client) Send(payload []byte) {
	if _, err := c.conn.Write(payload); err != nil {
		fmt.Println(err)
	}
}
