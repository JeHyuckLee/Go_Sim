package main

func start_server(ip, port string) {
	zctx, _ := zmq.NewContext()

	s, _ := zctx.NewSocket(zmq.PUSH)
	s.Bind("tcp://*:5555")
}
