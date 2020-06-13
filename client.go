package main

import (
	pb "bitbucket.org/alien_soft/Message/genproto/message"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

type Message struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Time string `json:"time"`
}


func main() {
	conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())

	if err != nil {
		log.Println("error while dialing", err)
		return
	}
	defer conn.Close()

	m := pb.NewMessageServiceClient(conn)

	r := gin.Default()

	r.POST("/message/", func(context *gin.Context) {
		var message Message

		err := context.ShouldBindJSON(&message)

		if err != nil {
			log.Println("error while binding json", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		resp, err := m.Create(context, &pb.Message{Name:message.Name, Time:message.Time})

		if err != nil {
			log.Println("error while creating message", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, resp)
	})

	r.PUT("/message/:id", func(context *gin.Context) {
		var message Message

		err := context.ShouldBindJSON(&message)

		if err != nil {
			log.Println("error while binding json", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		a := context.Param("id")
		id, err := strconv.Atoi(a)

		if err != nil {
			log.Println("param is not uint 64", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		resp, err := m.Update(context, &pb.Message{Id:uint64(id), Name:message.Name, Time:message.Time})

		if err != nil {
			log.Println("error while updating message", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, resp)
	})

	r.DELETE("/message/:id", func(context *gin.Context) {
		var message Message

		err := context.ShouldBindJSON(&message)

		if err != nil {
			log.Println("error while binding json", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		a := context.Param("id")
		id, err := strconv.Atoi(a)

		if err != nil {
			log.Println("param is not uint 64", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		resp, err := m.Delete(context, &pb.DeleteRequest{Id:uint64(id)})

		if err != nil {
			log.Println("error while deleting message", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, resp)
	})

	r.Run(":5054")
}
