// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package impalaservice

import (
	"bytes"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/koblas/impalathing/services/beeswax"
	"github.com/koblas/impalathing/services/cli_service"
	"github.com/koblas/impalathing/services/status"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var _ = status.GoUnusedProtection__
var _ = beeswax.GoUnusedProtection__
var _ = cli_service.GoUnusedProtection__

type ImpalaHiveServer2Service interface {
	cli_service.TCLIService

	ResetCatalog() (r *status.TStatus, err error)
}

type ImpalaHiveServer2ServiceClient struct {
	*cli_service.TCLIServiceClient
}

func NewImpalaHiveServer2ServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *ImpalaHiveServer2ServiceClient {
	return &ImpalaHiveServer2ServiceClient{TCLIServiceClient: cli_service.NewTCLIServiceClientFactory(t, f)}
}

func NewImpalaHiveServer2ServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *ImpalaHiveServer2ServiceClient {
	return &ImpalaHiveServer2ServiceClient{TCLIServiceClient: cli_service.NewTCLIServiceClientProtocol(t, iprot, oprot)}
}

func (p *ImpalaHiveServer2ServiceClient) ResetCatalog() (r *status.TStatus, err error) {
	if err = p.sendResetCatalog(); err != nil {
		return
	}
	return p.recvResetCatalog()
}

func (p *ImpalaHiveServer2ServiceClient) sendResetCatalog() (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	if err = oprot.WriteMessageBegin("ResetCatalog", thrift.CALL, p.SeqId); err != nil {
		return
	}
	args := ImpalaHiveServer2ServiceResetCatalogArgs{}
	if err = args.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	return oprot.Flush()
}

func (p *ImpalaHiveServer2ServiceClient) recvResetCatalog() (value *status.TStatus, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	method, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if method != "ResetCatalog" {
		err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "ResetCatalog failed: wrong method name")
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ResetCatalog failed: out of sequence response")
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error70 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error71 error
		error71, err = error70.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error71
		return
	}
	if mTypeId != thrift.REPLY {
		err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "ResetCatalog failed: invalid message type")
		return
	}
	result := ImpalaHiveServer2ServiceResetCatalogResult{}
	if err = result.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	value = result.GetSuccess()
	return
}

type ImpalaHiveServer2ServiceProcessor struct {
	*cli_service.TCLIServiceProcessor
}

func NewImpalaHiveServer2ServiceProcessor(handler ImpalaHiveServer2Service) *ImpalaHiveServer2ServiceProcessor {
	self72 := &ImpalaHiveServer2ServiceProcessor{cli_service.NewTCLIServiceProcessor(handler)}
	self72.AddToProcessorMap("ResetCatalog", &impalaHiveServer2ServiceProcessorResetCatalog{handler: handler})
	return self72
}

type impalaHiveServer2ServiceProcessorResetCatalog struct {
	handler ImpalaHiveServer2Service
}

func (p *impalaHiveServer2ServiceProcessorResetCatalog) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := ImpalaHiveServer2ServiceResetCatalogArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("ResetCatalog", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()
	result := ImpalaHiveServer2ServiceResetCatalogResult{}
	var retval *status.TStatus
	var err2 error
	if retval, err2 = p.handler.ResetCatalog(); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing ResetCatalog: "+err2.Error())
		oprot.WriteMessageBegin("ResetCatalog", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return true, err2
	} else {
		result.Success = retval
	}
	if err2 = oprot.WriteMessageBegin("ResetCatalog", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

type ImpalaHiveServer2ServiceResetCatalogArgs struct {
}

func NewImpalaHiveServer2ServiceResetCatalogArgs() *ImpalaHiveServer2ServiceResetCatalogArgs {
	return &ImpalaHiveServer2ServiceResetCatalogArgs{}
}

func (p *ImpalaHiveServer2ServiceResetCatalogArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *ImpalaHiveServer2ServiceResetCatalogArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ResetCatalog_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ImpalaHiveServer2ServiceResetCatalogArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ImpalaHiveServer2ServiceResetCatalogArgs(%+v)", *p)
}

// Attributes:
//  - Success
type ImpalaHiveServer2ServiceResetCatalogResult struct {
	Success *status.TStatus `thrift:"success,0" json:"success,omitempty"`
}

func NewImpalaHiveServer2ServiceResetCatalogResult() *ImpalaHiveServer2ServiceResetCatalogResult {
	return &ImpalaHiveServer2ServiceResetCatalogResult{}
}

var ImpalaHiveServer2ServiceResetCatalogResult_Success_DEFAULT *status.TStatus

func (p *ImpalaHiveServer2ServiceResetCatalogResult) GetSuccess() *status.TStatus {
	if !p.IsSetSuccess() {
		return ImpalaHiveServer2ServiceResetCatalogResult_Success_DEFAULT
	}
	return p.Success
}
func (p *ImpalaHiveServer2ServiceResetCatalogResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *ImpalaHiveServer2ServiceResetCatalogResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *ImpalaHiveServer2ServiceResetCatalogResult) readField0(iprot thrift.TProtocol) error {
	p.Success = &status.TStatus{}
	if err := p.Success.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Success), err)
	}
	return nil
}

func (p *ImpalaHiveServer2ServiceResetCatalogResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ResetCatalog_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ImpalaHiveServer2ServiceResetCatalogResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Success), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *ImpalaHiveServer2ServiceResetCatalogResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ImpalaHiveServer2ServiceResetCatalogResult(%+v)", *p)
}
