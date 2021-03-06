package syshttp

import (
	"errors"

	belogs "github.com/astaxie/beego/logs"
	"github.com/cpusoft/go-json-rest/rest"
	httpserver "github.com/cpusoft/goutil/httpserver"
	jsonutil "github.com/cpusoft/goutil/jsonutil"

	"model"
	sysmodel "sys/model"
	"sys/sys"
)

//
func InitReset(w rest.ResponseWriter, req *rest.Request) {
	belogs.Debug("InitReset()")
	sysStyle := sysmodel.SysStyle{}
	//shaodebug
	//content, _ := ioutil.ReadAll(req.Body)
	//belogs.Debug("InitReset():ReadAll body:", string(content))

	err := req.DecodeJsonPayload(&sysStyle)
	if err != nil {
		belogs.Error("InitReset(): DecodeJsonPayload:", err)
		w.WriteJson(httpserver.GetFailHttpResponse(err))
		return
	}
	belogs.Info("InitReset():get sysStyle:", jsonutil.MarshalJson(sysStyle))
	if sysStyle.SysStyle != "init" && sysStyle.SysStyle != "fullsync" && sysStyle.SysStyle != "resetall" {
		belogs.Error("InitReset(): SysStyle should be init or fullsync or resetall, it is ", sysStyle.SysStyle)
		w.WriteJson(httpserver.GetFailHttpResponse(errors.New("SysStyle should be init or fullsync or resetall")))
		return
	}
	belogs.Debug("InitReset(): sysStyle:", sysStyle)

	err = sys.InitReset(sysStyle)
	if err != nil {
		w.WriteJson(httpserver.GetFailHttpResponse(err))
	} else {
		w.WriteJson(httpserver.GetOkHttpResponse())
	}
}

// detail
func DetailStates(w rest.ResponseWriter, req *rest.Request) {
	belogs.Info("DetailStates()")

	detailStates, err := sys.DetailStates()
	if err != nil {
		w.WriteJson(httpserver.GetFailHttpResponse(err))
		return
	}
	belogs.Info("DetailStates():detailStates:", jsonutil.MarshalJson(detailStates))

	stateResponse := model.StateResponse{
		HttpResponse: httpserver.GetOkHttpResponse(),
		State:        detailStates,
	}
	w.WriteJson(stateResponse)
}

// summary
func SummaryStates(w rest.ResponseWriter, req *rest.Request) {
	belogs.Info("SummaryStates()")

	summaryStates, err := sys.SummaryStates()
	if err != nil {
		w.WriteJson(httpserver.GetFailHttpResponse(err))
		return
	}
	belogs.Info("SummaryStates():summaryStates:", jsonutil.MarshalJson(summaryStates))

	stateResponse := model.StateResponse{
		HttpResponse: httpserver.GetOkHttpResponse(),
		State:        summaryStates,
	}
	w.WriteJson(stateResponse)
}

// just return valid/warning/invalid count in cer/roa/mft/crl
func Results(w rest.ResponseWriter, req *rest.Request) {
	belogs.Info("Results()")

	results, err := sys.Results()
	if err != nil {
		w.WriteJson(httpserver.GetFailHttpResponse(err))
		return
	}
	belogs.Info("Results():results:", jsonutil.MarshalJson(results))
	w.WriteJson(results)
}

func ExportRoas(w rest.ResponseWriter, req *rest.Request) {
	belogs.Info("ExportRoas()")
	results, err := sys.ExportRoas()
	if err != nil {
		w.WriteJson(httpserver.GetFailHttpResponse(err))
		return
	}
	belogs.Info("ExportRoas():results:", jsonutil.MarshalJson(results))
	w.WriteJson(results)
}
