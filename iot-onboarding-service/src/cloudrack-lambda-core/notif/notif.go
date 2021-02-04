package notif

import (
	"strings"
	db "cloudrack-lambda-core/db"
	core "cloudrack-lambda-core/core"
	"log"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"net/http"
	"net/http/httputil"
	"time"
	"os"
	"encoding/json"

)

//this needs to be moved in env variable
var WEB_SOCKET_ENDPOINT = "https://a4dqpl7ml7.execute-api.us-east-1.amazonaws.com/test/%40connections/"

func Notify(dbc db.DBConfig, user string, notif core.Notification) error {
	conn := core.DynamoUserConnection{}
	err := dbc.Get(user,"",&conn)
	if(err == nil) {
		log.Printf("[NOTIF] Got conection info for user %v => %+v",user,conn)
		//escapedId := strings.Replace(conn.ConnetionId,"=","%3D",-1)
		escapedId := conn.ConnetionId
		rqBody, _  := json.Marshal(notif)
		creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), os.Getenv("AWS_SESSION_TOKEN"))
		signer := v4.NewSigner(creds)
		bodyPayload := strings.NewReader(string(rqBody))
		req, _ := http.NewRequest("POST", WEB_SOCKET_ENDPOINT+escapedId, bodyPayload)
		req.Header.Set("Content-Type", "application/json")
		signer.Sign(req, bodyPayload, "execute-api", "us-east-1", time.Now())
		log.Printf("[NOTIF] signed request to POST:  %+v",req)
		//req.Body = ioutil.NopCloser(bodyPayload)
    	//res, err2 := sigv4.Post(WEB_SOCKET_ENDPOINT+escapedId , "application/json", data)
    	dump, _ := httputil.DumpRequest(req, true)
    	log.Printf("[NOTIF] REQUEST DUMP:  %s",dump)
    	res, err2 := http.DefaultClient.Do(req)
    	log.Printf("[NOTIF] response from web socket POST:  %+v",res)
    	return err2
	}
	return err
}