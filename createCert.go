package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
)

type certData struct {
	key 	string
	cert 	string
}

//noinspection GoExportedFuncWithUnexportedType
func CreateCert(domain,t string,env map[string]string) certData {
	log.Println("begin to create a cert...")
	command :=""
	for k, v := range env {

		tmp := "export "+k+"="+v+";"
		command = command+tmp
	}
	command = command+`acme.sh --issue --dns `+t+` -d `+domain


	cmd := exec.Command("/bin/bash", "-c", command)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		log.Fatalln("fail to create a cert")
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		log.Println("Error:The command is err,", err)
		log.Fatalln("fail to create a cert")
	}

	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ReadAll Stdout:", err.Error())
		log.Fatalln("fail to create a cert")
	}

	if err := cmd.Wait(); err != nil {
		log.Println("wait:", err.Error())
		log.Fatalln("fail to create a cert")
	}
	output := string(bytes)
	reg, err := regexp.CompilePOSIX(`.+Your cert is in {2}(.+) \n.+Your cert key is in {2}(.+) `)
	if err != nil {
		log.Println("regexp error:", err.Error())
		log.Fatalln("fail to create a cert")
	}
	finds := reg.FindStringSubmatch(output)
	if len(finds) == 0 {
		log.Println("regexp find nothing:")
		log.Fatalln("fail to create a cert")
	}
    log.Println("finsh to create a cert...")
	return certData{
		key: 	base64encode(finds[2]),
		cert: 	base64encode(finds[1]),
	}
}


func base64encode(path string) string {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	//input := []byte(file)

	// base64编码
	encodeString := base64.StdEncoding.EncodeToString(file)
	return encodeString
}
