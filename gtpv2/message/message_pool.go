package message

type messagePool struct {
	modifyBearerRequests  chan *ModifyBearerRequest
	modifyBearerResponses chan *ModifyBearerResponse
}

var msgPool = newMessagePool(100)

func newMessagePool(bufferLen int) *messagePool {
	return &messagePool{
		modifyBearerRequests:  make(chan *ModifyBearerRequest, bufferLen),
		modifyBearerResponses: make(chan *ModifyBearerResponse, bufferLen),
	}
}

func SetMessagePoolSize(bufferLen int) {
	msgPool = newMessagePool(bufferLen)
}

func (mp *messagePool) getModifyBearerRequest() (m *ModifyBearerRequest) {
	select {
	case m = <-mp.modifyBearerRequests:
	default:
		m = &ModifyBearerRequest{}
	}
	return m
}

func (mp *messagePool) releaseModifyBearerRequest(m *ModifyBearerRequest) {
	if m == nil {
		return
	}
	select {
	case mp.modifyBearerRequests <- m:
	default:
	}
}

func (mp *messagePool) getModifyBearerResponse() (m *ModifyBearerResponse) {
	select {
	case m = <-mp.modifyBearerResponses:
	default:
		m = &ModifyBearerResponse{}
	}
	return m
}

func (mp *messagePool) releaseModifyBearerResponse(m *ModifyBearerResponse) {
	if m == nil {
		return
	}
	select {
	case mp.modifyBearerResponses <- m:
	default:
	}
}
