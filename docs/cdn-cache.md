# CDN缓存配置指南

## 概述

本文档详细说明了Docker镜像代理系统的CDN缓存配置策略，以优化镜像拉取性能和减少上游带宽消耗。

## 缓存策略概览

### 1. 长期缓存 (1年)

**路径模式:** `/v2/*/blobs/*`

**原因:** Blob数据是不可变的，基于内容哈希寻址，永远不会改变。

**示例URL:**
```
/v2/library/nginx/blobs/sha256:abc123...
/v2/helloz/myapp/blobs/sha256:def456...
```

**缓存头设置:**
```
Cache-Control: public, max-age=31536000
ETag: "sha256:abc123..."
```

### 2. 短期缓存 (5分钟)

**路径模式:** `/v2/*/manifests/*`

**原因:** Manifest文件可能会被更新，但更新频率不高。

**示例URL:**
```
/v2/library/nginx/manifests/latest
/v2/library/nginx/manifests/sha256:xyz789...
```

**缓存头设置:**
```
Cache-Control: public, max-age=300
ETag: "sha256:xyz789..."
```

### 3. 不缓存

**路径模式:**
- `/v2/` (API版本检查)
- `/v2/*/tags/list` (标签列表，可能频繁变化)
- `/api/*` (管理API)
- `/static/*` (WEB界面静态文件)

**缓存头设置:**
```
Cache-Control: no-cache, no-store, must-revalidate
```

## 主流CDN配置示例

### 1. Cloudflare

#### Page Rules配置

```
# 长期缓存 Blob 数据
URL: *example.com/v2/*/blobs/*
Settings:
  - Cache Level: Cache Everything
  - Edge Cache TTL: 1 year
  - Browser Cache TTL: 1 year

# 短期缓存 Manifest 数据
URL: *example.com/v2/*/manifests/*
Settings:
  - Cache Level: Cache Everything
  - Edge Cache TTL: 5 minutes
  - Browser Cache TTL: 5 minutes

# 不缓存 API 接口
URL: *example.com/api/*
Settings:
  - Cache Level: Bypass
```

#### Workers脚本示例

```javascript
addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})

async function handleRequest(request) {
  const url = new URL(request.url)
  const path = url.pathname
  
  // 长期缓存 blob 数据
  if (path.match(/^\/v2\/.*\/blobs\//)) {
    const cacheKey = new Request(url.toString(), request)
    const cache = caches.default
    let response = await cache.match(cacheKey)
    
    if (!response) {
      response = await fetch(request)
      if (response.status === 200) {
        response = new Response(response.body, {
          status: response.status,
          statusText: response.statusText,
          headers: {
            ...response.headers,
            'Cache-Control': 'public, max-age=31536000',
          },
        })
        event.waitUntil(cache.put(cacheKey, response.clone()))
      }
    }
    return response
  }
  
  // 短期缓存 manifest 数据
  if (path.match(/^\/v2\/.*\/manifests\//)) {
    const cacheKey = new Request(url.toString(), request)
    const cache = caches.default
    let response = await cache.match(cacheKey)
    
    if (!response) {
      response = await fetch(request)
      if (response.status === 200) {
        response = new Response(response.body, {
          status: response.status,
          statusText: response.statusText,
          headers: {
            ...response.headers,
            'Cache-Control': 'public, max-age=300',
          },
        })
        event.waitUntil(cache.put(cacheKey, response.clone()))
      }
    }
    return response
  }
  
  // 其他请求直接转发
  return fetch(request)
}
```

### 2. AWS CloudFront

#### Behaviors配置

```yaml
# Blob 数据缓存行为
- PathPattern: "/v2/*/blobs/*"
  TargetOriginId: zmirror-origin
  ViewerProtocolPolicy: redirect-to-https
  CachePolicyId: 4135ea2d-6df8-44a3-9df3-4b5a84be39ad  # CachingOptimized
  TTL:
    DefaultTTL: 31536000  # 1 year
    MaxTTL: 31536000      # 1 year

# Manifest 数据缓存行为
- PathPattern: "/v2/*/manifests/*"
  TargetOriginId: zmirror-origin
  ViewerProtocolPolicy: redirect-to-https
  CachePolicyId: custom-manifest-policy
  TTL:
    DefaultTTL: 300       # 5 minutes
    MaxTTL: 300           # 5 minutes

# API 接口不缓存
- PathPattern: "/api/*"
  TargetOriginId: zmirror-origin
  ViewerProtocolPolicy: redirect-to-https
  CachePolicyId: 4135ea2d-6df8-44a3-9df3-4b5a84be39ad  # CachingDisabled
```

#### Lambda@Edge函数

```javascript
exports.handler = (event, context, callback) => {
    const request = event.Records[0].cf.request;
    const uri = request.uri;
    
    // 添加缓存头
    if (uri.match(/^\/v2\/.*\/blobs\//)) {
        // Blob 数据长期缓存
        const response = {
            status: '200',
            statusDescription: 'OK',
            headers: {
                'cache-control': [{
                    key: 'Cache-Control',
                    value: 'public, max-age=31536000'
                }],
                'etag': [{
                    key: 'ETag',
                    value: `"${uri.split('/').pop()}"`
                }]
            }
        };
        callback(null, response);
    } else if (uri.match(/^\/v2\/.*\/manifests\//)) {
        // Manifest 数据短期缓存
        const response = {
            status: '200',
            statusDescription: 'OK',
            headers: {
                'cache-control': [{
                    key: 'Cache-Control',
                    value: 'public, max-age=300'
                }]
            }
        };
        callback(null, response);
    }
    
    callback(null, request);
};
```

### 3. 阿里云CDN

#### 域名配置

```yaml
# 缓存规则配置
CacheRules:
  - PathPattern: "/v2/*/blobs/*"
    CacheTTL: 31536000  # 1年
    CacheType: "file"
    
  - PathPattern: "/v2/*/manifests/*"
    CacheTTL: 300       # 5分钟
    CacheType: "file"
    
  - PathPattern: "/api/*"
    CacheTTL: 0         # 不缓存
    CacheType: "no-cache"

# 压缩配置
Compression:
  Enable: true
  Types: ["application/json", "text/plain"]

# 回源配置
Origin:
  Type: "domain"
  Domain: "your-zmirror-server.com"
  Port: 8080
  Protocol: "http"
```

### 4. 腾讯云CDN

#### 缓存配置规则

```json
{
  "cacheRules": [
    {
      "ruleType": "file",
      "rulePaths": ["/v2/*/blobs/*"],
      "cacheTime": 31536000,
      "ignoreCase": false
    },
    {
      "ruleType": "file", 
      "rulePaths": ["/v2/*/manifests/*"],
      "cacheTime": 300,
      "ignoreCase": false
    },
    {
      "ruleType": "file",
      "rulePaths": ["/api/*"],
      "cacheTime": 0,
      "ignoreCase": false
    }
  ]
}
```

## Nginx反向代理配置

### 完整配置示例

```nginx
# 上游服务器配置
upstream zmirror_backend {
    server zmirror:8080;
    keepalive 32;
}

# 缓存配置
proxy_cache_path /var/cache/nginx/zmirror 
                 levels=1:2 
                 keys_zone=zmirror:10m 
                 max_size=10g 
                 inactive=1y 
                 use_temp_path=off;

server {
    listen 80;
    server_name your-registry.example.com;
    
    # 重定向到HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-registry.example.com;
    
    # SSL配置
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # 基础代理配置
    location / {
        proxy_pass http://zmirror_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
    }
    
    # Blob 数据 - 长期缓存
    location ~ ^/v2/.*/blobs/ {
        proxy_pass http://zmirror_backend;
        proxy_cache zmirror;
        proxy_cache_valid 200 1y;
        proxy_cache_key $uri;
        proxy_cache_lock on;
        proxy_cache_use_stale error timeout invalid_header updating;
        
        # 添加缓存状态头
        add_header X-Cache-Status $upstream_cache_status;
        add_header Cache-Control "public, max-age=31536000";
        
        # 忽略客户端的 no-cache 请求
        proxy_ignore_headers Cache-Control Expires;
        proxy_cache_bypass 0;
        proxy_no_cache 0;
    }
    
    # Manifest 数据 - 短期缓存
    location ~ ^/v2/.*/manifests/ {
        proxy_pass http://zmirror_backend;
        proxy_cache zmirror;
        proxy_cache_valid 200 5m;
        proxy_cache_key $uri;
        proxy_cache_lock on;
        
        add_header X-Cache-Status $upstream_cache_status;
        add_header Cache-Control "public, max-age=300";
    }
    
    # API 接口 - 不缓存
    location /api/ {
        proxy_pass http://zmirror_backend;
        proxy_cache off;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
        add_header Pragma "no-cache";
        add_header Expires "0";
    }
    
    # 管理界面 - 不缓存
    location /static/ {
        proxy_pass http://zmirror_backend;
        proxy_cache off;
    }
    
    # Docker Registry API版本检查 - 不缓存
    location = /v2/ {
        proxy_pass http://zmirror_backend;
        proxy_cache off;
    }
    
    # 标签列表 - 不缓存（可能频繁变化）
    location ~ ^/v2/.*/tags/list$ {
        proxy_pass http://zmirror_backend;
        proxy_cache off;
    }
}
```

### 缓存清理脚本

```bash
#!/bin/bash
# clear_cache.sh - 清理特定镜像的缓存

IMAGE_NAME=$1
if [ -z "$IMAGE_NAME" ]; then
    echo "Usage: $0 <image_name>"
    echo "Example: $0 library/nginx"
    exit 1
fi

# 清理manifest缓存
find /var/cache/nginx/zmirror -name "*" -path "*v2*${IMAGE_NAME}*manifests*" -delete

# 可选：清理blob缓存（谨慎使用，blob数据很大且很少变化）
# find /var/cache/nginx/zmirror -name "*" -path "*v2*${IMAGE_NAME}*blobs*" -delete

echo "Cache cleared for image: $IMAGE_NAME"

# 重载Nginx配置
nginx -s reload
```

## 监控和维护

### 1. 缓存命中率监控

```bash
#!/bin/bash
# monitor_cache.sh - 监控缓存命中率

# 分析Nginx访问日志
tail -f /var/log/nginx/access.log | grep -E "(HIT|MISS|EXPIRED)" | \
while read line; do
    echo "$(date): $line"
done
```

### 2. 缓存大小监控

```bash
#!/bin/bash
# check_cache_size.sh - 检查缓存目录大小

CACHE_DIR="/var/cache/nginx/zmirror"
CACHE_SIZE=$(du -sh $CACHE_DIR | cut -f1)
echo "Cache size: $CACHE_SIZE"

# 如果缓存过大，清理旧的缓存文件
MAX_SIZE_GB=50
CURRENT_SIZE_GB=$(du -s $CACHE_DIR | awk '{print int($1/1024/1024)}')

if [ $CURRENT_SIZE_GB -gt $MAX_SIZE_GB ]; then
    echo "Cache size exceeds ${MAX_SIZE_GB}GB, cleaning old files..."
    find $CACHE_DIR -type f -atime +7 -delete
fi
```

### 3. 性能优化建议

1. **预热常用镜像**
   ```bash
   # 预热脚本
   POPULAR_IMAGES=(
     "library/nginx:latest"
     "library/ubuntu:latest" 
     "library/node:latest"
   )
   
   for image in "${POPULAR_IMAGES[@]}"; do
     docker pull your-registry.com/$image
   done
   ```

2. **监控上游延迟**
   ```bash
   # 检查上游响应时间
   curl -w "@curl-format.txt" -o /dev/null -s "https://registry-1.docker.io/v2/"
   ```

3. **调整缓存参数**
   - 根据实际使用情况调整缓存大小
   - 监控磁盘使用率
   - 优化缓存键策略

## 故障排除

### 常见问题

1. **缓存不生效**
   - 检查路径匹配规则
   - 验证响应头设置
   - 查看CDN配置日志

2. **缓存过期问题**
   - 清理特定镜像缓存
   - 检查TTL设置
   - 验证ETag处理

3. **性能问题**
   - 监控缓存命中率
   - 优化缓存策略
   - 增加缓存容量

### 调试工具

```bash
# 检查响应头
curl -I https://your-registry.com/v2/library/nginx/manifests/latest

# 测试缓存效果
curl -w "%{http_code} %{time_total}s\n" -o /dev/null -s \
  https://your-registry.com/v2/library/nginx/blobs/sha256:abc123...

# 查看CDN统计
# (根据具体CDN提供商的API文档)
```

通过以上配置，可以显著提升Docker镜像拉取的性能，减少上游带宽消耗，为用户提供更好的体验。
