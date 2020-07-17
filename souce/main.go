package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Client struct {
	ID          int
	OrderListId int
	Fname       string
	Lname       string
	Cellnumber  int
	Address     string
	Status      string
	Obs         string
}

type ClientData struct {
	Clientes []Client
}

type User struct {
	ID         int
	UserName   string
	Password   string
	Email      string
	Cellnumber int
	Status     string
	Role       string
	Obs        string
}
type UserData struct {
	Users []User
}
type Product struct {
	Id           int
	PName        string
	Pdescription string
	Valorunidade float64
	Product_cost float64
	Pdatafabi    string
	pImg         string
	Status       string
	Obs          string
}
type DT struct {
	Products []Product
}

//var pImg = ""

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "1234"
	dbName := "cakesdb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var html = template.Must(template.ParseGlob("web/main/*"))

func index(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "Index.html", nil)
}
func RegistrationCreate(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "RegistrationCreate.html", nil)
}
func createClient(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "ClientCreate.html", nil)
}

func productCreate(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "ProductCreate.html", nil)
}

func css(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "Css.html", nil)
}

func header(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "Header.html", nil)
}
func menu(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "Menu.html", nil)
}

func addProduct(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	if r.Method == "POST" {

		file, handler, err := r.FormFile("pimg")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		f, err := os.OpenFile("./images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		nameS := r.FormValue("pname")
		desc := r.FormValue("pdescription")
		vals, _ := strconv.ParseFloat(r.FormValue("valorunidade"), 64)
		cos, _ := strconv.ParseFloat(r.FormValue("product_cost"), 64)
		Fab := r.FormValue("pdatafabi")
		pImg := handler.Filename
		St := r.FormValue("status")
		oB := r.FormValue("obs")
		SC, _ := strconv.Atoi(r.FormValue("sector"))
		Ct, _ := strconv.Atoi(r.FormValue("category"))
		CtS, _ := strconv.Atoi(r.FormValue("sub-category"))
		Brnd, _ := strconv.Atoi(r.FormValue("brand"))
		Mdl, _ := strconv.Atoi(r.FormValue("model"))

		defer db.Close()
		insert, err := db.Prepare("INSERT INTO produtos(pname,pdescription,valorunidade,product_cost,pdatafabi,pimg,status,obs,sector_id,category_id,category_sub_id,brand_id,model_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")
		insert.Exec(nameS, desc, vals, cos, Fab, pImg, St, oB, SC, Ct, CtS, Brnd, Mdl)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
	http.Redirect(w, r, "/loadAllProduct", 301)
	defer db.Close()
}

func loadAllProduct(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	selDB, err := db.Query("SELECT id,pname,pdescription,valorunidade,product_cost,pdatafabi,pimg ,status ,obs FROM produtos  ")
	if err != nil {
		panic(err.Error())
	}

	pro := Product{}
	res := []Product{}

	for selDB.Next() {

		var id int

		var pname, pdescription, pdatafabi, pimg, status, obs string
		var valorunidade, product_cost float64
		err = selDB.Scan(&id, &pname, &pdescription, &valorunidade, &product_cost, &pdatafabi, &pimg, &status, &obs)
		if err != nil {
			panic(err.Error())
		}

		pro.Id = id
		pro.PName = pname
		pro.Pdescription = pdescription
		pro.Valorunidade = valorunidade
		pro.Product_cost = product_cost
		pro.Pdatafabi = pdatafabi
		pro.pImg = pimg
		pro.Status = status
		pro.Obs = obs

		res = append(res, pro)
	}
	dataProduct := DT{
		Products: res,
	}
	fmt.Println(dataProduct)

	html.ExecuteTemplate(w, "Product.html", dataProduct)
	defer db.Close()
}

func addClient(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		orderListId := r.FormValue("orderListId")
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		cellnumber := r.FormValue("cellnumber")
		address := r.FormValue("address")
		status := r.FormValue("status")
		obs := r.FormValue("obs")

		create, error := db.Prepare("INSERT INTO clientes(order_list_id,fname,lname,cellnumber,address, status,obs) VALUES(?,?,?,?,?,?,?)")
		create.Exec(orderListId, fname, lname, cellnumber, address, status, obs)
		if error != nil {
			panic(error.Error())
		}
		defer create.Close()
	}
	defer db.Close()
	http.Redirect(w, r, "/loadAllClient", 301)
}
func addUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")
		cellnumber := r.FormValue("cellnumber")
		status := r.FormValue("status")
		role := r.FormValue("role")
		obs := r.FormValue("obs")

		create, error := db.Prepare("INSERT INTO registration(username, password, email, cellnumber, status, role, obs) VALUES(?,?,?,?,?,?,?)")
		create.Exec(username, password, email, cellnumber, status, role, obs)
		if error != nil {
			panic(error.Error())
		}
		defer create.Close()
	}
	defer db.Close()
	http.Redirect(w, r, "/loadAllUser", 301)
}

func loadAllUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM registration ")
	if err != nil {
		panic(err.Error())
	}

	user := User{}
	res := []User{}

	for selDB.Next() {
		var id, cellnumber int

		var username, password, email, status, role, obs string

		err = selDB.Scan(&id, &username, &password, &email, &cellnumber, &status, &role, &obs)
		if err != nil {
			panic(err.Error())
		}
		user.ID = id
		user.UserName = username
		user.Password = password
		user.Email = email
		user.Cellnumber = cellnumber
		user.Status = status
		user.Role = role
		user.Obs = obs
		res = append(res, user)
	}
	dataUser := UserData{
		Users: res,
	}
	html.ExecuteTemplate(w, "Registration.html", dataUser)
	defer db.Close()
}

func loadAllClient(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM clientes ")
	if err != nil {
		panic(err.Error())
	}

	client := Client{}
	res := []Client{}

	for selDB.Next() {
		var id, order_list_id, cellnumber int

		var fname, lname, address, status, obs string

		err = selDB.Scan(&id, &order_list_id, &fname, &lname, &cellnumber, &address, &status, &obs)
		if err != nil {
			panic(err.Error())
		}
		client.ID = id
		client.OrderListId = order_list_id
		client.Fname = fname
		client.Lname = lname
		client.Cellnumber = cellnumber
		client.Address = address
		client.Status = status
		client.Obs = obs
		res = append(res, client)
	}
	dataClient := ClientData{
		Clientes: res,
	}
	html.ExecuteTemplate(w, "Client.html", dataClient)
	defer db.Close()
}

//--------------------------------------------------------------------
func main() {
	log.Println("Server started on: http://localhost:9095")
	//http.HandleFunc("/blue/", blueHandler)
	http.HandleFunc("/", index)
	http.HandleFunc("/loadAllClient", loadAllClient)
	http.HandleFunc("/loadAllUser", loadAllUser)
	http.HandleFunc("/RegistrationCreate", RegistrationCreate)
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/loadAllProduct", loadAllProduct)
	http.HandleFunc("/addClient", addClient)
	http.HandleFunc("/createClient", createClient)
	//http.HandleFunc("/upload", upload)
	http.HandleFunc("/product/create", productCreate)
	http.HandleFunc("/addProduct", addProduct)
	http.HandleFunc("/css", css)
	http.HandleFunc("/header", header)
	http.HandleFunc("/menu", menu)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	http.ListenAndServe(":9095", nil)
}
