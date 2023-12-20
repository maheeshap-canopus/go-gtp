package message

type messagePool struct {
	createSessionRequests         chan *CreateSessionRequest
	createSessionResponses        chan *CreateSessionResponse
	createBearerRequests          chan *CreateBearerRequest
	createBearerResponses         chan *CreateBearerResponse
	modifyBearerRequests          chan *ModifyBearerRequest
	modifyBearerResponses         chan *ModifyBearerResponse
	releaseAccessBearersRequests  chan *ReleaseAccessBearersRequest
	releaseAccessBearersResponses chan *ReleaseAccessBearersResponse
	deleteSessionRequests         chan *DeleteSessionRequest
	deleteSessionResponses        chan *DeleteSessionResponse
}

var msgPool = newMessagePool(100)

func newMessagePool(bufferLen int) *messagePool {
	return &messagePool{
		createSessionRequests:         make(chan *CreateSessionRequest, bufferLen),
		createSessionResponses:        make(chan *CreateSessionResponse, bufferLen),
		createBearerRequests:          make(chan *CreateBearerRequest, bufferLen),
		createBearerResponses:         make(chan *CreateBearerResponse, bufferLen),
		modifyBearerRequests:          make(chan *ModifyBearerRequest, bufferLen),
		modifyBearerResponses:         make(chan *ModifyBearerResponse, bufferLen),
		releaseAccessBearersRequests:  make(chan *ReleaseAccessBearersRequest, bufferLen),
		releaseAccessBearersResponses: make(chan *ReleaseAccessBearersResponse, bufferLen),
		deleteSessionRequests:         make(chan *DeleteSessionRequest, bufferLen),
		deleteSessionResponses:        make(chan *DeleteSessionResponse, bufferLen),
	}
}

func SetMessagePoolSize(bufferLen int) {
	msgPool = newMessagePool(bufferLen)
}

func (mp *messagePool) getCreateSessionRequest() (m *CreateSessionRequest) {
	select {
	case m = <-mp.createSessionRequests:
	default:
		m = &CreateSessionRequest{}
	}
	return m
}

func (mp *messagePool) releaseCreateSessionRequest(m *CreateSessionRequest) {
	if m == nil {
		return
	}
	select {
	case mp.createSessionRequests <- m:
	default:
	}
}

func (mp *messagePool) getCreateSessionResponse() (m *CreateSessionResponse) {
	select {
	case m = <-mp.createSessionResponses:
	default:
		m = &CreateSessionResponse{}
	}
	return m
}

func (mp *messagePool) releaseCreateSessionResponse(m *CreateSessionResponse) {
	if m == nil {
		return
	}
	select {
	case mp.createSessionResponses <- m:
	default:
	}
}

func (mp *messagePool) getCreateBearerRequest() (m *CreateBearerRequest) {
	select {
	case m = <-mp.createBearerRequests:
	default:
		m = &CreateBearerRequest{}
	}
	return m
}

func (mp *messagePool) releaseCreateBearerRequest(m *CreateBearerRequest) {
	if m == nil {
		return
	}
	select {
	case mp.createBearerRequests <- m:
	default:
	}
}

func (mp *messagePool) getCreateBearerResponse() (m *CreateBearerResponse) {
	select {
	case m = <-mp.createBearerResponses:
	default:
		m = &CreateBearerResponse{}
	}
	return m
}

func (mp *messagePool) releaseCreateBearerResponse(m *CreateBearerResponse) {
	if m == nil {
		return
	}
	select {
	case mp.createBearerResponses <- m:
	default:
	}
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

func (mp *messagePool) getReleaseAccessBearersRequest() (m *ReleaseAccessBearersRequest) {
	select {
	case m = <-mp.releaseAccessBearersRequests:
	default:
		m = &ReleaseAccessBearersRequest{}
	}
	return m
}

func (mp *messagePool) releaseReleaseAccessBearersRequest(m *ReleaseAccessBearersRequest) {
	if m == nil {
		return
	}
	select {
	case mp.releaseAccessBearersRequests <- m:
	default:
	}
}

func (mp *messagePool) getReleaseAccessBearersResponse() (m *ReleaseAccessBearersResponse) {
	select {
	case m = <-mp.releaseAccessBearersResponses:
	default:
		m = &ReleaseAccessBearersResponse{}
	}
	return m
}

func (mp *messagePool) releaseReleaseAccessBearersResponse(m *ReleaseAccessBearersResponse) {
	if m == nil {
		return
	}
	select {
	case mp.releaseAccessBearersResponses <- m:
	default:
	}
}

func (mp *messagePool) getDeleteSessionRequest() (m *DeleteSessionRequest) {
	select {
	case m = <-mp.deleteSessionRequests:
	default:
		m = &DeleteSessionRequest{}
	}
	return m
}

func (mp *messagePool) releaseDeleteSessionRequest(m *DeleteSessionRequest) {
	if m == nil {
		return
	}
	select {
	case mp.deleteSessionRequests <- m:
	default:
	}
}

func (mp *messagePool) getDeleteSessionResponse() (m *DeleteSessionResponse) {
	select {
	case m = <-mp.deleteSessionResponses:
	default:
		m = &DeleteSessionResponse{}
	}
	return m
}

func (mp *messagePool) releaseDeleteSessionResponse(m *DeleteSessionResponse) {
	if m == nil {
		return
	}
	select {
	case mp.deleteSessionResponses <- m:
	default:
	}
}
