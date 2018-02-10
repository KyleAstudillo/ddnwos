package ddnwos

import "fmt"
import "net/http"
import "strings"
import "regexp"
import "bytes"
import "time"
import "log"
import "io/ioutil"
import "os"
import "io"
import "strconv"


//**************************** WosREST **************************************

type WosREST struct {
	keepalive bool
	ssl bool
	protocol string
	port string
	hosts_cycle []string
	index int
	client *http.Client
	recent *http.Response
	socket *http.Response
}

//*************** Alternative Overloading Functions *************************

func (self *WosREST) SimpleInit (ssl bool, hosts []string, port string) {
	self.index = 0
	self.keepalive = false
	self.ssl = ssl
	timeout := time.Duration(300 * time.Second)
	self.client = &http.Client{
		Timeout: timeout,
	}
	if(self.ssl == false){
		self.protocol = "http"
	}else{
		self.protocol = "https"
	}
	self.port = port
	self.hosts_cycle = hosts
}

func (self *WosREST) SimplePut(policy string, data string) string{
	return self.Put(policy, data, false, "", "", false)
}

func (self *WosREST) SimpleGet(oid string) string{
	return self.Get(oid, false, false, -99, -99, false, false, false, false)
}

func (self *WosREST) SimpleDelete(oid string){
	self.Delete(oid, false, false, false)
}

func (self * WosREST) SimpleExists(oid string) string{
	return self.Exists(oid, 204, false, false)
}

//************************ Helper functions ********************************

func (self *WosREST) getHost() string{
	nextIndex := self.index + 1
	if(nextIndex >= len(self.hosts_cycle)){
		self.index = 0
		nextIndex = 0
	}else{
		self.index = nextIndex
	}
	//ping host
	//if good host return
	//else mark as bad host
	//try another host
	//if all hosts are bad
	//    panic
	return self.hosts_cycle[nextIndex]
}

func (self *WosREST) getscheme() string{
	return fmt.Sprintf("%s://%s:%s", self.protocol, self.getHost(), self.port)
}

func (self *WosREST) process_status(status string){
        v := strings.Split(status, " ")
	if (len(v) > 0){
		if(v[0] != "200" && v[0] != "0"){
			panic("WosREST error: " + status)
		}
	}

}

//********************* Advanced Functions **************************************

func (self *WosREST) Put(policy string,
	data string,
	decmode bool,
	ddpversion string,
	metadata string,
	multipart bool) string{
      
        scheme := self.getscheme()
	req, err := http.NewRequest("POST", scheme + "/cmd/put", bytes.NewBufferString(data))
	if err != nil{
		panic(err)
	}
	req.Header.Set("Content-type", "application/octet-stream")
        req.Header.Set("x-ddn-policy", policy)
	if(decmode == true){
		req.Header.Set("x-ddn-distributed-protection", "true")
	}
	if(ddpversion != ""){
		req.Header.Set("x-ddn-force-ddp-version", ddpversion)
	}
	if(multipart == true){
		req.Header.Set("x-ddn-is-multipart", "true")
	}
        if(metadata != ""){
		re := regexp.MustCompile("{}")
		ra := regexp.MustCompile("'")
		amdata := re.ReplaceAllString(metadata, "")
		bmdata := ra.ReplaceAllString(amdata, "\"")
		req.Header.Set("x-ddn-meta", bmdata)
	}
	resp, err := self.client.Do(req)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()
	oid := resp.Header.Get("x-ddn-oid")
	status := resp.Header.Get("x-ddn-status")
	code := resp.StatusCode
	if(status != ""){
		status = fmt.Sprintf("%d %d%s", code, code, "_http_error")
	}
	self.recent = resp
	self.process_status(status)
	//if (self.keepalive == true){
	//	self.socket = resp
	//}
	return oid
}

func (self * WosREST) Get(oid string,
	buffered bool,
	integrity_check bool,
	rangeStart int,
	rangeEnd int,
	decmode bool,
	head bool,
	noddp bool,
	index_only bool) string{
        scheme := self.getscheme()
	req, err := http.NewRequest("GET", scheme + "/cmd/get", nil)
	if err != nil{
		panic(err)
	}
	req.Header.Add("x-ddn-oid", oid)
        req.Header.Add("x-ddn-buffered", "false") //change later
	req.Header.Add("x-ddn-integrity-check", "false") // change later		

	if (decmode){
		req.Header.Add("x-ddn-distributed-protection", "true")
	}
	if (rangeStart > 0){
		if(rangeEnd >= rangeStart){
			req.Header.Add("range", fmt.Sprintf("bytes= %d-%d", rangeStart, rangeEnd))
		}
	}
	if (noddp){
		req.Header.Add("x-ddn-force-no-goa", "true")
	}
	if(index_only){
		req.Header.Add("x-ddn-index-only", "true")
	}
        //if self.keepalive:
        //   self.socket = h
	resp, err := self.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	return string(body) //and metadata
}

func (self * WosREST) Delete(oid string,
	retry_deletes bool,
	background_deletes bool,
	decmode bool){
	scheme := self.getscheme()
	req, err := http.NewRequest("POST", scheme + "/cmd/delete", nil)
	if err != nil{
		panic(err)
	}
        req.Header.Add("x-ddn-oid", oid)
	if (retry_deletes){
		req.Header.Add("x-ddn-retry-delete", "true")
	}
	if (background_deletes){
		req.Header.Add("x-ddn-background-delete", "true")
	}
	if (decmode){
		req.Header.Add("x-ddn-distributed-protection", "true")
	}
        //if self.keepalive:
        //   self.socket = h
	resp, err := self.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	status := resp.Header.Get("x-ddn-status")
	self.process_status(status)
}

func (self * WosREST) Exists(oid string,
	expected_code int,
	head bool,
	decmode bool) string{
	scheme := self.getscheme()
	req, err := http.NewRequest("GET", scheme + "/cmd/exists", nil)
	if err != nil{
		panic(err)
	}
        req.Header.Add("x-ddn-oid", oid)
	if (decmode){
		req.Header.Add("x-ddn-distributed-protection", "true")
	}
	//if (self.keepalive){
	//	self.socket = h
	//}
	resp, err := self.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	status := resp.Header.Get("x-ddn-status")
	code := resp.StatusCode
	if (code != expected_code){
		panic("WOS error: " + string(code))
	}
	return status//, h.info().getheaders()
}

func (self * WosREST) CreatePutStream(policy string, datalen int64, metadata string) *PutStream {
	//if self.keepalive:
        //   self.socket = h
	putstream := PutStream{
		parent: self,
		policy: policy,
		datalen: datalen,
		metadata: metadata,
	}
	return &putstream
}

func (self * WosREST) CreateGetStream(oid string) *GetStream{
	getstream := GetStream{
		parent: self,
		oid: oid,
	}
	return &getstream
}


//*************************** PUTSTREAM *****************************
type PutStream struct {
	parent *WosREST
	handle string
	policy string
	datalen int64
        metadata string
}

func (self *PutStream) init (parent *WosREST, datalen int64) {
	self.parent = parent
	self.datalen = datalen
}

func (self *PutStream) Putter (req *http.Request) string {
	req.Header.Set("Content-type", "application/octet-stream")
	req.Header.Set("x-ddn-policy", self.policy)
	if (self.metadata != ""){
		//if (type(metadata) == dict){
		//	metadata = str(metadata)[1:-1]
		//}
		//metadata = metadata.replace('\'', '\"')
		req.Header.Add("x-ddn-meta", self.metadata)
	}
	resp, err := self.parent.client.Do(req)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()
	oid := resp.Header.Get("x-ddn-oid")
	return oid
}

func (self *PutStream) Put (data io.Reader) string {
	scheme := self.parent.getscheme()
	req, err := http.NewRequest("POST", scheme + "/cmd/put", data)
	req.ContentLength = self.datalen
	if err != nil{
		panic(err)
	}
	return self.Putter(req)
}

func (self *PutStream) PutString (data string) string {
	scheme := self.parent.getscheme()
	req, err := http.NewRequest("POST", scheme + "/cmd/put", bytes.NewBufferString(data))
	if err != nil{
		panic(err)
	}
	return self.Putter(req)
}

func (self *PutStream) PutFile (data *os.File) string{
	scheme := self.parent.getscheme()
	req, err := http.NewRequest("POST", scheme + "/cmd/put", data)
	req.ContentLength = self.datalen
	if err != nil{
		panic(err)
	}
	return self.Putter(req)
}

func (self *PutStream) Close () string {
      return "TODO"
}

//*************************** GETSTREAM ****************************
type GetStream struct {
	parent *WosREST
	handle string
	oid string
	//dont set
	resp *http.Response
}

func (self *GetStream) init (parent *WosREST, oid string) {
	self.parent = parent
	self.oid = oid
}

func (self *GetStream) Getter (req *http.Request) io.ReadCloser {
	req.Header.Add("x-ddn-oid", self.oid)
        req.Header.Add("x-ddn-buffered", "false") //change later
	req.Header.Add("x-ddn-integrity-check", "false") // change later		

        //if self.keepalive:
        //   self.socket = h

	resp, err := self.parent.client.Do(req)
	if err != nil {
		panic(err)
	}
	self.resp = resp
	return resp.Body
}

func (self *GetStream) Read () string {
	scheme := self.parent.getscheme()
	req, err := http.NewRequest("GET", scheme + "/cmd/get", nil)
	if err != nil{
		panic(err)
	}
	respBody := self.Getter(req)
	defer self.resp.Body.Close()
	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	return string(body)
}

func (self *GetStream) ReadRange (start int, end int) string {
	scheme := self.parent.getscheme()
	req, err := http.NewRequest("GET", scheme + "/cmd/get", nil)
	if (start > 0){
		if(end >= start){
			req.Header.Add("range", fmt.Sprintf("bytes= %d-%d", start, end))
		}
	}
	if err != nil{
		panic(err)
	}
	respBody := self.Getter(req)
	defer self.resp.Body.Close()
	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	return string(body)
}

func (self *GetStream) GetReader () io.ReadCloser {
	scheme := self.parent.getscheme()
	req, err := http.NewRequest("GET", scheme + "/cmd/get", nil)
	if err != nil{
		panic(err)
	}
	return self.Getter(req)
}

func (self *GetStream) GetStatus () string {
	return self.resp.Header.Get("x-ddn-status")
}

func (self *GetStream) ReadToFile (filename string, perm os.FileMode) {
	scheme := self.parent.getscheme()
	req, err := http.NewRequest("GET", scheme + "/cmd/get", nil)
	if err != nil{
		panic(err)
	}
	respBody := self.Getter(req)
	defer self.resp.Body.Close()
	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, body, perm)
	if err != nil{
		panic(err)
	}
}

func (self *GetStream) GetLength () int {
	i, err := strconv.Atoi(self.resp.Header.Get("Content-Length"))
	if err != nil{
		panic(err)
	}
	return i
}

func (self *GetStream) Close () {
	self.resp.Body.Close()
}
