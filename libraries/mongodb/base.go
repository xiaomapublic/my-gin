//mongodb服务类,使用mgo库
package mongodb

import (
	"fmt"
	"github.com/rs/xid"
	"gopkg.in/mgo.v2"
	. "my-gin/libraries/config"
	"my-gin/libraries/log"
	"time"
)

var MongoSession map[string]*mgo.Database

func init() {

	MongoSession = make(map[string]*mgo.Database, len(UnmarshalConfig.Mongodb))

	for key, c := range UnmarshalConfig.Mongodb {
		addrArr := c.Addr

		addrs := make([]string, len(addrArr))
		for _, add := range addrArr {
			addrs = append(addrs, add)
		}

		user := c.User
		pwd := c.Pwd
		max_active := c.Max_active
		MongoSession[key] = NewMongodb(addrs, key, user, pwd, max_active)
	}
}

func NewMongodb(addr []string, databaseName, user, pwd string, max_active int) *mgo.Database {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     addr,
		Timeout:   5 * time.Second,
		Database:  databaseName,
		Username:  user,
		Password:  pwd,
		PoolLimit: max_active,
	})

	if err != nil {
		log.InitLog("mongodb").Errorf("NewMongodb", "msg", err.Error())
		fmt.Println(err)
	}

	//Strong 一致性模式
	//session 的读写操作总向 primary 服务器发起并使用一个唯一的连接，因此所有的读写操作完全的一致（不存在乱序或者获取到旧数据的问题）。
	//Monotonic 一致性模式
	//session 的读操作开始是向某个 secondary 服务器发起（且通过一个唯一的连接），只要出现了一次写操作，session 的连接就会切换至 primary 服务器。由此可见此模式下，能够分散一些读操作到 secondary 服务器，但是读操作不一定能够获得最新的数据。
	//Eventual 一致性模式
	//session 的读操作会向任意的 secondary 服务器发起，多次读操作并不一定使用相同的连接，也就是读操作不一定有序。session 的写操作总是向 primary 服务器发起，但是可能使用不同的连接，也就是写操作也不一定有序。Eventual 一致性模式最快，其是一种资源友好（resource-friendly）的模式。
	//但是strong模式和Monotonic模式会缓存socket到session中，导致我拿到的始终是同一个连接，这对并发请求mongo server会损失一定的效率（一个连接在使用过程中会有多次加锁解锁操作）。所以后续为了从连接池拿空闲的连接而不是一直使用同一个连接，会用到copy方法，拿到一个没有缓存连接的session，这样它就会去连接池拿空闲的可用的连接
	session.SetMode(mgo.Eventual, true)
	//通常使用session都需要Clone或Copy然后再Close关闭
	//如果不使用则需要使用Eventual模式，该模式最大的特点就是不会缓存连接，拒绝持有mongodb的一针一线。通过该mgo.Database实例的操作，每次都会发现自己的口袋里一无所有，都会经过一次上一节所述的长途跋涉获取连接，因此也规避了不使用Copy带来的两个副作用。一并发效率问题，Eventual的Session每次操作都从连接池取连接，相当于分散在连接池中完成了操作，二连接可用性问题，连接池机制确保了，从mongoServer取得的连接，都是活的连接。
	return session.DB(databaseName)
}

//创建_id
func CreateId() string {
	return xid.New().String()
}
