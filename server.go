package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Image struct {
	Id     string `json:"id"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Buffer []byte `json:"buffer"`
}

var images []Image

//загрузка запроса без парсирга в структуру
func uploadRawBody(w http.ResponseWriter, r *http.Request) {

	//считываем весь реквест в body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	//создаем структуру
	image := &Image{}

	//парсим json в эту структуру
	err = json.Unmarshal(body, image)

	//формируем ответ передаем в метод структуру и возвращаем ошибку
	err = Resize(*image)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)

	//нет ошибки, вызываем метод добавления картинки в массив
	if err == nil {
		addImage(*image)
	}
}

func Resize(image Image) error {
	s := image
	if s.Id == "" || s.Height == 0 || s.Width == 0 || len(s.Buffer) == 0 {
		return errors.New("error: data is not correct")
	}
	return nil
}

func getImages(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(&images)
}

func getImageId(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	for _, item := range images {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Image{})
}

func addImage(image Image) {

	//рандомим id картинки
	//image.Id = rand.Intn(10000)

	//добавляем картинку в массив картинок
	images = append(images, image)
}

func updateImage(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	for index, item := range images {
		if item.Id == params["id"] {
			images = append(images[:index], images[index+1:]...)
			var image Image
			_ = json.NewDecoder(r.Body).Decode(&image)
			image.Id = params["id"]
			images = append(images, image)
			json.NewEncoder(w).Encode(image)
			return
		}
	}
	json.NewEncoder(w).Encode(images)
}

func main() {

	http.HandleFunc("/struct/", uploadRawBody)

	http.HandleFunc("/getimageId/{id}/", getImageId)

	http.HandleFunc("/getimages/", getImages)

	http.HandleFunc("/updateimage/{id}/", updateImage)

	images = append(images, Image{Id: "5", Height: 300, Width: 350, Buffer: []byte{100}})
	images = append(images, Image{Id: "6", Height: 400, Width: 450, Buffer: []byte{150}})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":45998", nil)
}
