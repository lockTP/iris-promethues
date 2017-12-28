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
	var result Result
	var labelHost, labelHandleMethod, labelStatus metrics.Label
	labelHost = metrics.Label{Name: "host", Value: ctx.HostString()}
	handle := strings.Split(ctx.GetHandlerName(), "/")
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
	err := json.Unmarshal(ctx.Response.Body(), &result)
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