# Dacrd Backend 2024

## HW

[apply link](https://boards.greenhouse.io/dcard/jobs/5650761)

[HW Link](https://drive.google.com/file/d/1dnDiBDen7FrzOAJdKZMDJg479IC77_zT/view)


:::info

有效率的解決
對 在特定區間有效的時序性資料
做 多條件查詢

( 100 Request Per Second )


Age 會是常態分佈 
( 以 Dcard 來說可能是 15-40 歲的 User 最多 )

:::

### Condition 
- Age <1~100>
- Gender <enum:M、F>
- Country <enum:TW、JP 等符合 
- Platform <enum:android, ios, web>

> Country : https://zh.wikipedia.org/wiki/ISO_3166-1 

### Amount condition 

- 10000 Requests Per Second 
- 同時存在系統的總活躍廣告數量 (也就是 StartAt < NOW < EndAt) < 1000
- 每天 create 的廣告數量 不會超過 3000 個

### Example : Create 

```
curl -X POST -H "Content-Type: application/json" \
"http://<host>/api/v1/ad" \
--data '{
    "title" "AD 55",
    "startAt" "2023-12-10T03:00:00.000Z",
    "endAt" "2023-12-31T16:00:00.000Z",
    "conditions": {
        {
            "ageStart": 20,
            "ageEnd": 30,
            "country: ["TW", "JP"],
            "platform": ["android", "ios"]
        }
    }
}'
```

### Example : Read 
```
curl -X GET -H "Content-Type: application/json" \
"http://<host>/api/v1/ad?offset=10&limit=3&age=24&gender=F&country=TW&platform=ios"

# response

{
    "items": [
        {
            "title": "AD 1",
            "endAt" "2023-12-22T01:00:00.000Z"
        },
        {
            "title": "AD 31",
            "endAt" "2023-12-30T12:00:00.000Z"
        },
        {
            "title": "AD 10",
            "endAt" "2023-12-31T16:00:00.000Z"
        }
    ]
}
```

## Question ? 

```

```

## Preload Ad Cache

:::info

因為 當前活躍的廣告數 不會超過 1000
所以可以用 CornJob ( 可能每 30 mins 看一次 )
檢查有沒有快要需要投放的廣告需要 pre-load 到 Redis
> 並且 pre-load 到 Redis ( Create 的部分 ) 會以 Lua Script 撰寫
> 保證 ACID Transaction 


如果太早 pre-laod 會很容易塞滿 Redis
> 因為一天至多會 Create 到 3000 個廣告 
> 但廣告區間不一定會在當天

:::

:::warning

Pre-load 到 Redis 的 Corn Job 如果 fail 的處理：

使用 k8s 的 job 

:::
    
## Solution 1
    
:::info
Structure
```
`Gender:Platform:Country:Age:uuid` : "title,end",
`Gender:Platform:Country:Age:uuid` : "title,end",
...
```
當建立廣告時，也會同步在 Redis 
開一個 `Gender:Platform:Country:Age:uuid` 的 String 
然後以 `SCAN` 來去得 Query 
> 在 Query 時會以

再 Query 前會先檢查是否有 `cache:Gender:Platform:Country:Age` 
> 之前有以相同條件下去 Search 的 Cache

會再把當前的 Query 結果以 `ZSET` 加到 Redis
key 設為`cache:Gender:Platform:Country:Age` 
並且會設定相對短，但有不會到太短的 TTL ( 5-10mins )
> 方便下次用相同查詢條件不同 offset,limit 時去做查詢
:::

:::warning
**Pros**
- 自動 TTL ，不需要額外處理
- 每個廣告都以獨立的 Key 存，如果有修改廣告很好改

**Cons**

- 可能會有太多Search 的 Cache ( `cache:Gender:Platform:Country:Age`  )
佔滿 Redis 
> 假設以 10 Country Enum 為例
> 至少會有 2 * 3 * 100 * 10 ( Gender * Platform * Age * Country )
> 6000 個 key 
> ( 雖然不算很多 Key，但是每個都是 ZSET，也蠻佔 Memory 的 )
- Search 的時候還需要多很多處理
> 因為 Age 的 Field 有點難處理
> - 分 4 次 SCAN 再取交集
> - 前後多 SCAN 一些 再對前後 trim
- Search 時主要靠 SCAN ( loop through , O(N) ) 
> 這邊的 N 是 total redis key 的數量
- Pagenation 不能在 Redis 層做到
:::

key field order : big > small 

Age : `ssbb`

Key : `Gender:Platform:Country:Age`

Key : `Gender:Platform:Country:Age:uuid`

match by regx 

redis search key by pattern time complexity : O(N)

zset :  2^32 - 1 -> 32 bit
    

lua script ( store procedure )

[estimate redis memory usage](https://lucasmagnum.medium.com/redistip-estimate-the-memory-usage-for-repeated-keys-in-redis-2dc3f163fdab)

[SCAN + offset + limit ](https://jinguoxing.github.io/redis/2018/09/04/redis-scan/)


regex example : 
    
```
[0-9][0-9]-[0-9][0-9]*
```

```
[a-c][b-9]-[a-c][b-d]*
```

Age Set: `ab` - `cd`
> 1-100 -> 00-99

Age Match: `ef`

> start <= ef , ef <= end

start : `[0-(e-1)][1-9]` or `e[0-f]`
end : `[(e+1)-9][f-9]` or `e[f-9]`

```
( [0-(e-1)][1-9] or e[0-f] )-[(e+1)-9][f-9]
```

> 在沒有 cache 到的時候要開 4 個 goroutine 去 抓結果
> ( 或是寫 lua script )
> ( 因為 start end 各有 2 個，總共 4 個可能 )
> 再把結果存到 `cache:Country:FM:Platform:Age` ( ZSET )
> 在有 Cache 過的狀態只需要 O( LogN + M ) 的複雜度就可以拿到結果



```
[0-e][0-9]-[e-9][0-9]
```
> 可以選到最近的區間，再對前後做 trim 
> ( 可以在 Lua script 中做 )

## Solution 2

:::info

Structure : 
```
0 : {
        `Gender:Platform:Country` : [],
        `Gender:Platform:Country` : [],
        ...
},
1 : {
        `Gender:Platform:Country` : [],
        `Gender:Platform:Country` : [],
        ...
},
2 : {
    ....
}
...
99 : {
        `Gender:Platform:Country` : [],
        `Gender:Platform:Country` : [],
        ...
},
all : {
        `Gender:Platform:Country` : [],
        `Gender:Platform:Country` : [],
        ...
}
```
對 Age 做 Partition ，開 101 個 HASH ( all 代表所有 Age ， 如果沒限定 Age 就會從 all 去找 )
- Field Key 設為 `Gender:Platform:Country`
- Value 是一個序列化的 List

**Create Section:**
這邊的 Create 都是指 Pre-load 到 Redis 時的操作
Create 時，會到 [st,ed] 的 key 中加入廣告
> 可以把 Age 很難處理的問題解決
> Create 的過程為了確保 ACID 
> 寫成 Lua Script



因為 HASH 的 
- Key 設為 `Age`
- Field Key 設為 `Gender:Platform:Country`

在 Create 廣告時
很有可能會有 **在同一個 HASH 的 Field Key 衝突** 的問題
所以 Field 的 Value 才會設為 List 來解決

**Query Section**
在查詢時只需透過 `HSCAN` 並搭配 cursor + count ( API 的 offset + limit )
就可以拿到所有符合的廣告ㄌ
( 並且這樣可以在 Redis 層做 Pagenation )

而 TTL 機制則需要額外寫一個 CornJob
會跑過所有 Age 的 HASH 中的 Field 來檢查是否過期
( 用 HSCAN + HDEL )
> 這邊就不用寫 Lua Script 了
> 在執行 Lua Script 時，是 Blocking 的狀態
> 
> 與 Create 不同，不需要保證當下的 ACID
> 如果這次 Delete 到一半的 Age 就突然 Fail 了也沒關係
> 下一輪的 Corn Job 還是會清乾淨
> 達到最終一致性


也會再把當前的 Query 結果以 `ZSET` 加到 Redis
key 設為`cache:Gender:Platform:Country:Age` 
與 Solution 1 相同

:::

:::warning

當沒有指定 Age 的時候???
:::

## Solution 3 

:::info

只對 Search 的結果做 Cache 
直接到 DB 以查詢條件 Search
( 不用 offset + limit 查詢 )

並把整個結果 Cache 到 Redis ( 以 ZSET 存 )
並且不會特別做 Pre-Load 

:::


## Solution 4

:::info


Pre-Load 時以 **最有可能出現的查詢條件** 加入 Redis
廣告 Condition 有: ( 以可能性少到多排序 )
- Gender : 2+1 種可能 ( F , M , * )
- Platform : 3+1 種可能 ( ios , andriod , web , * )
- Country : 這邊假設 10+1 種 ( 假設只針對 10 個國家)
- Age : 100+1 種 ( 1-99 歲 )

==但是 Dcard 的活躍用戶年齡是 **18-35** 歲 !==
所以 Age 這邊只以 [18-5,35+5] 的 Age 區間做 Pre-Load
> ((35+5)-(18-5)+1) = 28 
> 所以是 28+1 = 29 種

總可能數 : `3*4*11*29 = 3828` 個 Key !
> Key 都設為 `F:ios:TW:1` ... ( `Gender:Platform:Country:Age`)
> 多條件時以


所以可以開 3828 個 ZET ， 裡面都存 SQL 查詢結果
因為當天活躍的廣告數至多 1000 個
所以 **最差** 情況 : **3828** 個 ZET，裡面都各存 **1000** String 
( 實際將上述的 Data Pre-Load 到 Redis 只用 345MB )
> 並且這 3828 個 Key 不設 TTL ， 因為都會是主要資料熱點
> String 包含 Title 和 endTime 
> 視情況看要不要 encode 

而遇到 Age 查詢條件在 [13,40] 以外的條件時
會設一個較短 TTL ( 5-10 mins )
同樣不會以 Offset & Limit 去 SQL 查詢
會把整個符合 Condition 的結果存在 ZSET

而 Redis 新增廣告和淘汰的機制同樣使用 CornJob 去把
- 要 active 的廣告 Pre-Laod 到 Redis
- 過期的廣告從 Redis 移除
:::

Pre-Load 完後的 Redis Memory Usage : 
```
> INFO memory
# Memory
used_memory:361795920
used_memory_human:345.04M
```

Corn Job 要跑過所有 Key 時
( 回傳的 value 要帶入下次的 cursor  )
> HSCAN 也通用
```
SCAN 0 COUNT 100
SCAN 1696 COUNT 100
SCAN 1072 COUNT 100
```

為了 Redis 的高可用 會使用 Redis Sentinel 
並分流 Read Request 到各 Redis instance
    
## Idea

- read heavy ( read >> write )
- redis red lock
- mysql > postgresql ( i.e read >> write )
- nosql ?
- primary replica db
- store procedure 
- 隔離節別設定
    - [postgresql 事務隔離級別 ](https://juejin.cn/post/7130952235508301831)
    - [mysql 事務隔離級別 ](https://cloud.tencent.com/developer/article/1833688)
- mysql vs postgresql benckmark : 
    - https://medium.com/@servikash/postgresql-vs-mysql-performance-benchmarking-e2929ee377d4
- 分表
- sharding
    - on MySQL
    - on redis 
        - 多資料庫操作
        - 預設 0-15 個，不互相影響
- 以 redis 分表
- redis cluster 
- redis multi query 
    - lua script 
    - [stackoverflow](https://stackoverflow.com/questions/27469495/multi-criteria-query-with-redis)
    - [redis io : adding-auxiliary-information-in-the-index](https://redis.io/docs/manual/patterns/indexes/#adding-auxiliary-information-in-the-index)
    - [redis index](https://www.slideshare.net/itamarhaber/20140922-redis-tlv-redis-indices)
- [csdn : redis 多條件查詢 + 分頁 ](https://blog.csdn.net/moshowgame/article/details/133563383?utm_medium=distribute.pc_relevant.none-task-blog-2~default~baidujs_utm_term~default-0-133563383-blog-107105950.235^v43^pc_blog_bottom_relevance_base9&spm=1001.2101.3001.4242.1&utm_relevant_index=1)
- [github redis ha deploy](https://github.com/chechiachang/go-redis-ha)
- [redis ha ithome](https://ithelp.ithome.com.tw/articles/10222161?sc=rss.iron)
- [go redis sentinel csdn](https://blog.csdn.net/pengpengzhou/article/details/109363155)
- [go redis sentinel docs](https://redis.uptrace.dev/zh/guide/go-redis-sentinel.html)
- [go redis sentinel client 優化](https://www.jianshu.com/p/a26754b8131f)
- [redis GKE](https://medium.com/@prodriguezdefino/redis-sentinel-on-gke-9c627421b76d)
- [prometheus GKE Cluster](https://medium.com/@selvigp/monitoring-gke-cluster-with-managed-service-for-prometheus-3ecffd61a7c6)
- [Redis大Key优化
](https://blog.csdn.net/lovoo/article/details/131236830?utm_medium=distribute.pc_relevant.none-task-blog-2~default~baidujs_baidulandingword~default-4-131236830-blog-83475960.235^v43^pc_blog_bottom_relevance_base9&spm=1001.2101.3001.4242.3&utm_relevant_index=7)
- [Redis Lua Script : DEL + SCAN](https://juejin.cn/post/6844904117739978766)
- MarshalBinary
- 在 Test 的時後可以用 Dcard 的年齡分佈去模擬
    - https://about.dcard.tw/news/25 ( 18-35 )
    - https://www.similarweb.com/zh-tw/website/dcard.tw/#ranking
- redis lua script deserialize : [cjson](https://onecompiler.com/questions/3x6yguydg/how-to-encode-decode-json-objects-in-redis-lua-script)
- redis `ZSET` select with offest & limit 
    - `ZRANGE sorted_set_test - + BYLEX LIMIT 3 2`
        > offset:3 , limit:2
- [estimate-the-memory-usage](https://lucasmagnum.medium.com/redistip-estimate-the-memory-usage-for-repeated-keys-in-redis-2dc3f163fdab)
- [github (繁中): redis + HAProxy in docker compose](https://github.com/880831ian/docker-compose-redis-sentinel-haproxy?tab=readme-ov-file)
- [ithome + github : HAProxy + K8s ](https://github.com/chechiachang/haproxy-kubernetes)
- [prometheus + redis exporter]
    - [ithome: prometheus](https://ithelp.ithome.com.tw/articles/10224491)
    - [github: prometheus](https://github.com/chechiachang/prometheus-kubernetes?tab=readme-ov-file)
    - [redis exporter](https://github.com/oliver006/redis_exporter)
    - [haproxy exporter](https://grafana.com/grafana/dashboards/367-haproxy-servers-haproxy/)
- [csdn : bitnami redis](https://blog.csdn.net/fly910905/article/details/125474213)
- redis cluster : 
    - medium 中文教學系列
        - [初識 Redis Cluster Ep 1 : Redis Cluster 架構簡介 — 當 Redis 群聚在一起

](https://vicxu.medium.com/%E5%88%9D%E8%AD%98-redis-cluster-ep-1-redis-cluster-%E6%9E%B6%E6%A7%8B%E7%B0%A1%E4%BB%8B-%E7%95%B6-redis-%E7%BE%A4%E8%81%9A%E5%9C%A8%E4%B8%80%E8%B5%B7-67be41e68654)
    - [Redis : Master Slave](https://mp.weixin.qq.com/s/VJTBmAB-A1aRT9DR6v5gow)
    - [Redis Sentinel-深入浅出原理和实战](https://mp.weixin.qq.com/s/k-wGpBBnS53Ap86KNiBYvA)
- [redis 優化配置](https://iter01.com/571384.html)

## Benchmark

- go : request per second
    - https://muetsch.io/http-performance-java-jersey-vs-go-vs-nodejs.html
    - 30000
- redis 
    - for 10000 connection : rps is about 8~90000
    - https://redis.io/docs/management/optimization/benchmarks/#factors-impacting-redis-performance
