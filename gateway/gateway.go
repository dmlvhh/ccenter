package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero2/initialize"
	"go-zero2/system/system_models"
	"go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Config 配置文件
type Config struct {
	Addr string
	Etcd string
	Auth struct {
		AccessSecret string
		AccessExpire int
	}
	Mysql struct {
		DataSource string
	}
	Log logx.LogConf
}

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

var configFile = flag.String("f", "settings.yaml", "the config file")
var config Config
var DB *gorm.DB

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli   *clientv3.Client
	addrs map[string][]string
	mu    sync.RWMutex
}
type ServiceContext struct {
	Config Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewServiceDiscovery(endpoints []string) (*ServiceDiscovery, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	sd := &ServiceDiscovery{
		cli:   cli,
		addrs: make(map[string][]string),
	}
	go sd.watchServices()
	return sd, nil
}

func (sd *ServiceDiscovery) watchServices() {
	watcher := clientv3.NewWatcher(sd.cli)
	for {
		watchChan := watcher.Watch(context.Background(), "services/", clientv3.WithPrefix())
		for resp := range watchChan {
			if resp.Err() != nil {
				logx.Error("Watch services failed:", resp.Err())
				continue
			}

			for _, event := range resp.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					sd.mu.Lock()
					sd.addrs[string(event.Kv.Key)] = append(sd.addrs[string(event.Kv.Key)], string(event.Kv.Value))
					sd.mu.Unlock()
				case clientv3.EventTypeDelete:
					sd.mu.Lock()
					delete(sd.addrs, string(event.Kv.Key))
					sd.mu.Unlock()
				}
			}
		}
	}
}

func (sd *ServiceDiscovery) GetServiceAddrs(service string) []string {
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	return sd.addrs[service]
}

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	addrs []string
	index int
	mu    sync.Mutex
}

func NewLoadBalancer(addrs []string) *LoadBalancer {
	return &LoadBalancer{
		addrs: addrs,
		index: 0,
	}
}

func (lb *LoadBalancer) Next() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.addrs) == 0 {
		return ""
	}

	addr := lb.addrs[lb.index]
	lb.index = (lb.index + 1) % len(lb.addrs)
	return addr
}

// Proxy 反向代理
type Proxy struct{}

func (Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// 从请求路径中提取服务名称
	regex, _ := regexp.Compile(`/api/([^/]+)`)
	addrList := regex.FindStringSubmatch(req.URL.Path)
	if len(addrList) != 2 {
		FileResponse("请求路径错误", res)
		return
	}
	service := addrList[1]

	fmt.Println("service", service)
	// 从服务发现获取实例地址
	addrs := sd.GetServiceAddrs(service + "_api")
	if len(addrs) == 0 {
		logx.Errorf("服务 %s 不可用", service)
		FileResponse(fmt.Sprintf("服务 %s 不可用", service), res)
		return
	}

	// 使用负载均衡选择一个地址
	lb := NewLoadBalancer(addrs)
	addr := lb.Next()

	// 打印请求信息
	remoteAddr := strings.Split(req.RemoteAddr, ":")
	logx.Infof("请求服务: %s, 请求IP: %s, 请求URL: %s", service, remoteAddr[0], req.URL.String())

	// 认证拦截
	if !auth(req) {
		FileResponse("认证失败", res)
		return
	}

	// 转发请求
	remote, _ := url.Parse("http://" + addr)
	reverseProxy := httputil.NewSingleHostReverseProxy(remote)
	reverseProxy.ServeHTTP(res, req)
}

// auth 认证拦截
func auth(req *http.Request) bool {
	// 从请求头中获取子系统的 token
	token := req.Header.Get("Authorization")
	if token == "" {
		return false
	}

	// 解析 token，获取子系统 ID
	claims, err := ParseToken(token, config.Auth.AccessSecret)
	if err != nil {
		return false
	}

	// 检查子系统是否被主系统授权
	if !isSubSystemAuthorized(claims.SubSystemID) {
		return false
	}

	// 检查子系统的授权是否在有效期内
	if isSubSystemExpired(claims.SubSystemID) {
		return false
	}

	return true
}

// JWYPayLoad 是 JWT 载荷的结构体，包含用户相关信息
type JWYPayLoad struct {
	SubSystemID    string `json:"sub_system_id"`
	SubSystemName  string `json:"sub_system_name"`
	ServerAddress  string `json:"server_address"`
	DatabaseConfig string `json:"database_config"`
}

// CustomClaims 是自定义声明的结构体，包含 JWYPayLoad 和 JWT 标准声明 RegisteredClaims
type CustomClaims struct {
	JWYPayLoad
	jwt.RegisteredClaims
}

// ParseToken 函数用于解析 JWT Token，并返回 CustomClaims 结构体和可能的错误
func ParseToken(tokenStr string, accessSecret string) (*CustomClaims, error) {
	// 使用 jwt 包的 ParseWithClaims 函数解析 Token
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	// 检查解析过程中是否出现错误
	if err != nil {
		return nil, err
	}
	// 检查 Token 是否有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	// 如果 Token 有效，返回 CustomClaims 结构体
	return nil, errors.New("invalid token")
}

// isSubSystemAuthorized 检查子系统是否被授权
func isSubSystemAuthorized(subSystemID string) bool {
	var count int64
	err := DB.Take(&system_models.SubSystem{}, subSystemID).Count(&count).Error
	if count <= 0 {
		logx.Error(err)
		return false
	}
	// 查询数据库或缓存，检查子系统是否被授权
	return true // 示例代码，实际需要实现
}

// isSubSystemExpired 检查子系统的授权是否过期
func isSubSystemExpired(subSystemID string) bool {
	// 查询数据库或缓存，检查子系统的授权是否过期
	var subS system_models.Authorization
	err := DB.Model(&system_models.Authorization{}).Where("sub_system_id=?", subSystemID).First(&subS).Error
	if err != nil {
		logx.Error(err)
		return true
	}
	// 检查授权是否过期
	if time.Now().After(subS.ExpiresAt) {
		logx.Infof("子系统 %s 授权已过期", subSystemID)
		return true
	}
	return false // 示例代码，实际需要实现
}

// FileResponse 返回 JSON 响应
func FileResponse(msg string, res http.ResponseWriter) {
	response := BaseResponse{Code: 0, Msg: msg, Data: map[string]interface{}{}}
	byteData, _ := json.Marshal(&response)
	res.Write(byteData)
}

var sd *ServiceDiscovery

func main() {
	flag.Parse()
	conf.MustLoad(*configFile, &config)
	logx.SetUp(config.Log)

	var err error
	sd, err = NewServiceDiscovery(strings.Split(config.Etcd, ","))
	if err != nil {
		logx.Error("初始化服务发现失败:", err)
		return
	}
	DB = initialize.InitMysql(config.Mysql.DataSource)
	Proxy := Proxy{}
	logx.Infof("网关运行地址: %s\n", config.Addr)
	http.ListenAndServe(config.Addr, Proxy)
}
