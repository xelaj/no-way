package noway

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ololosha228/keystorage"
	"github.com/xelaj/ifttt"
)

var (
	AppName string
	Token   string
)

func ALERT() {
	r := recover()
	if r == nil {
		return
	}

	recoverValue := fmt.Sprint(r)
	stackBytes := debug.Stack()

	stack := "Environment:\n" + strings.Join(os.Environ(), "\n") + "\n\nCallstack:\n" + string(stackBytes)
	stack = strings.ReplaceAll(stack, "\n", "</br>\n")

	if AppName == "" || Token == "" {
		rip()
		os.Stderr.Write([]byte("PANIC! " + recoverValue + "\n\nreport DOES NOT sent to ifttt\n"))
		os.Stderr.Write(stackBytes)
		os.Exit(2)
	}

	storage := keystorage.NewPrimitive("ifttt").Set(AppName, Token)
	iftttClient, _ := ifttt.NewClient(ifttt.ClientConfig{
		KeyStorage: storage,
	})

	err := iftttClient.By(AppName).Trigger(AppName+"_crashed", recoverValue, stack)
	if err != nil {
		rip()
		os.Stderr.Write([]byte("PANIC! " + recoverValue + "\n\nreport DOES NOT sent to ifttt\n"))
		os.Stderr.Write(stackBytes)
	} else {
		os.Stderr.Write([]byte("PANIC! " + recoverValue + "\n\nreport sent to ifttt\n"))
		os.Stderr.Write(stackBytes)
	}

	os.Exit(2)
}

func rip() {
	err := ioutil.WriteFile("/dev/null", []byte("https://bit.ly/3gmfkml"), 0660)
	if err != nil {
		println("https://youtu.be/2ift8T7eVEo")
	}
}
