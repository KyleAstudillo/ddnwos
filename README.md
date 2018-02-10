# ddnwos #
## Description ##
DDN WOS Rest Helper class to facilitate use of the HTTP interface to WOS
## Set Up ##
Clone repo into $GOROOT/src or /usr/local/go/src
## Example ##
#### Inialize WosREST ####
``` go
wos := ddnwos.WosREST{}
wos.SimpleInit(false, []string"1.2.3.4", "5.6.7.8", "9.10.11.12"}, "80")
```
or
``` go
wos := WosREST{  keepalive: false,
                 ssl: false,
                 protocol: "http",
                 port: "80",
                 hosts_cycle: []string{ "1.2.3.4",
                                        "5.6.7.8",
                                        "9.10.11.12"},
}
```
#### Simple Functions ####
``` go       
oid := wos.SimplePut("POLICYNAME", "DATA STRING")
wos.SimpleGet(oid)
wos.SimpleDelete(oid)
wos.SimpleExists(oid)
```
#### Put Stream ####
``` go 
putstream := wos.CreatePutStream("POLICYNAME", INT, "METADATA")  //METADATA can be an empty string
oid = putstream.PutString("DATA STRING")
```
##### PutFile Function #####
``` go 
file, err := os.Open("FILENAME.txt")
if err != nil {
    log.Fatal(err)
}
oid = putstream.PutFile(file)
```
#### Get Stream ####
##### Read Function #####
``` go 
getstream := wos.CreateGetStream(oid)
data := getstream.Read()
getstream.Close()  //not needed on GetStream.Read Function 
```

##### ReadRange Function #####
``` go 
getstream := wos.CreateGetStream(oid)
data := getstream.ReadRange(START_INT, END_INT)
getstream.Close()  //not needed on GetStream.Read Function 
```

##### get io.CloseReader Function #####
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
``` go 
getstream.ReadToFile("GetStreamMadeThisFile.txt", 0644)
```

 ##### Long Put Get Delete Functions #####
 ``` go 
oid = wos.Put("POLICY NAME", "Data String", false, "", "",false)
wos.Get(oid, false, false, -99, -99 false, false, false, false)
wos.Delete(oid , false, false, false)
```

