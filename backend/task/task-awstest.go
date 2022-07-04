package main

import (
	"os"
	"sync"
	"github.com/iost-official/explorer/backend/config"



	"fmt"
	"log"
	"strconv"
	"os/exec"



	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

)




var ws2 = new(sync.WaitGroup)


func main() {
	config.ReadConfig("")
	//avg := sigar.LoadAverage{}


	const interval =int64(1000)
	
	fromblock, _ := strconv.ParseInt(os.Args[1], 10, 64)
	toblock, _ := strconv.ParseInt(os.Args[2], 10, 64)




	inner_from_block :=fromblock
	inner_to_block := inner_from_block +interval -1

	
	for;inner_from_block<toblock;{
	
		
		
		from_to_name:=strconv.FormatInt(inner_from_block,10)+"_"+strconv.FormatInt(inner_to_block,10)
		exec.Command("mkdir","src/"+from_to_name).Run()


		ws2.Add(1)
		go func(from_to_name string) {
			
			
			// sessionの作成
			sess := session.Must(session.NewSessionWithOptions(session.Options{
				Config: aws.Config{Region: aws.String("ap-southeast-1")},
				Profile:           "default",
				SharedConfigState: session.SharedConfigEnable,
			}))

			// ファイルを開く
			targetFilePath := "src/"+from_to_name+".tar.gz"
			file, err := os.Open(targetFilePath)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			bucketName := "stg.stir-hosyu"
			objectKey := "iost-explorer/"+from_to_name+".tar.gz"

			// Uploaderを作成し、ローカルファイルをアップロード
			uploader := s3manager.NewUploader(sess)
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
				Body:   file,
			})
			if err != nil {
				log.Fatal(err)
			}
			log.Println("done") 
			ws2.Done()
		}(from_to_name)
		ws2.Wait()


		//S3へアップロード
		//SCPする

		inner_from_block += interval
		inner_to_block += interval 
	}
	
	// time.Sleep(10 * time.Second)
	// start tasks
	
	
	
}
func Printhoge(ws *sync.WaitGroup,getnum int64) {
	defer ws.Done()
	fmt.Print(getnum,"\n")
}

