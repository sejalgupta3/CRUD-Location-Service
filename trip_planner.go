package main
import (
   	"fmt"
    "httprouter"
    "net/http"
    "encoding/json"
    "log"
    "io/ioutil"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "math/rand"
    "strconv"
    "strings"
)

type User struct{
	UserId int
	Name string
	UserAddress Address
}

type Address struct{
	Address string
	City string
	State string
	Zip string
	Coordinates Location
}

type Location struct{
	Latitude float64
	Longitude float64
}

type AddLocationRequest struct{
	Name string
	Address string
	City string
	State string
	Zip string
}

type UpdateLocationRequest struct{
	Address string
	City string
	State string
	Zip string
}

type AddLocationResponse struct{
	UserId int
	Name string
	Address string
	City string
	State string
	Zip string
	Coordinates Location
}

func getCoordinates(a *Address){
	addressString := strings.Replace(a.Address+"+"+a.City+"+"+a.State+"+"+a.Zip, " ", "%20", -1)
   	resp, err := http.Get("http://maps.google.com/maps/api/geocode/json?address="+addressString+"&sensor=false")
   	fmt.Println("http://maps.google.com/maps/api/geocode/json?address="+addressString+"&sensor=false")
	if(err == nil){
	    body, err := ioutil.ReadAll(resp.Body)
	    if(err == nil) {
	        var data interface{}
	        json.Unmarshal(body, &data)
	        var m = data.(map[string] interface{})            
	        var articles = m["results"].([]interface{})[0].(map[string]interface{})["geometry"].(map[string]interface{})["location"]
	        lat := articles.(map[string]interface{})["lat"].(float64)
	        lng := articles.(map[string]interface{})["lng"].(float64)
	        a.Coordinates.Latitude = lat
	        a.Coordinates.Longitude = lng
	    } else {
	        fmt.Println(err)
	    }
	} else {
	    fmt.Println(err)
	}
}

func createLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		addLocationRequest := new(AddLocationRequest)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&addLocationRequest)
		if error != nil {
			log.Println(error.Error())
			http.Error(rw, error.Error(), http.StatusInternalServerError)
			return
		}
		
		location := Location{}
		address := Address{addLocationRequest.Address,addLocationRequest.City,addLocationRequest.State,addLocationRequest.Zip,location}
		user := User{0,addLocationRequest.Name,address}
		user.UserId = rand.Intn(1000)
		getCoordinates(&user.UserAddress)
		session, err := mgo.Dial("mongodb://sejal:1234@ds045064.mongolab.com:45064/planner")
		if err != nil {
		        panic(err)
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("planner").C("user")
		err = c.Insert(&user)
		if err != nil {
		        log.Fatal(err)
		}
		result := User{}
		err = c.Find(bson.M{"userid":user.UserId}).One(&result)
		if err != nil {
		        log.Fatal(err)
		}
		outgoingJSON, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		fmt.Fprint(rw, string(outgoingJSON))
}

func getLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	locationId,_ := strconv.Atoi(p.ByName("location_id"))
	session, err := mgo.Dial("mongodb://sejal:1234@ds045064.mongolab.com:45064/planner")
    if err != nil {
            panic(err)
    }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    c := session.DB("planner").C("user")
    result := User{}
    err = c.Find(bson.M{"userid":locationId}).One(&result)
    if err != nil {
            log.Fatal(err)
    }
    outgoingJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
    fmt.Fprint(rw, string(outgoingJSON))
}

func putLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		locationId,_ := strconv.Atoi(p.ByName("location_id"))
		updateLocationRequest := new(UpdateLocationRequest)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&updateLocationRequest)
		if error != nil {
			log.Println(error.Error())
			http.Error(rw, error.Error(), http.StatusInternalServerError)
			return
		}
		
		location := Location{}
		address := Address{updateLocationRequest.Address,updateLocationRequest.City,updateLocationRequest.State,updateLocationRequest.Zip,location}
		
		session, err := mgo.Dial("mongodb://sejal:1234@ds045064.mongolab.com:45064/planner")
		if err != nil {
		        panic(err)
		}
		
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("planner").C("user")
		
		result := User{}
    	err = c.Find(bson.M{"userid":locationId}).One(&result)
	    if err != nil {
	            log.Fatal(err)
	    }
	    
	    userName := result.Name    
		user := User{locationId,userName,address}
		getCoordinates(&user.UserAddress)
		
		colQuerier := bson.M{"userid":locationId}
		change := bson.M{"$set": bson.M{"useraddress": user.UserAddress}}
		err = c.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}
		
		result2 := User{}
		err = c.Find(bson.M{"userid":locationId}).One(&result2)
		if err != nil {
		        log.Fatal(err)
		}
		outgoingJSON, err := json.Marshal(result2)
		if err != nil {
			log.Fatal(err)
		}
		
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		fmt.Fprint(rw, string(outgoingJSON))
}

func deleteLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	locationId,_ := strconv.Atoi(p.ByName("location_id"))
	session, err := mgo.Dial("mongodb://sejal:1234@ds045064.mongolab.com:45064/planner")
    if err != nil {
            panic(err)
    }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    c := session.DB("planner").C("user")
    err = c.Remove(bson.M{"userid":locationId})
    if err != nil {
            log.Fatal(err)
    }
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
    fmt.Fprint(rw, "User Location deleted successfully")
}

func main() {
    router := httprouter.New()
    router.POST("/locations", createLocation)
    router.GET("/locations/:location_id", getLocation)
    router.PUT("/locations/:location_id", putLocation)
    router.DELETE ("/locations/:location_id", deleteLocation)
    server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: router,
    }
    server.ListenAndServe()
}

