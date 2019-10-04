package apiserver

import (
	"os"
	"fmt"
	"net/http"
	"github.com/rgreen312/owlplace/server/consensus"
	"html/template"
)


type ApiServer struct {
	sendc chan consensus.BackendMessage
	recvc chan consensus.ConsensusMessage
}

func NewApiServer(send_channel chan consensus.BackendMessage, recv_channel chan consensus.ConsensusMessage) *ApiServer {
	return &ApiServer{
		sendc: send_channel,
		recvc: recv_channel,
	}
}

func (api *ApiServer) ListenAndServe(){
	http.HandleFunc("/headers", api.Headers)
	http.HandleFunc("/get_image", api.GetImage)
	http.HandleFunc("/update_pixel", api.UpdatePixel)
	http.ListenAndServe(":3000", nil)
}


func (api *ApiServer) GetImage(w http.ResponseWriter, req *http.Request) {
	// Debug message
	fmt.Fprintf(os.Stdout, "Getting Image From Raft\n")
	// Construct the message
	m := consensus.BackendMessage{ Type: consensus.GET_IMAGE }
	// Send a message through the channel
	api.sendc <- m
	var ImageTemplate string = `<!DOCTYPE html>
								<html lang="en"><head></head>
								<body><img src="data:image/jpg;base64,{{.Image}}"></body>`
	image_msg := <- api.recvc

	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		fmt.Fprintf(os.Stdout, "Unable to parse image template.\n")
	} else {
		data := map[string]interface{}{"Image": image_msg.Data}
		if err = tmpl.Execute(w, data); err != nil {
			fmt.Fprintf(os.Stdout, "Unable to execute template.\n")
		}
	}

}

func (api *ApiServer) UpdatePixel(w http.ResponseWriter, req *http.Request) {

	// Decode the request
	update := req.URL.Query().Get("update")
	if(update != ""){
		fmt.Fprintf(os.Stdout, update)
		// Testing with some dummy data
		m := consensus.BackendMessage{ Type: consensus.UPDATE_PIXEL, Data: update }
		api.sendc <- m
		image_msg := <- api.recvc
		fmt.Fprintf(os.Stdout, image_msg.Data)
	}
}

func (api *ApiServer) GetLastUserModification(user_id string) string{

	// Testing with some dummy data
	m := consensus.BackendMessage{ Type: consensus.GET_LAST_USER_UPDATE, Data: fmt.Sprintf("get %s", user_id)}
	api.sendc <- m
	image_msg := <- api.recvc
	return image_msg.Data
}

func (api *ApiServer) SetLastUserModification(user_id string, last_modification string) bool{

	// Testing with some dummy data
	m := consensus.BackendMessage{ Type: consensus.SET_LAST_USER_UPDATE, Data: fmt.Sprintf("put %s %s", user_id, last_modification)}
	api.sendc <- m
	image_msg := <- api.recvc
	return image_msg.Type == consensus.SUCCESS
}

func (api *ApiServer) Headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
