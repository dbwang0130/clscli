---
name: clscli
description: 检索分析腾讯云CLS日志
homepage: https://github.com/
metadata:
    {"requires": {"bin": ["clscli"], "env": ["TENCENTCLOUD_SECRET_ID", "TENCENTCLOUD_SECRET_KEY"]}}
---

# CLS Skill

检索分析腾讯云CLS日志

## Setup
1. 获取凭证与地域：https://cloud.tencent.com/document/api/614/56474
2. 设置环境变量（与腾讯云 API 公共参数一致）：
    ```bash
    export TENCENTCLOUD_SECRET_ID="your-secret-id"
    export TENCENTCLOUD_SECRET_KEY="your-secret-key"
    ```
3. 地域通过命令行参数 `--region` 指定（如 ap-guangzhou）。

## Usage

!IMPORTANT: 如果不知道日志TOPIC，先查询topics列表。

### List log topics
List topics in a region to determine which `--region` and topic ID to use for query/context.

```bash
clscli topics --region <region> [--topic-name name] [--logset-name name] [--logset-id id] [--limit 20] [--offset 0]
```
示例：`--output=json`、`--output=csv`、`-o topics.csv`

参数说明
| 参数 | 必选 | 说明 |
|------|------|------|
| --region | 是 | CLS 地域，如 ap-guangzhou |
| --topic-name | 否 | 按主题名称过滤（模糊匹配） |
| --logset-name | 否 | 按日志集名称过滤（模糊匹配） |
| --logset-id | 否 | 按日志集 ID 过滤 |
| --limit | 否 | 单页条数，默认 20，最大 100 |
| --offset | 否 | 分页偏移，默认 0 |
| --output, -o | 否 | 输出：json、csv，或文件路径 |

输出列：Region, TopicId, TopicName, LogsetId, CreateTime, StorageType。

### Get log by query
```bash
clscli query -q "[检索条件] | [SQL 语句]" --region <region> -t <TopicId> --last 1h
```
示例：
- 时间：`--last 1h`、`--last 30m`；或 `--from`/`--to`（Unix 毫秒）
- 多主题：`--topics <id1>,<id2>` 或多次 `-t <id>`
- 自动翻页与上限：`--max 5000`（自动翻页直到累计 5000 条或结束）
- 输出：`--output=json`、`--output=csv`、`-o result.json`（输出到文件）

参数说明
| 参数 | 必选 | 说明 |
|------|------|------|
| --region | 是 | CLS 地域，如 ap-guangzhou |
| -q, --query | 是 | 检索条件或 SQL，如 `level:ERROR` 或 `* \| select count(*) as cnt` |
| -t, --topic | 与 --topics 二选一 | 单日志主题 ID |
| --topics | 与 -t 二选一 | 逗号分隔的多个主题 ID，最多 50 个 |
| --last | 与 --from/--to 二选一 | 时间范围，如 1h、30m、24h |
| --from, --to | 与 --last 二选一 | 起止时间（Unix 毫秒） |
| --limit | 否 | 单次请求条数，默认 100，最大 1000 |
| --max | 否 | 累计最大条数，非 0 时自动翻页直到达到或 ListOver |
| --output, -o | 否 | 输出：json、csv，或文件路径 |
| --sort | 否 | 排序：asc 或 desc，默认 desc |
#### 检索条件

检索条件支持两种语法规则
CQL：CLS Query Language，日志服务 CLS 专用检索语法，专为日志检索设计，使用容易，推荐使用。
Lucene：开源 Lucene 语法，由于该语法并非专为日志检索设计，对特殊符号、大小写、通配符等有较多限制，使用较为繁琐，容易出现语法错误，不推荐使用。

##### CQL语法规则
| 语法       | 说明 |
|------------|------|
| `key:value` | 键值检索，查询字段（key）的值中包含 value 的日志，例如：`level:ERROR` |
| `value`     | 全文检索，查询日志全文中包含 value 的日志，例如：`ERROR` |
| `AND`       | “与”逻辑操作符，不区分大小写，例如：`level:ERROR AND pid:1234` |
| `OR`        | “或”逻辑操作符，不区分大小写，例如：`level:ERROR OR level:WARNING`，`level:(ERROR OR WARNING)` |
| `NOT`       | “非”逻辑操作符，不区分大小写，例如：`level:ERROR NOT pid:1234`，`level:ERROR AND NOT pid:1234` |
| `()`        | 逻辑分组操作符，控制逻辑运算优先级，例如：`level:(ERROR OR WARNING) AND pid:1234`<br/>**注意：未使用括号时，AND 优先级高于 OR**。 |
| `"  "`      | 短语检索，使用双引号包裹一个字符串，日志需包含字符串内的各个词，且各个词的顺序保持不变，例如：`name:"john Smith"`<br/>短语检索中不存在逻辑操作符，其等同于查询字符本身，例如：`name:"and"` |
| `'  '`      | 短语检索，使用单引号包裹一个字符串，功能等价于 `""`，当被检索短语中包含双引号时，可使用单引号包裹该短语，以避免语法错误，例如：`body:'user_name:"bob"'` |
| `*`         | 模糊检索，匹配零个、单个、多个字符，例如：`host:www.test*.com`，不支持前缀模糊检索 |
| `>`         | 范围操作符，表示大于某个数值，例如：`status>400` 或 `status:>400` |
| `>=`        | 范围操作符，表示大于等于某个数值，例如：`status>=400` 或 `status:>=400` |
| `<`         | 范围操作符，表示小于某个数值，例如：`status<400` 或 `status:<400` |
| `<=`        | 范围操作符，表示小于等于某个数值，例如：`status<=400` 或 `status:<=400` |
| `=`         | 范围操作符，表示等于某个数值，例如：`status=400`，等价于 `status:400` |
| `\`         | 转义符号，转义后的字符表示符号本身。<br/>被检索的值包含空格、:、(、)、>、=、<、"、'、* 时，需进行转义，例如：`body:user_name\:bob`<br/>使用双引号进行短语检索时，仅需转义 `"` 及 `*`<br/>使用单引号进行短语检索时，仅需转义 `'` 及 `*`<br/>未转义的 `*` 代表模糊检索 |
| `key:*`     | text 类型字段：查询字段（key）存在的日志，无论值是否为空，例如：`url:*`<br/>long/double 类型字段：查询字段（key）存在，且值为数值的日志，例如：`response_time:*` |
| `key:""`    | text 类型字段：查询字段（key）存在且值为空的日志，值仅包含分词符时也等价为空，例如：`url:""`<br/>long/double 类型字段：查询字段值不为数值的日志，包含字段（key）不存在的情况，例如：`response_time:""` |

#### SQL语句语法
| 语法 | 说明 |
|------|------|
| SELECT | 从表中选取数据，默认从当前日志主题中获取符合检索条件的数据 |
| AS | 为列名称（KEY）指定别名 |
| GROUP BY | 结合聚合函数，根据一个或多个列（KEY）对原始数据进行分组聚合 |
| ORDER BY | 根据指定的 KEY 对结果集进行排序 |
| LIMIT | 限制结果集数据行数，默认限制为100，最大100万 |
| WHERE | 对查询到的原始数据进行过滤 |
| HAVING | 对分组聚合后的数据进行过滤，与 WHERE 的区别在于其作用于分组（GROUP BY）之后，排序（ORDER BY）之前，而 WHERE 作用于聚合前的原始数据 |
| 嵌套子查询 | 针对一些复杂的统计分析场景，需要先对原始数据进行一次统计分析，再针对该分析结果进行二次统计分析，这时候需要在一个 SELECT 语句中嵌套另一个 SELECT 语句，这种查询方式称为嵌套子查询 |
| SQL 函数 | 使用 SQL 函数对日志进行更丰富的分析处理，例如从 IP 解析地理信息、时间格式转换、字符串分隔及连接、JSON 字段提取、算数运算、统计唯一值个数等 |


### Describe log context 上线文检索

用于搜索日志上下文附近的内容

```bash
clscli context <PkgId> <PkgLogId> --region <region> -t <TopicId>
```
示例：`--output=json`、`--output=csv`、`-o context.json`（输出到文件）

参数说明
| 参数 | 必选 | 类型 | 说明 | 示例值 |
|------|------|------|------|--------|
| --region | 是 | String | CLS 地域 | ap-guangzhou |
| -t, --topic | 是 | String | 日志主题 ID | - |
| PkgId | 是 | String | 日志包序号，即 SearchLog 返回的 Results[].PkgId | 528C1318606EFEB8-1A7 |
| PkgLogId | 是 | Integer | 日志包内序号，即 SearchLog 返回的 Results[].PkgLogId | 65536 |
| --output, -o | 否 | - | 输出：json、csv，或文件路径 | - |