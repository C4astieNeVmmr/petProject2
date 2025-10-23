package main
import "fmt"
import "os"
import "os/exec"
import "encoding/json"
import "errors"

type r_change struct{
    Links_dict map[string]string
    Files_dict map[string]string
    Current_resource string
}

var template_link = "location /%s {\n\treturn 302 %s;\n}\n"
var template_file = "location /%s {\n\troot %s;\n\tindex %s;\n}\n"

func (self *r_change) load(filename string) (error){
    var data,err = os.ReadFile(PATH_TO_RESOURCES+filename)
    if err!=nil {
        return err
    }
    err = json.Unmarshal(data,self)
    if err!=nil {
        return err
    }
    return nil
}

func (self *r_change) dump(filename string) (error){
    var data,err = json.MarshalIndent(self,"","\t")
    if err!=nil {
        return err
    }
    err = os.WriteFile(PATH_TO_RESOURCES+filename,data,0644)
    if err!=nil {
        return err
    }
    return nil
}

func (self *r_change) change(name string) (error){
    if self.Current_resource == name {
        return nil
    }
    self.Current_resource = name
    var s string
    if _, ok := self.Links_dict[name]; ok {
        s = fmt.Sprintf(template_link,SITE_LOCATION_NAME,self.Links_dict[name])
    }else if _, ok := self.Files_dict[name]; ok {
        s = fmt.Sprintf(template_file,SITE_LOCATION_NAME,PATH_TO_PAGES,self.Files_dict[name])
    } else {
        return errors.New("no such link or file")
    }
    var data = []byte(s)
    os.WriteFile(PATH_TO_NGINX+SITE_CONFIG_NAME,data,0644)
    var cmd = exec.Command("nginx", "-s", "reload")
    var err = cmd.Run()
    if err!=nil {
        return err
    }
    return nil
}

type command struct{
    Number_of_args int
    Description string
}
