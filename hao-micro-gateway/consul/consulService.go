package consul

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

type HaoServiceRegistration struct {
	Id          string            `json:"id"`          //注册服务id
	Name        string            `json:"name"`        //注册服务名称
	Name_mc     string            `json:"name_mc"`     //注册服务名称中文
	Meta        map[string]string `json:"meta"`        //注册服务名称中文
	Address     string            `json:"address"`     //服务注册ip
	Port        int               `json:"Port"`        //服务注册端口
	Connections int               `json:"connections"` //链接数
	Check       HaoServiceCheck   `json:"check"`       //服务注册核查信息

}

type HaoServiceCheck struct {
	Node        string              `json:"node"`        //
	CheckID     string              `json:"checkID"`     //注册服务核查id
	CheckType   string              `json:"checkType"`   //核查  TCP/HTTP（get 请求"https://localhost:5000/health"）
	Name        string              `json:"name"`        //注册服务核查名称
	Name_mc     string              `json:"name_mc"`     //注册服务核查名称中文
	Interval    string              `json:"interval"`    //核查多长时间检查一次 单位s
	Timeout     string              `json:"timeout"`     //核查超时时间 单位s
	Tcp         string              `json:"tcp"`         //tcp 核查地址
	Http        string              `json:"http"`        //http 核查地址
	Method      string              `json:"method"`      //http 核查地址 类型（get/post）
	Header      map[string][]string `json:"header"`      //http 请求头参数
	Body        string              `json:"body"`        //http 请求报文
	Status      string              `json:"status"`      //检查状态
	PassingOnly bool                `json:"passingOnly"` //是否检查异常节点；true-只检查正常
	Output      string              `json:"output"`      //核查信息
}

// 创建 consul 客户端
var CONSUL_CLIENT *api.Client

// 定义一个全局的服务实例列表和锁
var LeastConnectionBalancer = make(map[string][]HaoServiceRegistration)
var lock = &sync.Mutex{}

// 初始化客户端
func IntoConsulClient(address string, timeTicker int) error {
	client, err := api.NewClient(&api.Config{
		// Address: "101.43.7.108:8500", // 替换为你的Consul服务器地址和端口
		Address: address,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("初始化conusl客户端失败，地址： %s , 错误：%s", address, err.Error()))
	}
	CONSUL_CLIENT = client

	//开启服务定时任务
	SyncTickerConsulData(timeTicker)
	return nil
}

// 定义接口
type ConsulService interface {
	ConsulServiceRegister() error                          //注册服务
	ConsulServiceDeregister() error                        // 注销服务
	ConsulServiceQuery() ([]HaoServiceRegistration, error) // 获取服务列表
}

// 注册服务
func (hst HaoServiceRegistration) ConsulServiceRegister() error {
	// 创建一个服务实例
	reg := &api.AgentServiceRegistration{
		// ID:   uuid.New().String(),
		ID:   hst.Id,
		Name: hst.Name,
		// Tags:    []string{"golang"},
		Address: hst.Address,
		Port:    hst.Port,
		Meta:    hst.Meta,
		Check: &api.AgentServiceCheck{
			CheckID: hst.Check.CheckID,
			// TCP:      "101.43.7.108:8002",
			Name:     hst.Check.Name,
			Interval: hst.Check.Interval, //10s
			Timeout:  hst.Check.Timeout,  //5s
		},
	}

	if hst.Check.CheckType == "TCP" {
		reg.Check.TCP = hst.Check.Tcp
	} else if hst.Check.CheckType == "HTTP" {
		if hst.Check.Method == "GET" {
			reg.Check.HTTP = hst.Check.Http
		} else if hst.Check.Method == "POST" {
			reg.Check.HTTP = hst.Check.Http
			reg.Check.Method = hst.Check.Method
			reg.Check.Header = hst.Check.Header
			reg.Check.Body = hst.Check.Body
		}
	} else {
		reg.Check.TCP = fmt.Sprintf("%s:%d", hst.Address, hst.Port)
	}

	// 注册服务
	err := CONSUL_CLIENT.Agent().ServiceRegister(reg)
	if err != nil {
		return err
	}
	fmt.Println("Service registered to Consul")
	return nil
}

// 注销服务
func (hst HaoServiceRegistration) ConsulServiceDeregister() error {
	err := CONSUL_CLIENT.Agent().ServiceDeregister(hst.Id)
	if err != nil {
		return err
	}
	fmt.Println("Service deregistered from Consul")
	return nil
}

func (hst HaoServiceRegistration) ConsulServiceQuery() ([]HaoServiceRegistration, error) {
	// 获取服务的健康实例列表
	var hsts []HaoServiceRegistration
	instances, _, err := CONSUL_CLIENT.Health().Service(hst.Name, "", hst.Check.PassingOnly, &api.QueryOptions{})
	if err != nil {
		return nil, errors.New("获取服务列表失败," + err.Error())
	}

	// 遍历服务实例列表
	for _, instance := range instances {
		hst := HaoServiceRegistration{
			Id:      instance.Service.ID,
			Name:    hst.Name,
			Address: instance.Service.Address,
			Port:    instance.Service.Port,
			Meta:    instance.Service.Meta,
		}
		for _, check := range instance.Checks {
			if check.ServiceName == hst.Name {
				// fmt.Print(fmt.Sprintf("检查信息 %+v", check))
				checkinfo := HaoServiceCheck{
					CheckID: check.CheckID,
					Name:    check.Name,
					Node:    check.Node,
					Status:  check.Status,
					Output:  check.Output,
				}
				hst.Check = checkinfo
			}
		}
		hsts = append(hsts, hst)
		// fmt.Printf("Service: %s, Address: %s, Port: %d", serviceName, address, port)
	}
	return hsts, nil
}

// 获取服务列表  prefix 需要过滤下前缀
func ConsulServiceList(prefix string, passingOnly bool) (map[string][]HaoServiceRegistration, error) {
	services, _, err := CONSUL_CLIENT.Catalog().Services(&api.QueryOptions{})
	if err != nil {
		return nil, errors.New("获取服务信息失败," + err.Error())
	}
	list := make(map[string][]HaoServiceRegistration)
	// 遍历服务列表
	for serviceName := range services {
		// 过滤指定名称的服务
		var instances []*api.ServiceEntry
		var err error
		if prefix != "" && len(prefix) > 0 {
			if strings.HasPrefix(serviceName, prefix) {
				// 获取服务的健康实例列表
				instances, _, err = CONSUL_CLIENT.Health().Service(serviceName, "", passingOnly, &api.QueryOptions{})
			} else {
				err = errors.New("不符合数据的服务！")
			}
		} else {
			// 获取服务的健康实例列表
			instances, _, err = CONSUL_CLIENT.Health().Service(serviceName, "", passingOnly, &api.QueryOptions{})
		}

		if err == nil {
			// 遍历服务实例列表
			var hsts []HaoServiceRegistration
			// 遍历服务实例列表
			for _, instance := range instances {
				hst := HaoServiceRegistration{
					Id:      instance.Service.ID,
					Name:    serviceName,
					Address: instance.Service.Address,
					Port:    instance.Service.Port,
					Meta:    instance.Service.Meta,
				}
				for _, check := range instance.Checks {
					if check.ServiceName == hst.Name {
						// fmt.Print(fmt.Sprintf("检查信息 %+v", check))
						checkinfo := HaoServiceCheck{
							CheckID: check.CheckID,
							Name:    check.Name,
							Node:    check.Node,
							Status:  check.Status,
							Output:  check.Output,
						}
						hst.Check = checkinfo
					}
				}
				hsts = append(hsts, hst)
				// fmt.Printf("Service: %s, Address: %s, Port: %d", serviceName, address, port)
			}
			list[serviceName] = hsts
		}
	}
	return list, nil
}

func SyncTickerConsulData(timeTicker int) {
	// 创建一个定时器，每隔 timeTicker 秒执行一次任务
	timer := time.NewTicker(time.Duration(timeTicker) * time.Second)

	// 在一个新的goroutine中执行定时任务
	go func() {
		for range timer.C {
			// 执行你的定时任务代码
			// dateTime := utils.TimeFormatNow(utils.Format_YMDHS)
			// fmt.Println(fmt.Sprintf("定时每隔[%d]秒任务开始，当前时间[%s]", timeTicker, dateTime))
			list, _ := ConsulServiceList("hao_", true)
			SetLeastConnectionInstance(list)
		}
	}()
	// // 等待10秒后停止定时器
	// time.Sleep(10 * time.Second)
	// timer.Stop()
	// fmt.Println("定时任务结束")
}

func SetLeastConnectionInstance(list map[string][]HaoServiceRegistration) {
	lock.Lock()
	defer lock.Unlock()
	for serviceName, _ := range LeastConnectionBalancer {
		_, ok := list[serviceName]
		if !ok {
			delete(LeastConnectionBalancer, serviceName)
		}
	}
	for serviceName, check_hst := range list {
		//获取服务列表
		instances, ok := LeastConnectionBalancer[serviceName]
		var hsts []HaoServiceRegistration
		if ok && len(instances) > 0 {
			//先遍历注册服务上列表
			for _, check_hst_data := range check_hst {
				id := check_hst_data.Id
				check_hst_data.Connections = 0
				//获取本地服务信息
				for _, instance := range instances {
					if id == instance.Id {
						check_hst_data.Connections = instance.Connections
					}
				}
				hsts = append(hsts, check_hst_data)
			}
		} else {
			//todo 如果不存在，则直接给集合赋值
			hsts = check_hst
		}
		LeastConnectionBalancer[serviceName] = hsts
	}
}

// 根据服务名称获取最小连接数的服务实例
func GetLeastConnectionInstance(serviceName string) (HaoServiceRegistration, error) {
	lock.Lock()
	defer lock.Unlock()

	instances, ok := LeastConnectionBalancer[serviceName]
	if !ok || len(instances) == 0 {
		return HaoServiceRegistration{}, errors.New(fmt.Sprintf("获取服务[%s]失败！", serviceName))
	}

	if len(instances) == 1 {
		return instances[0], nil
	}

	// 找到连接数最少的服务实例
	minConnections := instances[0].Connections
	minInstance := instances[0]
	for _, instance := range instances {
		if instance.Connections < minConnections {
			minConnections = instance.Connections
			minInstance = instance
		}
	}
	// 更新连接数
	minInstance.Connections++
	var hsts []HaoServiceRegistration
	for _, instance := range instances {
		if instance.Id == minInstance.Id {
			instance.Connections = minInstance.Connections
		}
		hsts = append(hsts, instance)
	}
	LeastConnectionBalancer[serviceName] = hsts

	// 返回最小连接数的服务实例
	return minInstance, nil
}
