package main

import (
  "flag"
  "fmt"
  "github.com/syndtr/goleveldb/leveldb"
  "github.com/syndtr/goleveldb/leveldb/util"
  "time"
)

const (
  chequebookKey           = "swap_chequebook"
  ChequebookDeploymentKey = "swap_chequebook_transaction_deployment"
  overlayKey = "overlay"

  balanceCheckBackoffDuration = 20 * time.Second
  balanceCheckMaxRetries      = 10

  dbSchemaKey = "statestore_schema"

  dbSchemaGrace = "grace"


  lastIssuedChequeKeyPrefix = "swap_chequebook_last_issued_cheque_"
  totalIssuedKey            = "swap_chequebook_total_issued_"

)

// store uses LevelDB to store values.
type store struct {
  db *leveldb.DB
}



type DbMap = map[string] []byte


func main() {

  var action string
  var source string
  var target string
  var key string
  var value string
  flag.StringVar(&action,"action","","action:get,put,copy,delete")
  flag.StringVar(&source,"source","","源目录")
  flag.StringVar(&target,"target","","目标目录")
  flag.StringVar(&key,"key","","关键字")
  flag.StringVar(&value,"value","","值")
  flag.Parse()


  var s *store = nil
  var t *store = nil
  if  source != "" {
    s = getStore(source)
  }else{
    return
  }
  if  target != "" {
    t = getStore(target)
  }
  defer func() {
    if s != nil {
      s.db.Close()
    }
    if t != nil {
      t.db.Close()
    }

  }()

  switch action {
  case "get" :
    s.Get(key)
  case "put" :
    s.Put(key, value)
  case "delete" :
    s.delete(key)
  case "copy" :
    if t != nil {
      s.copy(key, t)
    }
  default:
    fmt.Println("未知操作")
  
  }

  //deleteData(lastIssuedChequeKeyPrefix)
  //deleteData(totalIssuedKey)
  //downData()
  //copyData(chequebookKey)
  //copyData(overlayKey)

  //getSchemaName()

}


func getStore(path string) *store {
  db, err := leveldb.OpenFile(path, nil)
  if err != nil {
    fmt.Println(err, "init!!!")
    return nil
  }
  stateStore := &store{db: db}
  return stateStore
}

func (s *store) Get(key string) DbMap {
  result := DbMap{}
  iter := s.db.NewIterator(util.BytesPrefix([]byte(key)), nil)
  for iter.Next() {
    key := iter.Key()
    value := iter.Value()

    fmt.Println( string(key) + ":" + string(value))
    result[string(key)] = value
  }
  iter.Release()
  return  result
}

func (s *store) delete(key string) {
  iter := s.db.NewIterator(util.BytesPrefix([]byte(key)), nil)

  for iter.Next() {
    key := iter.Key()
    value := iter.Value()
    fmt.Println( string(key) + ":" + string(value))
    s.db.Delete(key, nil)
  }
  iter.Release()
}

func (s *store) copy(key string, t *store) {
  iter := s.db.NewIterator(util.BytesPrefix([]byte(key)), nil)
  for iter.Next() {
    key := iter.Key()
    value := iter.Value()

    fmt.Println( string(key) + ":" + string(value))
    t.Put(key,value)
  }
  iter.Release()

}

func (s *store) Put(key,value interface{}) {
  s.db.Put(toBytes(key) ,toBytes(value),nil)
}

func toBytes(d interface{}) []byte {
  switch d.(type) {
    case string:
      str := d.(string)
      return []byte(str)
  case []byte :
      return d.([]byte)
  default:
    return nil
  }
}