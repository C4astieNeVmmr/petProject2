package main
import "fmt"
import "net/http"
import "encoding/json"
import "io"
import "os"

func main() {
    var urlTemplate = "https://api.telegram.org/bot%s/%s"
    var urlReceiveTemplate = fmt.Sprintf(urlTemplate+"?timeout=%s&offset=%s&limit=1",BOT_TOKEN,"getUpdates",TIMEOUT,"%d")
    var urlSendTemplate = fmt.Sprintf(urlTemplate+"?chat_id=%s&text=%s",BOT_TOKEN,"sendMessage","%d","%s")
    var offset = 0
    var data []byte
    var responseRaw *http.Response
    var err error
    data,err = os.ReadFile("allowed.json")
    if err!=nil {
        fmt.Println(err)
        return
    }
    var allowedUsers map[string]string
    err = json.Unmarshal(data,&allowedUsers)
    for ;true; {
        responseRaw, err = http.Get(fmt.Sprintf(urlReceiveTemplate,offset))
        if err != nil {
            fmt.Println(err)
            continue
        }
        data,err = io.ReadAll(responseRaw.Body)
        if err != nil {
            fmt.Println(err)
            continue
        }
        responseRaw.Body.Close()
        var dataJson Response
        err = json.Unmarshal(data,&dataJson)
        if err != nil {
            fmt.Println(err)
            continue
        }
        if len(dataJson.Result) == 0{
	    continue
	}

        offset = dataJson.Result[0].UpdateID+1
        fmt.Println(dataJson.Result[0].Message.From.Username)
        fmt.Println(dataJson.Result[0].Message.Text,"\n")
        if _,ok := allowedUsers[dataJson.Result[0].Message.From.Username]; !ok {
            _, err = http.Get(fmt.Sprintf(urlSendTemplate,dataJson.Result[0].Message.Chat.ID,"you are not allowed to do this"))
            if err!=nil {
                fmt.Println(err)
            }
            continue
        }
    }
}
