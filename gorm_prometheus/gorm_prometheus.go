package gorm_prometheus

import (
  "time"

  "github.com/penglongli/gin-metrics/ginmetrics"
  "gorm.io/gorm"
)

type MetricsPlugin struct{}

func (op *MetricsPlugin) Name() string {
  return "metricsPlugin"
}

func (op *MetricsPlugin) Initialize(db *gorm.DB) (err error) {
  // 开始前
  _ = db.Callback().Create().Before("gorm:before_create").Register("before_create", before)
  _ = db.Callback().Query().Before("gorm:query").Register("before_query", before)
  _ = db.Callback().Delete().Before("gorm:before_delete").Register("before_delete", before)
  _ = db.Callback().Update().Before("gorm:setup_reflect_value").Register("before_update", before)
  _ = db.Callback().Row().Before("gorm:row").Register("before_row", before)
  _ = db.Callback().Raw().Before("gorm:raw").Register("before_raw", before)

  // 结束后
  _ = db.Callback().Create().After("gorm:after_create").Register("after_create", after)
  _ = db.Callback().Query().After("gorm:after_query").Register("after_query", after)
  _ = db.Callback().Delete().After("gorm:after_delete").Register("after_delete", after)
  _ = db.Callback().Update().After("gorm:after_update").Register("after_update", after)
  _ = db.Callback().Row().After("gorm:row").Register("row", after)
  _ = db.Callback().Raw().After("gorm:raw").Register("raw", after)
  return
}

// 回调前
func before(db *gorm.DB) {
  db.InstanceSet("startTime", time.Now())
  return
}

// 回调后
func after(db *gorm.DB) {
  _ts, isExist := db.InstanceGet("startTime")
  if !isExist {
    return
  }

  ts, ok := _ts.(time.Time)
  if !ok {
    return
  }

  sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
  elapsed := time.Since(ts)
  var action string
  switch {
  case strings.HasPrefix(sql, "SELECT"):
    action = "select"
  case strings.HasPrefix(sql, "UPDATE"):
    action = "update"
  case strings.HasPrefix(sql, "DELETE"):
    action = "delete"
  case strings.HasPrefix(sql, "INSERT"):
    action = "insert"
  }
  constTime := float64(elapsed.Nanoseconds() / 1e6)
  addMysqlPrometheusLabelValue(action, constTime)
  return
}

// HistogramMetric 声明prometheus 指标
var (
  HistogramMetric = &ginmetrics.Metric{
    Type:        ginmetrics.Histogram,
    Name:        "mysql_operate_duration_milliseconds",
    Description: "an example of gauge type metric",
    Buckets:     []float64{0.1, 0.5, 1, 2, 3, 5, 10, 20, 50, 100},
    Labels:      []string{"action"},
  }
)

// RegisterMysqlPrometheus 注册 mysql prometheus 指标
func RegisterMysqlPrometheus() error {
  return ginmetrics.GetMonitor().AddMetric(HistogramMetric)
}

func addMysqlPrometheusLabelValue(action string, constTime float64) error {
  var label []string
  label = append(label, action)
  return ginmetrics.GetMonitor().GetMetric("mysql_operate_duration_milliseconds").Observe(label, constTime)
}
