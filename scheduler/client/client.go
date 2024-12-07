package client

import "google.golang.org/grpc"

type Client struct {
	conn grpc.ClientConn
}

func NewClient() *Client {

}
