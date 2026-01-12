# API 接口完整文档

## 1. 概述

本文档提供选课系统所有 API 接口的完整说明，包括请求格式、响应格式、参数说明和错误处理。

---

## 2. 统一响应格式

### 2.1 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { }
}
```

### 2.2 失败响应

```json
{
  "code": 1,
  "message": "参数不合法"
}
```

### 2.3 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 错误码，0 表示成功 |
| message | string | 提示信息 |
| data | object | 响应数据，可选 |

---

## 3. 认证模块

### 3.1 POST /api/v1/auth/login - 用户登录

**路径**: `POST /api/v1/auth/login`

**权限**: 公开

**请求体**:
```json
{
  "username": "admin",
  "password": "AdminPass123"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 (登录账号) |
| password | string | 是 | 密码 |

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "1"
  }
}
```

**错误响应**:
| code | message | 说明 |
|------|---------|------|
| 1 | 参数不合法 | 参数缺失或格式错误 |
| 4 | 用户不存在 | 用户名不存在 |
| 5 | 密码错误 | 密码不正确 |

---

### 3.2 POST /api/v1/auth/logout - 用户登出

**路径**: `POST /api/v1/auth/logout`

**权限**: 需登录

**请求**: 无需参数，自动从 Cookie 获取 sessionId

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

**错误响应**:
| code | message | 说明 |
|------|---------|------|
| 6 | 用户未登录 | Session 不存在或已过期 |

---

### 3.3 GET /api/v1/auth/whoami - 获取当前用户信息

**路径**: `GET /api/v1/auth/whoami`

**权限**: 需登录

**请求**: 无需参数，从 Session 获取当前用户

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "1",
    "username": "admin",
    "nickname": "管理员",
    "user_type": 1
  }
}
```

**响应字段说明**:
| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | string | 用户ID |
| username | string | 用户名 |
| nickname | string | 昵称 |
| user_type | int | 用户类型 (1=管理员, 2=教师, 3=学生) |

---

## 4. 成员管理模块

### 4.1 GET /api/v1/member - 获取单个成员

**路径**: `GET /api/v1/member`

**权限**: 需登录

**请求参数**:
| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| user_id | string | query | 是 | 用户ID |

**请求示例**:
```
GET /api/v1/member?user_id=1
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "1",
    "username": "admin",
    "nickname": "管理员",
    "user_type": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 4.2 GET /api/v1/member/list - 获取成员列表

**路径**: `GET /api/v1/member/list`

**权限**: 需登录

**请求参数**:
| 参数 | 类型 | 位置 | 必填 | 默认值 | 说明 |
|------|------|------|------|--------|------|
| offset | int | query | 否 | 0 | 分页偏移 |
| limit | int | query | 否 | 10 | 每页数量 |

**请求示例**:
```
GET /api/v1/member/list?offset=0&limit=20
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "user_id": "1",
        "username": "admin",
        "nickname": "管理员",
        "user_type": 1
      },
      {
        "user_id": "2",
        "username": "teacher01",
        "nickname": "张老师",
        "user_type": 2
      }
    ],
    "total": 100
  }
}
```

**响应字段说明**:
| 字段 | 类型 | 说明 |
|------|------|------|
| list | array | 成员列表 |
| list[].user_id | string | 用户ID |
| list[].username | string | 用户名 |
| list[].nickname | string | 昵称 |
| list[].user_type | int | 用户类型 |
| total | int | 总数量 |

---

### 4.3 POST /api/v1/member/create - 创建成员

**路径**: `POST /api/v1/member/create`

**权限**: 管理员

**请求体**:
```json
{
  "username": "newuser",
  "password": "Password123",
  "nickname": "新用户",
  "user_type": 3
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名，唯一 |
| password | string | 是 | 密码，最少8位 |
| nickname | string | 是 | 昵称 |
| user_type | int | 是 | 用户类型 (1=管理员, 2=教师, 3=学生) |

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "100"
  }
}
```

**错误响应**:
| code | message | 说明 |
|------|---------|------|
| 2 | 该 Username 已存在 | 用户名重复 |

---

### 4.4 POST /api/v1/member/update - 更新成员

**路径**: `POST /api/v1/member/update`

**权限**: 管理员

**请求体**:
```json
{
  "user_id": "1",
  "nickname": "新昵称"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | string | 是 | 要更新的用户ID |
| nickname | string | 否 | 新的昵称 |

**成功响应**:
```json
{
  "code": 0,
  "message": "success"
}
```

---

### 4.5 POST /api/v1/member/delete - 删除成员

**路径**: `POST /api/v1/member/delete`

**权限**: 管理员

**请求体**:
```json
{
  "user_id": "1"
}
```

**说明**: 软删除，仅更新 is_deleted 标志

**成功响应**:
```json
{
  "code": 0,
  "message": "success"
}
```

---

## 5. 课程管理模块

### 5.1 GET /api/v1/course/get - 获取课程

**路径**: `GET /api/v1/course/get`

**权限**: 需登录

**请求参数**:
| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| course_id | string | query | 是 | 课程ID |

**请求示例**:
```
GET /api/v1/course/get?course_id=1
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "course_id": "1",
    "name": "高等数学",
    "capacity": 100,
    "cap_selected": 50,
    "teacher_id": "2"
  }
}
```

**响应字段说明**:
| 字段 | 类型 | 说明 |
|------|------|------|
| course_id | string | 课程ID |
| name | string | 课程名称 |
| capacity | int | 课程容量 |
| cap_selected | int | 已选人数 |
| teacher_id | string | 授课教师ID (可选) |

---

### 5.2 POST /api/v1/course/create - 创建课程

**路径**: `POST /api/v1/course/create`

**权限**: 管理员

**请求体**:
```json
{
  "name": "高等数学",
  "cap": 100
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 课程名称 (1-100字符) |
| cap | int | 是 | 课程容量 (必须 > 0) |

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "course_id": "100"
  }
}
```

---

### 5.3 POST /api/v1/course/schedule - 批量排课

**路径**: `POST /api/v1/course/schedule`

**权限**: 管理员

**请求体**:
```json
{
  "teacher_course_relationship": {
    "2": ["1", "2", "3"],
    "3": ["4", "5"]
  }
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| teacher_course_relationship | object | 是 | 教师课程映射关系 |

**成功响应**:
```json
{
  "code": 0,
  "message": "success"
}
```

---

## 6. 教师管理模块

### 6.1 GET /api/v1/teacher/get_course - 获取教师课程

**路径**: `GET /api/v1/teacher/get_course`

**权限**: 需登录

**请求参数**:
| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| teacher_id | string | query | 是 | 教师ID |

**请求示例**:
```
GET /api/v1/teacher/get_course?teacher_id=2
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "course_id": "1",
      "name": "高等数学",
      "capacity": 100,
      "teacher_id": "2"
    }
  ]
}
```

---

### 6.2 POST /api/v1/teacher/bind_course - 绑定课程

**路径**: `POST /api/v1/teacher/bind_course`

**权限**: 管理员

**请求体**:
```json
{
  "course_id": "1",
  "teacher_id": "2"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| course_id | string | 是 | 课程ID |
| teacher_id | string | 是 | 教师ID |

**成功响应**:
```json
{
  "code": 0,
  "message": "success"
}
```

**错误响应**:
| code | message | 说明 |
|------|---------|------|
| 12 | 课程不存在 | course_id 不存在 |
| 8 | 课程已绑定过 | 课程已被其他教师绑定 |

---

### 6.3 POST /api/v1/teacher/unbind_course - 解绑课程

**路径**: `POST /api/v1/teacher/unbind_course`

**权限**: 管理员

**请求体**:
```json
{
  "course_id": "1",
  "teacher_id": "2"
}
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success"
}
```

**错误响应**:
| code | message | 说明 |
|------|---------|------|
| 9 | 课程未绑定过 | 课程未绑定到任何教师 |
| 10 | 没有操作权限 | 不是绑定该课程的教师 |

---

## 7. 学生选课模块

### 7.1 POST /api/v1/student/book_course - 选课

**路径**: `POST /api/v1/student/book_course`

**权限**: 需登录 (学生)

**请求体**:
```json
{
  "student_id": "4",
  "course_id": "1"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| student_id | string | 是 | 学生ID |
| course_id | string | 是 | 课程ID |

**成功响应**:
```json
{
  "code": 0,
  "message": "success"
}
```

**错误响应**:
| code | message | 说明 |
|------|---------|------|
| 7 | 课程已满 | 课程容量已满 |
| 15 | 重复请求 | 学生已选过该课程 |
| 12 | 课程不存在 | course_id 不存在 |

---

### 7.2 GET /api/v1/student/course - 获取课表

**路径**: `GET /api/v1/student/course`

**权限**: 需登录

**请求参数**:
| 参数 | 类型 | 位置 | 必填 | 说明 |
|------|------|------|------|------|
| student_id | string | query | 是 | 学生ID |

**请求示例**:
```
GET /api/v1/student/course?student_id=4
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "course_id": "1",
      "name": "高等数学",
      "teacher_id": "2"
    },
    {
      "course_id": "2",
      "name": "大学物理",
      "teacher_id": "3"
    }
  ]
}
```

---

## 8. 健康检查

### 8.1 GET /health - 健康检查

**路径**: `GET /health`

**权限**: 公开

**成功响应**:
```json
{
  "status": "ok"
}
```

---

## 9. 监控指标

### 9.1 GET /metrics - Prometheus 指标

**路径**: `GET /metrics`

**权限**: 公开

**响应**: Prometheus 格式的监控指标

---

## 10. 权限一览表

| 接口 | 方法 | 路径 | 权限 |
|------|------|------|------|
| 登录 | POST | /api/v1/auth/login | 公开 |
| 登出 | POST | /api/v1/auth/logout | 需登录 |
| 当前用户 | GET | /api/v1/auth/whoami | 需登录 |
| 获取成员 | GET | /api/v1/member | 需登录 |
| 成员列表 | GET | /api/v1/member/list | 需登录 |
| 创建成员 | POST | /api/v1/member/create | 管理员 |
| 更新成员 | POST | /api/v1/member/update | 管理员 |
| 删除成员 | POST | /api/v1/member/delete | 管理员 |
| 获取课程 | GET | /api/v1/course/get | 需登录 |
| 创建课程 | POST | /api/v1/course/create | 管理员 |
| 批量排课 | POST | /api/v1/course/schedule | 管理员 |
| 教师课程 | GET | /api/v1/teacher/get_course | 需登录 |
| 绑定课程 | POST | /api/v1/teacher/bind_course | 管理员 |
| 解绑课程 | POST | /api/v1/teacher/unbind_course | 管理员 |
| 选课 | POST | /api/v1/student/book_course | 需登录 |
| 课表 | GET | /api/v1/student/course | 需登录 |
