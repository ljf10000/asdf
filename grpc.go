package asdf

import (
	"container/list"
	"errors"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGrpcServer(port string, register func(server *grpc.Server)) error {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer listen.Close()

	server := grpc.NewServer()
	register(server)
	reflection.Register(server)

	Log.Info("run grpc server %s", port)

	return server.Serve(listen)
}

func GrpcCall(server string, call func(conn *grpc.ClientConn) error) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if nil != err {
		return err
	}

	err = call(conn)

	conn.Close()

	return err
}

type GrpcClientPool struct {
	server string

	q    *list.List
	lock *AccessLock
}

func NewGrpcClientPoll(server string, count int) (*GrpcClientPool, error) {
	pool := &GrpcClientPool{
		server: server,
		q:      list.New(),
		lock:   NewAccessLock(server, false),
	}

	for i := 0; i < count; i++ {
		conn, err := grpc.Dial(server, grpc.WithInsecure())
		if nil != err {
			Log.Error("grpc pool dial %s error:%v", server, err)

			return nil, err
		}

		pool.q.PushBack(conn)
	}

	return pool, nil
}

func (me *GrpcClientPool) Close() {
	for {
		if conn := me.Get(); nil != conn {
			conn.Close()
		} else {
			return
		}
	}
}

func (me *GrpcClientPool) Get() *grpc.ClientConn {
	var conn *grpc.ClientConn

	me.lock.Handle(func() {
		if head := me.q.Front(); nil != head {
			v := me.q.Remove(head)
			if obj, ok := v.(*grpc.ClientConn); ok {
				conn = obj
			}
		}
	})

	return conn
}

func (me *GrpcClientPool) Put(conn *grpc.ClientConn) {
	if nil != conn {
		me.lock.Handle(func() {
			me.q.PushBack(conn)
		})
	}
}

func (me *GrpcClientPool) Exec(call func(conn *grpc.ClientConn) error) error {
	var err error

	if conn := me.Get(); nil != conn {
		err = call(conn)
		me.Put(conn)
	} else {
		err = errors.New("empty grcp client pool")
	}

	return err
}
