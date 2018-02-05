package asdf


type MsgChanResponse struct {
	Err  error
	Data interface{}
}

type MsgChanRequest struct {
	Chan chan *MsgChanResponse
	Cmd  int
	Data interface{}
}

func RequestByChan(ch chan *MsgChanRequest, cmd int, obj interface{}, buf []byte) (error, interface{}) {
	chPrivate := make(chan *MsgChanResponse)

	if nil != buf {
		err := json.Unmarshal(buf, obj)
		if nil != err {
			return err, nil
		}
	}

	ch <- &MsgChanRequest{
		Chan: chPrivate,
		Cmd:  cmd,
		Data: obj,
	}

	response := <-chPrivate

	close(chPrivate)

	return response.Err, response.Data
}
