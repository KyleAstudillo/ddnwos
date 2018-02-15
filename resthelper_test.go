package ddnwos

import ("testing"
	"fmt"
	"os"
)

var tc RestHelperTestConfig = RestHelperTestConfig{} //TestConfig

type RestHelperTestConfig struct {
	wos *WosREST
	policy string
	data string
	oids map[int]string
	oidCount int
	ps *PutStream
	gs *GetStream
}

func (self *RestHelperTestConfig) IncOid(){
	self.oidCount = self.oidCount + 1
}

func (self *RestHelperTestConfig) DecOid(){
	self.oidCount = self.oidCount - 1
}

func TestMain(m *testing.M){
	policy := "replicate"//default
	data := "data" //default
	
	wos := WosREST{}
	wos.SimpleInit(false, []string{"1.2.3.4"}, "80") //change ip to cluster mts ip
	tc = RestHelperTestConfig{wos: &wos,
		policy: policy,
		data: data,
		oids: make(map[int]string),
		oidCount: 0,
	}
	os.Exit(m.Run())
}

//Test Alternative Overloading Functions

func TestSimplePut(t *testing.T){
	tc.oids[tc.oidCount] = tc.wos.SimplePut(tc.policy, fmt.Sprintf("%s %d",tc.data, tc.oidCount))
	fmt.Println(tc.wos.SimpleExists(tc.oids[tc.oidCount]))
	tc.IncOid()
	//Output: 0 ok
}

func TestSimpleGet(t *testing.T){
	data := tc.wos.SimpleGet(tc.oids[0])
	fmt.Println(data)
	//Output: data 0
}

func TestSimpleDelete(t *testing.T){
	tc.DecOid()
	oid := tc.oids[tc.oidCount]
	tc.wos.SimpleDelete(oid)
	fmt.Println(tc.wos.SimpleExists(oid))
	delete(tc.oids, tc.oidCount)
	//Output: 207 ObjNotFound
}

//TestPutStreams

func TestPutStreamPutString(t *testing.T){
	putstream := tc.wos.CreatePutStream(tc.policy, -1, "")
	tc.oids[tc.oidCount] = putstream.PutString(fmt.Sprintf("%s %d", tc.data, tc.oidCount))
	fmt.Println(tc.wos.SimpleExists(tc.oids[tc.oidCount]))
	tc.IncOid()
	//Output: 0 ok
}

//testGetStreams

func TestGetStreamRead(t *testing.T){
	getstream := tc.wos.CreateGetStream(tc.oids[0], false, false)
	data := getstream.Read()
	fmt.Println(data)
	//Output: data 0
}


