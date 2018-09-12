package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.xiaojukeji.com/gulfstream/pope-offline-action-runner/models/global"
	"git.xiaojukeji.com/gulfstream/pope-offline-action-runner/models/structs"
	"git.xiaojukeji.com/gulfstream/pope-offline-action-runner/utils/http"
	"git.xiaojukeji.com/gulfstream/pope-offline-action-runner/utils/plog"
	"github.com/urfave/negroni"
)

var serverInstance *server

type server struct {
	mux *http.ServeMux
	n   *negroni.Negroni
	s   *http.Server
}

func (s *server) addRoute(path string, ctl *controller) {
	if ctl == nil || ctl.gererateRespFn == nil || ctl.getRequestErrRespFn == nil || ctl.processFn == nil {
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), global.GlobalConfig.Proxy.TimeOut.Duration)
		go func() {
			select {
			case <-ctx.Done():
				cancel()
			}
		}()

		if ctl.initCtxTraceFlag {
			ctx = httputil.GenerateTraceFromReq(ctx, r)
		}

		reqStr, err := ctl.getRequestFn(r)
		if err != nil {
			respStr := ctl.getRequestErrRespFn()
			timeStr := fmt.Sprintf("%vms", time.Since(startTime).Nanoseconds()/time.Millisecond.Nanoseconds())
			w.Write([]byte(respStr))
			plog.HttpFailed(ctx, "httprequest", "request", reqStr, "url", r.URL, "response", respStr, "proc_time", timeStr, "error", err)
			return
		}

		resp, err := ctl.processFn(ctx, reqStr)
		respStr := ctl.gererateRespFn(resp, err)
		timeStr := fmt.Sprintf("%vms", time.Since(startTime).Nanoseconds()/time.Millisecond.Nanoseconds())
		w.Write([]byte(respStr))
		plog.HttpRequestOut(ctx, "httprequest", "request", reqStr, "url", r.URL, "response", respStr, "proc_time", timeStr)
	}

	s.mux.HandleFunc(path, handler)
}

func Serve() error {
	serverInstance = &server{
		mux: http.NewServeMux(),
		n:   negroni.New(),
	}

	// api for scheduled recalc
	addRoute("/gulfstream/pope-offline/scheduled/recalc-activity", ApiReCalcActivity)
	addRoute("/gulfstream/pope-offline/scheduled/recalc-range", ApiReCalcRange)
	addRoute("/gulfstream/pope-offline/scheduled/recalc-orders", ApiReCalcOrders)
	addGetRoute("/gulfstream/pope-offline/scheduled/calc-jobs", ApiCalcJobs)
	addGetRoute("/gulfstream/pope-offline/scheduled/calc-job", ApiCalcJob)
	addGetRoute("/gulfstream/pope-offline/scheduled/calc-jobs-count", ApiCalcJobsCount)
	// api for ondemand calc
	addGetRoute("/gulfstream/pope-offline/on-demand/calc-metrics", ApiCalcMetrics)
	addGetRoute("/gulfstream/pope-offline/on-demand/calc-driver", ApiCalcDriver)
	addGetRoute("/gulfstream/pope-offline/on-demand/calc-activity", ApiCalcActivity)
	// api for driver H5
	addRoute("/gulfstream/pope-offline/my-activities/active", ApiActiveActivities)
	addRoute("/gulfstream/pope-offline/my-activities/historic", ApiHistoricActivities)
	addRoute("/gulfstream/pope-offline/my-activity", ApiActivityDetails)
	addRoute("/gulfstream/pope-offline/my-activity-orders", ApiActivityOrders)
	// api for MIS and monitoring
	addGetRoute("/gulfstream/pope-offline/reload-activities", ApiReloadActivities)
	addRoute("/gulfstream/pope-offline/activity-results", ApiActivityResults)
	addRoute("/gulfstream/pope-offline/activity-statistic", ApiActivityStatistic)
	// api for monitoring
	addGetRoute("/gulfstream/pope-offline/status/shard", ApiStatusShard)
	addGetRoute("/gulfstream/pope-offline/status/shards", ApiStatusShards)
	addGetRoute("/gulfstream/pope-offline/status/calculated-activities", ApiStatusCalculatdActivities)
	addGetRoute("/gulfstream/pope-offline/status/configured-activities", ApiStatusConfiguredActivities)
	addGetRoute("/gulfstream/pope-offline/status/activity", ApiStatusActivity)
	addGetRoute("/gulfstream/pope-offline/status/activity-database", ApiActivityDatabase)
	addHtmlRoute("/gulfstream/pope-offline/status/index.html", ApiStatusIndex)

	serverInstance.n.UseHandler(serverInstance.mux)
	serverInstance.s = &http.Server{
		Addr:         global.GlobalConfig.Proxy.HTTPPort,
		Handler:      serverInstance.n,
		ReadTimeout:  global.GlobalConfig.Proxy.HTTPServerReadTimeout.Duration,
		WriteTimeout: global.GlobalConfig.Proxy.HTTPServerWriteTimeout.Duration,
	}

	return serverInstance.s.ListenAndServe()
}

func addRoute(url string, fn func(ctx context.Context, reqData *structs.ReqData) (*structs.Response, error)) {
	ctrl := newDefaultController()
	ctrl.processFn = fn
	serverInstance.addRoute(url, ctrl)
}

func addGetRoute(url string, fn func(ctx context.Context, reqData *structs.ReqData) (*structs.Response, error)) {
	ctrl := newDefaultController()
	ctrl.getRequestFn = func(request *http.Request) (*structs.ReqData, error) {
		return &structs.ReqData{
			Queryform: request.URL.Query(),
		}, nil
	}
	ctrl.processFn = fn
	serverInstance.addRoute(url, ctrl)
}

func addHtmlRoute(url string, fn func(ctx context.Context, reqData *structs.ReqData) (*structs.Response, error)) {
	ctrl := newDefaultController()
	ctrl.getRequestFn = func(request *http.Request) (*structs.ReqData, error) {
		return &structs.ReqData{
			Queryform: request.URL.Query(),
		}, nil
	}
	ctrl.gererateRespFn = func(resp *structs.Response, err error) string {
		return resp.Data.(string)
	}
	ctrl.processFn = fn
	serverInstance.addRoute(url, ctrl)
}
