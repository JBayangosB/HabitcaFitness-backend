package main

import (
	"fmt"
	"net/http"
	"strings"
	"appengine"
	"encoding/json"
	"strconv"
	"appengine/datastore"
	// "github.com/gorilla/mux"
)

type Quest struct {
	Name string
	Desc string
	HealthCurrent int64
	HealthMax int64
	MonsterName string
	QuestLogs []QuestLog
	IsEnabled bool
	Region string
}

type QuestLog struct {
	Name string
	Content string
}

type User struct {
	Name string
	Email string
	Password string
	HealthCurrent int64
	HealthMax int64
	ExpCurrent int64
	ExpMax int64
	Gold int64
	DistanceTotal int64 
	Level int64
}

type Reward struct {
	Exp int64
	Gold int64
}

var users map[string]User
var quests map[string]Quest

func inital(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Fitness Habitca!")
	// setupQuests(r)
}

// socialmovieclub-1092.appspot.com/Region/<Melbourne>
// Returns quest particular to region
func returnRegionQuest(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Get id from URL
	parts := strings.Split(r.URL.Path, "/")
	id := parts[2]
	// Get quest
	q := getQuest(id, c)
	js, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return json
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	// log
	c.Infof("ID: [%s]", id)
}



// func getQuest(region string) Quest {
// 	return quests[region]
// }

func userKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Users", "asalieri", 0, nil)
}

// Working!
// socialmovieclub-1092.appspot.com/CreateUser/
// ?email="email@gmail.com"&password="pass1234"&name="Steve"
// Create user on server and returns id
func createUser(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Get variables from url
	query := r.URL.Query()
	name := query.Get("name")
	email := query.Get("email")
	password := query.Get("password")
	if (len(name) < 0 || len(email) < 0 || len(password) < 0) {
		w.Write([]byte("Missing fields"))
	} else {
		c.Infof(name)
	// Create user
		u := User {
			Name: name,
			Email: email,
			Password: password,
			HealthCurrent: 50,
			HealthMax: 50,
			ExpCurrent: 0,
			ExpMax: 100,
			Level: 1,
		}
		// Create key for datastore
		key := datastore.NewIncompleteKey(c, "User", nil)
		// Store in datastore
		if (!doesUserExist(name, c)) {
			if k, err := datastore.Put(c, key, &u); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				id := make(map[string]string)
				stringID := strconv.FormatInt(int64(k.IntID()), 10)
				c.Infof("ID: %s", stringID)
				id["id"] = stringID
				js, err := json.Marshal(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// Return json
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
		} else {
			w.Write([]byte("User already exists"))
		}
	}
	
}


// Check if user already exists in database
func doesUserExist(n string, c appengine.Context) bool {
	q := datastore.NewQuery("User").Filter("Name =", n)
	var results []User
	if _, err := q.GetAll(c, &results); err != nil {
		c.Infof("Error occured with query")
	}
	exists := false
	if ((len(results))> 0) {
		exists = true
	} 
	return exists
}


// }
// socialmovieclub-1092.appspot.com/ReturnUser/?id="333333"
// Returns user
func returnUser(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	// Parse id
	query := r.URL.Query()
	id := query.Get("id")
	// Setup datastore key
	//key := datastore.NewKey(c, "User", "", entityID, nil)
	user := getUser(id, c)
	c.Infof("User name is: %s", user.Name)
	// Convert to json
	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return json
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	// log
	c.Infof("Returned user")

}


func getUser(id string, c appengine.Context) User{
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Infof("Error parsing id")
	}
	k := datastore.NewKey(c, "User", "", intID, nil)
	var user User
	datastore.Get(c, k, &user)
	return user
}

//Needs fixing
// adroit-chemist-144605.appspot.com/RegionAttack/
// ?region="Melbourne"&id="555555"&km="5"
// Attack monster in region and receive reward
func regionAttack(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	// Identify region
	query := r.URL.Query()
	region := query.Get("region")
	quest := getQuest(region, c)
	// Identify user
	u := query.Get("id")
	// Get distance
	km := query.Get("km")
	distance, err := strconv.ParseInt(km, 10, 64)
	c.Infof("distance: [%d]", distance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get user reward by adding distance ran, exp earned
	giveReward(u, distance, r)
	// Return updated quest
	dam := 3*distance
	textDam := strconv.FormatInt(dam, 10)
	l := QuestLog{
		Name: u,
		Content: fmt.Sprintf("%s has attacked %s, dealing %s damage!", u, quest.MonsterName, textDam),
	}
	quest.QuestLogs = append(quest.QuestLogs, l)
	// Update quest
	quest.HealthCurrent -= dam
	updateQuest(quest, c)
	// Convert to json
	js, err := json.Marshal(quest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return json
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	// log
	c.Infof("Returned quest")
}

func getQuest(region string, c appengine.Context) Quest {
	q := datastore.NewQuery("Quest").Filter("Region =", region)
	var results []Quest
	if _, err := q.GetAll(c, &results); err != nil {
		c.Infof("Error occured with query")
	}

	return results[0]
}

func giveReward(id string, km int64, r *http.Request) {
	c := appengine.NewContext(r)
	user := getUser(id, c)
	c.Infof("User found: %s", user.Name)
	user.DistanceTotal += km
	user.ExpCurrent += km*5
	c.Infof("User exp is now [%d] of max: [%d]", user.ExpCurrent, user.ExpMax)
	if(user.ExpMax < user.ExpCurrent){
		// Level up!
		user.Level += 1
		user.ExpCurrent = user.ExpCurrent - user.ExpMax
		user.ExpMax = user.ExpMax + 10
		user.HealthCurrent = user.HealthMax
		c.Infof("User leveled up!")
	}
	updateUser(user, id, c)
	c.Infof("User rewarded")

}

func updateUser(user User, id string, c appengine.Context) {
	// Setup key for user
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Infof("Error parsing id")
	}
	key := datastore.NewKey(c, "User", "", intID, nil)
	// Save to datastore
	if _, err := datastore.Put(c, key, &user); err != nil {
		c.Infof("Error updating user")
		return
	} else {
		c.Infof("User updated")
	}
}

func updateQuest(quest Quest, c appengine.Context) {
	// Setup key for quest
	q := datastore.NewQuery("Quest").Filter("Region =", quest.Region).KeysOnly()
	// var key *datastore.Key = nil
	key, err := q.GetAll(c, nil) 
	if err != nil {
		c.Infof("Error occured with query")
	}
	// Update quest
	if _, err := datastore.Put(c, key[0], &quest); err != nil {
		c.Infof("Error updating quest")
		return
	} else {
		c.Infof("Success updating quest")
	}
}

func handleRequests() {
	http.HandleFunc("/", inital)
	http.HandleFunc("/Region/", returnRegionQuest)
	http.HandleFunc("/CreateUser/", createUser)
	// http.HandleFunc("/ModifyUser/", modifyUser)
	http.HandleFunc("/ReturnUser/", returnUser)
	http.HandleFunc("/RegionAttack/", regionAttack)

}

func init() {
	handleRequests()
}



// func setupQuests(r *http.Request) {
// 	c := appengine.NewContext(r)
// 	// quests = make(map[string]Quest)
// 	q := Quest {
// 		Name: "First Quest",
// 		Desc: "This is your first quest!",
// 		HealthCurrent: 100,
// 		HealthMax: 100,
// 		MonsterName: "Fastrackia",
// 		IsEnabled: true,
// 		Region: "Melbourne",
// 	}

// 	// Create key for datastore
// 	key := datastore.NewIncompleteKey(c, "Quest", nil)
// 	// Store in datastore
// 	if _, err := datastore.Put(c, key, &q); err != nil {
// 		c.Infof("Error setting up quests")
// 		return
// 	} else {
// 		c.Infof("Success setting up quest")
// 	}
// 	// quests["Melbourne"] = q
// }

