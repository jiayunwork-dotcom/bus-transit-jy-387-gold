# bus-transit

城市公交线路运营效率评估工具。

读取站点时刻 CSV（含计划到达、实际到达、载客量、额定容量），按线路计算
**准点率** 与 **平均满载率**，并对低于准点率阈值或超员的线路标记低效能。

- 退出码 `0`：全部线路达标
- 退出码 `1`：存在低效能线路（或运行时错误）
- 退出码 `2`：参数缺失（用法错误）

## 用法

```
bus-transit -stops <csv> [-format text|json]
```

## CSV 格式

stops.csv:
```
route,trip_id,stop,sched_arr,act_arr,passengers,capacity
R1,T1,S1,08:00,08:00,20,50
```

时间列为 `HH:MM`，准点判定阈值默认 3 分钟。

## 实现说明

纯标准库实现，离线可构建。
