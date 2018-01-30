package prometheusMetrics

import (
	prometheusSink "github.com/lockTP/go-metrics/prometheus"
	metrics "github.com/lockTP/go-metrics"
	"gopkg.in/kataras/iris.v5"
	"encoding/json"
	"strings"
	"strconv"
	"sort"
)

var (
	Met     *metrics.Metrics
)


func NewPrometheusMetrics(serviceName string) {
	sink,_ := prometheusSink.NewPrometheusSink()
	conf := metrics.DefaultConfig(serviceName)
	Met, _ = metrics.New(conf, sink)
}

/**
* api拦截记录
 */
func Record(ctx *iris.Context) {
	hostStr := ctx.HostString()
	handlerName := ctx.GetHandlerName()
	body := ctx.Response.Body()
	//异步记录日志
	go ApiRecord(hostStr, handlerName, body)
}

func ApiRecord(hostStr string, handlerName string, body []byte) {
	var result Result
	var labelHost, labelHandleMethod, labelStatus metrics.Label
	labelHost = metrics.Label{Name: "host", Value: hostStr}
	handle := strings.Split(handlerName, "/")
	var _api []string
	var _f []string
	var _c string = ""
	var _system string = ""
	var _controller string = ""
	var _func string = ""
	var _interface string = ""
	if len(handle) >= 2 {
		_api = strings.Split(handle[len(handle)-1], ".")
		if len(_api) >= 2 {
			_f = strings.Split(_api[len(_api)-1], "-")
			_c = _api[len(_api)-2]
		}
		if len(_c) >= 3 {
			_controller = _c[2 : len(_c)-1]
		}
		_system = strings.Replace(handle[1], "-", "_", -1)
		if len(_f) >= 1 {
			_func = _f[0]
		}
	}
	//接口名称
	_interface = _system + "_" + _controller + "_" + _func
	labelHandleMethod = metrics.Label{Name: "handleMethod", Value: _interface}
	err := json.Unmarshal(body, &result)
	if err != nil {//返回非合法json格式
		labelStatus = metrics.Label{Name: "status", Value: "9999"}
	} else {
		labelStatus = metrics.Label{Name: "status", Value: strconv.Itoa(result.Status)}
	}
	labels := []metrics.Label{}
	labels = append(labels, labelHost, labelHandleMethod, labelStatus)
	Met.AddSampleWithLabels([]string{"api"}, 1, labels)
}

/**
* api拦截记录，只适用于较高版本的iris
 */
func Record_New(ctx iris_context.Context) {
	var result Result
	var labelHost, labelHandleMethod, labelStatus metrics.Label
	labelHost = metrics.Label{Name: "host", Value: ctx.Host()}
	handle := strings.Split(ctx.HandlerName(), "/")
	var _api []string
	var _f []string
	var _c string = ""
	var _system string = ""
	var _controller string = ""
	var _func string = ""
	var _interface string = ""
	if len(handle) >= 2 {
		_api = strings.Split(handle[len(handle)-1], ".")
		if len(_api) >= 2 {
			_f = strings.Split(_api[len(_api)-1], "-")
			_c = _api[len(_api)-2]
		}
		if len(_c) >= 3 {
			_controller = _c[2 : len(_c)-1]
		}
		_system = strings.Replace(handle[1], "-", "_", -1)
		if len(_f) >= 1 {
			_func = _f[0]
		}
	}
	//接口名称
	_interface = _system + "_" + _controller + "_" + _func
	labelHandleMethod = metrics.Label{Name: "handleMethod", Value: _interface}
	recorder, flag := ctx.IsRecording()
	if !flag {//未记录response body
		labelStatus = metrics.Label{Name: "status", Value: "9998"}
	} else {
		err := json.Unmarshal(recorder.Body(), &result)
		if err != nil {//返回非合法json格式
			labelStatus = metrics.Label{Name: "status", Value: "9999"}
		} else {
			labelStatus = metrics.Label{Name: "status", Value: strconv.Itoa(result.Status)}
		}
	}
	labels := []metrics.Label{}
	labels = append(labels, labelHost, labelHandleMethod, labelStatus)
	Met.AddSampleWithLabels([]string{"api"}, 1, labels)
}


/**
* 	单一label记录（极简版）
*   记录至  <namespace>_<name>_simple_count{lable="<你的输入>"}  这个metrics中
 */
func SimpleRecord(input string) {
	var label metrics.Label
	label = metrics.Label{Name: "lable", Value: input}
	labels := []metrics.Label{}
	labels = append(labels, label)
	Met.AddSampleWithLabels([]string{"simple"}, 1, labels)
}

/**
*   自定义记录
*   记录至  <namespace>_<name>_<metricsName>_count{labelMap<key>="labelMap<value>"...}  这个metrics中
 */
func CustomRecord(metricsName string, lableMap map[string]string) {
	//将乱序的map按顺序输出
	var keys []string
	for k := range lableMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	labels := []metrics.Label{}
	for _, k := range keys {
		label := metrics.Label{Name: k, Value: lableMap[k]}
		labels = append(labels, label)
	}
	Met.AddSampleWithLabels([]string{metricsName}, 1, labels)
}

/**
* 	数据库链接数监控
 */
func DBConnectCount (count int, dbType string ) {
	label := metrics.Label{Name: "dbType", Value: dbType}
	labels := []metrics.Label{}
	labels = append(labels, label)
	Met.SetGaugeWithLabels([]string{"DBOpenCount"}, float32(count), labels)
}
