# ddnwos #
## Description ##
DDN WOS Rest Helper class to facilitate use of the HTTP interface to WOS
## Set Up ##
Clone repo into $GOROOT/src or /usr/local/go/src
## All Functions ##
* type WosREST
  * func (self *WosREST) SimpleInit (ssl bool, hosts []string, port string)
  * func (self *WosREST) SimplePut(policy string, data string) string
  * func (self *WosREST) SimpleGet(oid string) string
  * func (self *WosREST) SimpleDelete(oid string)
  * func (self * WosREST) SimpleExists(oid string) string
  * func (self * WosREST) Close()
  * func (self *WosREST) Put(policy string, data string, decmode bool, ddpversion string, metadata string, multipart bool) string
  * func (self * WosREST) Get(oid string, buffered bool, integrity_check bool, rangeStart int, rangeEnd int, decmode bool, head bool, noddp bool, index_only bool) string
  * func (self * WosREST) Delete(oid string, retry_deletes bool, background_deletes bool, decmode bool)
  * func (self * WosREST) Exists(oid string, expected_code int, head bool, decmode bool) string
  * func (self * WosREST) CreateGetStream(oid string, buffered bool, integritycheck bool) *GetStream
  * func (self * WosREST) CreateGetStream(oid string, buffered bool, integritycheck bool) *GetStream
* type PutStreams
  * func (self *PutStream) Put (data io.Reader) string
  * func (self *PutStream) PutString (data string) string
  * func (self *PutStream) PutFile (data *os.File) string
  * func (self *PutStream) Close ()
* typeGetStreams
  * func (self *GetStream) Read () string
  * func (self *GetStream) ReadRange (start int, end int) string
  * func (self *GetStream) GetReader () io.ReadCloser
  * func (self *GetStream) ReadToFile (filename string, perm os.FileMode)
  * func (self *GetStream) GetStatus () string
  * func (self *GetStream) GetLength () int
  * func (self *GetStream) Close ()

## Example ##
#### Initialize WosREST ####
Description: SimpleInit(ssl bool, []string{IPs}, port string)
use to quickly create a WosREST obj with defaults
``` go
wos := ddnwos.WosREST{}
wos.SimpleInit(false, []string"1.2.3.4", "5.6.7.8", "9.10.11.12"}, "80")
```
Can also create a WosREST obj that is fully customizable
``` go
wos := WosREST{  keepalive: false,
                 ssl: false,
                 protocol: "http",
                 port: "80",
                 hosts_cycle: []string{ "1.2.3.4",
                                        "5.6.7.8",
                                        "9.10.11.12"
                 },
                 index: 0
                 client: &http.Client{}
                 badConn: make(map[string]bool)
                 buffered: false
                 integritycheck: false
                 debugtoggle: false
}
```
#### Simple Functions ####
Description: SimplePut is a function that can store a string of data to a WOS cluster policy then returns an oid string that can be used to retrieve the string from WOS, for more complex puts check out PutStreams.
``` go       
oid := wos.SimplePut("POLICYNAME", "DATA STRING")
wos.SimpleGet(oid)
wos.SimpleDelete(oid)
wos.SimpleExists(oid)
```
#### Put Stream ####
Description: PutStream object can be used with more complex put operations such as put files, or io.readers. **Do not forget to use Close function if keepalive is not set to false.**
``` go 
putstream := wos.CreatePutStream("POLICYNAME", INT, "METADATA")  //METADATA can be an empty string
oid = putstream.PutString("DATA STRING")
putstream.Close()
```
##### PutFile Function #####
Description: PutFile is a function that can store a file to a WOS cluster then returns an oid string that can be used to retrieve the file from WOS, **Do not forget to use Close function if keepalive is not set to false.**
``` go 
file, err := os.Open("FILENAME.txt")
if err != nil {
    log.Fatal(err)
}
oid = putstream.PutFile(file)
putstream.Close()
```
#### Get Stream ####
Description: GetStream object can be used with more complex get operations such as readRange of object, read object to file, or obtain io.CloseReader 
##### Read Function #####
Description: Read function simply gets oid and returns the object data in a string format, Close automatically called with read function
``` go 
getstream := wos.CreateGetStream(oid)
data := getstream.Read()
getstream.Close()  //not needed on GetStream.Read Function 
```

##### ReadRange Function #####
Description: ReadRange function reads range of bytes specified from object in WOS, Close automatically called with read function
``` go 
getstream := wos.CreateGetStream(oid)
data := getstream.ReadRange(START_INT, END_INT)
getstream.Close()  //not needed on GetStream.Read Function 
```

##### get io.CloseReader Function #####
Description: GetReader function return io.CloseReader, **Make sure to use Close operation when finished reading.**
``` go 
reader := getstream.GetReader()
body, err := ioutil.ReadAll(reader)
if err != nil {
    log.Fatalf("ERROR: %s", err)
}
println(string(body))
getstream.Close()  //This is needed make sure to close if this Function is used
```

##### ReadToFile Function #####
Description: ReadToFile function reads object from WOS and stores it locally in a file specified with permission specified.
``` go 
getstream.ReadToFile("GetStreamMadeThisFile.txt", 0644)
```

 ##### Long Put Get Delete Functions #####
 Description: Put, Get, Delete functions filled out with defaults see All Functions at the top to see which parameters map to what
 ``` go 
oid = wos.Put("POLICY NAME", "Data String", false, "", "",false)
wos.Get(oid, false, false, -99, -99 false, false, false, false)
wos.Delete(oid , false, false, false)
```

