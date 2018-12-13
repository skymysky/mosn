/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

import (
	"fmt"
	"github.com/alipay/sofa-mosn/pkg/protocol"
	"github.com/alipay/sofa-mosn/pkg/types"
	"net"
	"runtime"
	"testing"
	"time"
)

func BenchmarkAccessLog(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	InitDefaultLogger("", INFO)
	// ~ replace the path if needed
	accessLog, err := NewAccessLog("/tmp/mosn_bench/benchmark_access.log", nil, "")

	if err != nil {
		fmt.Errorf(err.Error())
	}
	reqHeaders := map[string]string{
		"service": "test",
	}

	respHeaders := map[string]string{
		"Server": "MOSN",
	}

	requestInfo := newRequestInfo()
	requestInfo.SetRequestReceivedDuration(time.Now())
	requestInfo.SetResponseReceivedDuration(time.Now().Add(time.Second * 2))
	requestInfo.SetBytesSent(2048)
	requestInfo.SetBytesReceived(2048)

	requestInfo.SetResponseFlag(0)
	requestInfo.SetUpstreamLocalAddress(&net.TCPAddr{[]byte("127.0.0.1"), 23456, ""})
	requestInfo.SetDownstreamLocalAddress(&net.TCPAddr{[]byte("127.0.0.1"), 12200, ""})
	requestInfo.SetDownstreamRemoteAddress(&net.TCPAddr{[]byte("127.0.0.2"), 53242, ""})
	requestInfo.OnUpstreamHostSelected(nil)

	for n := 0; n < b.N; n++ {
		accessLog.Log(protocol.CommonHeader(reqHeaders), protocol.CommonHeader(respHeaders), requestInfo)
	}
}

func BenchmarkAccessLogParallel(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	InitDefaultLogger("", INFO)
	// ~ replace the path if needed
	accessLog, err := NewAccessLog("/tmp/mosn_bench/benchmark_access.log", nil, "")

	if err != nil {
		fmt.Errorf(err.Error())
	}
	reqHeaders := map[string]string{
		"service": "test",
	}

	respHeaders := map[string]string{
		"Server": "MOSN",
	}

	requestInfo := newRequestInfo()
	requestInfo.SetRequestReceivedDuration(time.Now())
	requestInfo.SetResponseReceivedDuration(time.Now().Add(time.Second * 2))
	requestInfo.SetBytesSent(2048)
	requestInfo.SetBytesReceived(2048)

	requestInfo.SetResponseFlag(0)
	requestInfo.SetUpstreamLocalAddress(&net.TCPAddr{[]byte("127.0.0.1"), 23456, ""})
	requestInfo.SetDownstreamLocalAddress(&net.TCPAddr{[]byte("127.0.0.1"), 12200, ""})
	requestInfo.SetDownstreamRemoteAddress(&net.TCPAddr{[]byte("127.0.0.2"), 53242, ""})
	requestInfo.OnUpstreamHostSelected(nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			accessLog.Log(protocol.CommonHeader(reqHeaders), protocol.CommonHeader(respHeaders), requestInfo)
		}
	})
}

// mock_requestInfo
type mock_requestInfo struct {
	protocol                 types.Protocol
	startTime                time.Time
	responseFlag             types.ResponseFlag
	upstreamHost             types.HostInfo
	requestReceivedDuration  time.Duration
	responseReceivedDuration time.Duration
	bytesSent                uint64
	bytesReceived            uint64
	responseCode             uint32
	localAddress             net.Addr
	downstreamLocalAddress   net.Addr
	downstreamRemoteAddress  net.Addr
	isHealthCheckRequest     bool
	routerRule               types.RouteRule
}

// NewrequestInfo
func newRequestInfo() types.RequestInfo {
	return &mock_requestInfo{
		startTime: time.Now(),
	}
}

func (r *mock_requestInfo) StartTime() time.Time {
	return r.startTime
}

func (r *mock_requestInfo) SetStartTime() {
	r.startTime = time.Now()
}

func (r *mock_requestInfo) RequestReceivedDuration() time.Duration {
	return r.requestReceivedDuration
}

func (r *mock_requestInfo) SetRequestReceivedDuration(time time.Time) {
	r.requestReceivedDuration = time.Sub(r.startTime)
}

func (r *mock_requestInfo) ResponseReceivedDuration() time.Duration {
	return r.responseReceivedDuration
}

func (r *mock_requestInfo) SetResponseReceivedDuration(time time.Time) {
	r.responseReceivedDuration = time.Sub(r.startTime)
}

func (r *mock_requestInfo) BytesSent() uint64 {
	return r.bytesSent
}

func (r *mock_requestInfo) SetBytesSent(bytesSent uint64) {
	r.bytesSent = bytesSent
}

func (r *mock_requestInfo) BytesReceived() uint64 {
	return r.bytesReceived
}

func (r *mock_requestInfo) SetBytesReceived(bytesReceived uint64) {
	r.bytesReceived = bytesReceived
}

func (r *mock_requestInfo) Protocol() types.Protocol {
	return r.protocol
}

func (r *mock_requestInfo) ResponseCode() uint32 {
	return r.responseCode
}

func (r *mock_requestInfo) SetResponseCode(code uint32) {
	r.responseCode = code
}

func (r *mock_requestInfo) Duration() time.Duration {
	return time.Now().Sub(r.startTime)
}

func (r *mock_requestInfo) GetResponseFlag(flag types.ResponseFlag) bool {
	return r.responseFlag&flag != 0
}

func (r *mock_requestInfo) SetResponseFlag(flag types.ResponseFlag) {
	r.responseFlag |= flag
}

func (r *mock_requestInfo) UpstreamHost() types.HostInfo {
	return r.upstreamHost
}

func (r *mock_requestInfo) OnUpstreamHostSelected(host types.HostInfo) {
	r.upstreamHost = host
}

func (r *mock_requestInfo) UpstreamLocalAddress() net.Addr {
	return r.localAddress
}

func (r *mock_requestInfo) SetUpstreamLocalAddress(addr net.Addr) {
	r.localAddress = addr
}

func (r *mock_requestInfo) IsHealthCheck() bool {
	return r.isHealthCheckRequest
}

func (r *mock_requestInfo) SetHealthCheck(isHc bool) {
	r.isHealthCheckRequest = isHc
}

func (r *mock_requestInfo) DownstreamLocalAddress() net.Addr {
	return r.downstreamLocalAddress
}

func (r *mock_requestInfo) SetDownstreamLocalAddress(addr net.Addr) {
	r.downstreamLocalAddress = addr
}

func (r *mock_requestInfo) DownstreamRemoteAddress() net.Addr {
	return r.downstreamRemoteAddress
}

func (r *mock_requestInfo) SetDownstreamRemoteAddress(addr net.Addr) {
	r.downstreamRemoteAddress = addr
}

func (r *mock_requestInfo) RouteEntry() types.RouteRule {
	return r.routerRule
}

func (r *mock_requestInfo) SetRouteEntry(routerRule types.RouteRule) {
	r.routerRule = routerRule
}
