package oauth

import "net/http"

type loggableResponse struct {
	Status     string
	StatusCode int
	Proto      string
	ProtoMajor int
	ProtoMinor int

	Header http.Header

	ContentLength    int64
	TransferEncoding []string
	Close            bool

	Trailer http.Header

	Request loggableRequest
}

func newLoggableResponse(resp *http.Response) loggableResponse {
	loggableResp := loggableResponse{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMajor:       resp.ProtoMajor,
		ProtoMinor:       resp.ProtoMinor,
		Header:           resp.Header,
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Close:            resp.Close,
		Trailer:          resp.Trailer,
	}

	if resp.Request != nil {
		loggableResp.Request = newLoggableRequest(*resp.Request)
	}

	return loggableResp
}
