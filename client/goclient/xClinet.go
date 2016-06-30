package goclient

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/coreos/etcd/client"

	"io/ioutil"

	"os"

	"net"

	"github.com/sosop/libconfig"
	"golang.org/x/net/context"
)

var (
	api     client.KeysAPI
	prePath string
	prjName string
	env     string
	isWeb   bool
	entry   = make(map[string]string, 128)
	// IniConf 读取配置
	IniConf *libconfig.IniConfig
)

func init() {
	IniConf = &libconfig.IniConfig{Entry: make(map[string]string, 32)}
}

// Config 初始化设置
func Config() {
	isWeb = IniConf.GetBool("web", false)

	clientUrlsStr := IniConf.GetString("etcd_clinet_urls", "http://127.0.0.1:2379")
	clientUrls := strings.Split(clientUrlsStr, ",")
	newKeysAPI(clientUrls)

	if !isWeb {
		prjName = IniConf.GetString("prjName", "prjName")
		env = IniConf.GetString("env", "prod")
		prePath = MakeKey(prjName, env) + "/"

		err := pullAll()
		for i := 0; i < 3 && err != nil; i++ {
			log.Println(err)
			time.Sleep(time.Second)
			err = pullAll()
		}
		if err != nil {
			if err = readFromDump(); err != nil {
				panic(err)
			}
		} else {
			Dump()
		}
		info := getLocalInfo()
		go heart(info)
	}
}

// SetInit 设置初始化
func SetInit(key, value string) {
	IniConf.Set(key, value)
}

// newnewKeysApi 创建keyapi
func newKeysAPI(clientUrls []string) {
	cfg := client.Config{
		Endpoints: clientUrls,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {

	}
	api = client.NewKeysAPI(c)
}

// Get retrieves a set of Nodes from etcd
func Get(key string, opts *client.GetOptions) (*client.Response, error) {
	return api.Get(context.Background(), prePath+key, opts)
}

// Set assigns a new value to a Node identified by a given key. The caller
// may define a set of conditions in the SetOptions. If SetOptions.Dir=true
// then value is ignored.
func Set(key, value string, opts *client.SetOptions) (*client.Response, error) {
	return api.Set(context.Background(), prePath+key, value, opts)
}

// Delete removes a Node identified by the given key, optionally destroying
// all of its children as well. The caller may define a set of required
// conditions in an DeleteOptions object.
func Delete(key string, opts *client.DeleteOptions) (*client.Response, error) {
	return api.Delete(context.Background(), prePath+key, opts)
}

// Create is an alias for Set w/ PrevExist=false
func Create(key, value string) (*client.Response, error) {
	return api.Create(context.Background(), prePath+key, value)
}

// CreateInOrder is used to atomically create in-order keys within the given directory.
func CreateInOrder(dir, value string, opts *client.CreateInOrderOptions) (*client.Response, error) {
	return api.CreateInOrder(context.Background(), prePath+dir, value, opts)
}

// CreateDir is used to atomically create in-order keys within the given directory.
func CreateDir(dir string) error {
	key := prePath + dir + "/1"
	_, err := api.Create(context.Background(), key, "1")
	if err != nil {
		return err
	}
	_, err = Delete(key, nil)
	return err
}

// Update is an alias for Set w/ PrevExist=true
func Update(key, value string) (*client.Response, error) {
	return api.Update(context.Background(), prePath+key, value)
}

// Watcher builds a new Watcher targeted at a specific Node identified
// by the given key. The Watcher may be configured at creation time
// through a WatcherOptions object. The returned Watcher is designed
// to emit events that happen to a Node, and optionally to its children.
func Watcher(key string, opts *client.WatcherOptions) client.Watcher {
	return api.Watcher(prePath+key, opts)
}

// Dump 持久化本地
func Dump() error {
	confs := ""
	for k, v := range entry {
		confs += fmt.Sprintln(k, "=", v)
	}
	err := ioutil.WriteFile(IniConf.GetString("dunpPath", "confs.dump"), []byte(confs), 0666)
	if err != nil {
		return err
	}
	return nil
}

// Watching 监听节点变化
func Watching(f func()) {
	go watching(f, "/publish"+prePath)
}

// WatchingShare 监听共享节点节点变化
func WatchingShare(f func()) {
	go watching(f, "/publish/share/"+env)
}

func watching(f func(), path string) {
	for {
		resp, err := api.Watcher(path, &client.WatcherOptions{Recursive: true}).Next(context.Background())
		if err != nil {
			time.Sleep(time.Millisecond * 100)
			log.Println(err)
			continue
		}
		if strings.ToLower(resp.Action) == "update" {
			err = pullAll()
			if err != nil {
				log.Println(err)
			}
			err = Dump()
			if err != nil {
				log.Println(err)
			}
			f()
		}
	}
}

func pullAll() error {
	resp, err := api.Get(context.Background(), MakeKey("share", env), &client.GetOptions{Recursive: true})
	if err != nil {
		return err
	}
	for _, node := range resp.Node.Nodes {
		entry[strings.Replace(node.Key, "/share/"+env+"/", "", -1)] = node.Value
	}
	resp, err = Get("", &client.GetOptions{Recursive: true})
	if err != nil {
		return err
	}
	for _, node := range resp.Node.Nodes {
		entry[strings.Replace(node.Key, prePath, "", -1)] = node.Value
	}
	return nil
}

func readFromDump() error {
	filename := IniConf.GetString("dunpPath", "confs.dump")
	if fileIsExist(filename) {
		entry = libconfig.NewIniConfig(filename).Entry
	} else {
		log.Println(filename, "is not exist!")
	}
	return nil
}

func fileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func getLocalInfo() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
		log.Println(err)
	}
	hostname = IniConf.GetString("instanceName", "instance") + "-" + hostname
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return hostname
	}
	for _, addr := range addrs {
		ok, err := regexp.Match(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`, []byte(addr.String()))
		if err == nil && ok && !strings.HasPrefix(addr.String(), "127.0.0.1") {
			ip := strings.Split(addr.String(), "/")[0]
			return fmt.Sprint(hostname, "-", ip)
		}
	}
	return hostname
}

func heart(info string) {
	for {
		select {
		case <-time.Tick(time.Second * 5):
			_, err := api.Set(context.Background(), MakeKey("heartbeat", prjName, env, info), "1", &client.SetOptions{TTL: time.Second * 5})
			if err != nil {
				log.Println(err)
				time.Sleep(time.Microsecond * 10)
			}
		}
	}
}

// GetLM 本地内存获取获取
func GetLM(key string) string {
	if v, ok := entry[key]; ok {
		return v
	}
	return ""
}

// MakeKey 生成key
func MakeKey(levels ...string) string {
	key := ""
	for _, l := range levels {
		key += "/" + l
	}
	return key
}
