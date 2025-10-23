package main
import "fmt"
import "net/http"
import "encoding/json"
import "io"
import "os"
import "os/exec"
import "strings"
import "net/url"

func main() {
    var urlTemplate = "https://api.telegram.org/bot%s/%s"
    var urlReceiveTemplate = fmt.Sprintf(urlTemplate+"?timeout=%s&offset=%s&limit=1",BOT_TOKEN,"getUpdates",TIMEOUT,"%d")
    var urlSendTemplate = fmt.Sprintf(urlTemplate+"?chat_id=%s&text=%s",BOT_TOKEN,"sendMessage","%d","%s")
    var offset = 0
    var data []byte
    var responseRaw *http.Response
    var dataJson Response
    var reactionToCommand string
    var command_arr []string
    var err error
    var changeObj r_change
    err = changeObj.load(PATH_TO_RESOURCES+"resources.json")
    if err!=nil {
        fmt.Println(err)
        return
    }
    data,err = os.ReadFile(PATH_TO_ALLOWED_ACCOUNTS+"allowed.json")
    if err!=nil {
        fmt.Println(err)
        return
    }
    var allowedUsers map[string]string
    err = json.Unmarshal(data,&allowedUsers)
    if err!=nil {
        fmt.Println(err)
        return
    }
    fmt.Println("The program has started")
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
        command_arr = strings.Split(dataJson.Result[0].Message.Text," ")
        reactionToCommand = ""
        if _,ok := COMMANDS[command_arr[0]]; ok && (len(command_arr)!=COMMANDS[command_arr[0]].Number_of_args) {
            if len(command_arr)>COMMANDS[command_arr[0]].Number_of_args {
                reactionToCommand = "too much arguments for this command"
            } else if len(command_arr)<COMMANDS[command_arr[0]].Number_of_args{
                reactionToCommand = "too few arguments for this command"
            }
            _, err = http.Get(fmt.Sprintf(urlSendTemplate,dataJson.Result[0].Message.Chat.ID,reactionToCommand))
            if err!=nil {
                fmt.Println(err)
            }
            continue
        }
        switch(command_arr[0]){
            case "/start":
                reactionToCommand = "Hello! This bot are capable of changing/adding/removal of resources to which leadsspecial qr-code"
            case "/help":
                for key, value := range COMMANDS {
                    reactionToCommand += fmt.Sprintf("%s - args %d %s\n",key,value.Number_of_args-1,value.Description)
                }
            case "/change":
                err = changeObj.change(command_arr[1])
                if err!=nil {
                    fmt.Println(err)
                    continue
                }
                reactionToCommand = "resource changed"
            case "/add_link":
                changeObj.Links_dict[command_arr[1]] = command_arr[2]
                err = changeObj.dump(PATH_TO_RESOURCES+"resources.json")
                if err!=nil {
                    fmt.Println(err)
                    continue
                }
                reactionToCommand = "link added"
            case "/add_file":
                reactionToCommand = "/add_file not implemented yet"
            case "/remove":
                delete(changeObj.Links_dict,command_arr[1])
                delete(changeObj.Files_dict,command_arr[1])
                err = changeObj.dump(PATH_TO_RESOURCES+"resources.json")
                if err!=nil {
                    fmt.Println(err)
                    continue
                }
                reactionToCommand = "resource removed"
            case "/show_list":
                reactionToCommand = "links:\n"
                for key, value := range changeObj.Links_dict {
                    reactionToCommand += key+"\t-\t"+value+"\n"
                }
                reactionToCommand += "\nfiles:\n"
                for key, value := range changeObj.Files_dict {
                    reactionToCommand += key+"\t-\t"+value+"\n"
                }
            case "/reboot_nginx":
                var cmd = exec.Command("nginx", "-s", "reload")
                err = cmd.Run()
                if err!=nil {
                    fmt.Println(err)
                    continue
                }
                reactionToCommand = "nginx was rebooted"
            default:
                reactionToCommand = "there is no such command"
        }
        _, err = http.Get(fmt.Sprintf(urlSendTemplate,dataJson.Result[0].Message.Chat.ID,url.PathEscape(reactionToCommand)))
        if err!=nil {
            fmt.Println(err)
            continue
        }
    }
}
