package initiator

import "github.com/abelmalu/golang-posts/platform"



func InitLogger()*platform.Logger{

    	return platform.InitZapLogger()

}