package main
import "fmt"
import "net/http"
import "encoding/json"
import "io"

func main() {
    var urlTemplate,method = "https://api.telegram.org/bot%s/%s?timeout=%s&offset=%s&limit=1","getUpdates"
    var offset = 0
    var url string
    var data []byte
    var responseRaw *http.Response
    var err error
    urlTemplate = fmt.Sprintf(urlTemplate,BOT_TOKEN,method,TIMEOUT,"%d")
    for ;true; {
        url = fmt.Sprintf(urlTemplate,offset)
        responseRaw, err = http.Get(url)
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
    }
}
