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
	Image string
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
	ID string
}

type Reward struct {
	Exp int64
	Gold int64
}

type Item struct {
	Name string
	Desc string
	Damage int64
	Health int64
	LevelNeeded int64
	GoldNeeded int64
	ID string
	Image string
}

type Inventory struct {
	UserID string
	Items []Item
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
				// Saving id
				u.ID = strconv.FormatInt(int64(k.IntID()), 10)
				key := datastore.NewKey(c, "User", "", k.IntID(), nil)
				// Save to datastore
				if _, err := datastore.Put(c, key, &u); err != nil {
					c.Infof("Error updating user")
					return
				} else {
					c.Infof("user updated")
				}
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

// adroit-chemist-144605.appspot.com/LoginUser/
// ?email=<email>&password=<password>
// returns id
func loginUser(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// Get details
	query := r.URL.Query()
	email := query.Get("email")
	password := query.Get("password")
	// Query database
	q := datastore.NewQuery("User").Filter("Email =", email)
	var results []User
	if _, err := q.GetAll(c, &results); err != nil {
		c.Infof("Error finding email")
	}
	// Check if users found
	if ((len(results)) < 0) {
		w.Write([]byte("No users with that email"))
		return
	}
	// Check if password matches
	u := results[0]
	if u.Password != password {
		w.Write([]byte("Invalid password"))
		return
	}
	// Get user id
	id := u.ID

	// Return user id
	js, err := json.Marshal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return json
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}


// socialmovieclub-1092.appspot.com/ReturnUser/?id="333333"
// Returns user
func returnUser(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	// Parse id
	query := r.URL.Query()
	id := query.Get("id")
	// Setup datastore key
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
	name := getUser(u, c)
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
		Content: fmt.Sprintf("%s has attacked %s, dealing %s damage!", name, quest.MonsterName, textDam),
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
// adroit-chemist-144605.appspot.com/Shop/
func getShop(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)	
	// Get all items in order of gold needed
	q := datastore.NewQuery("Item").Order("GoldNeeded")
	var items []Item
	if _, err := q.GetAll(c, &items); err != nil {
		c.Infof("Error occured getting all items")
	}



	// Convert to json
	js, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return json
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// adroit-chemist-144605.appspot.com/Shop/Create/
// ?name=<Sword>&desc=<Something>&health=<10>&damage=<10>&level=<5>&gold=<50>
func createItem(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// Identify item stats
	query := r.URL.Query()
	health := query.Get("health")
	h, err := strconv.ParseInt(health, 10, 64)
	if err != nil {
		return
	}
	damage := query.Get("damage")
	d, err := strconv.ParseInt(damage, 10, 64)
	if err != nil {
		return
	}
	levelNeeded := query.Get("level")
	l, err := strconv.ParseInt(levelNeeded, 10, 64)
	if err != nil {
		return
	}
	goldNeeded := query.Get("gold")
	g, err := strconv.ParseInt(goldNeeded, 10, 64)
	if err != nil {
		return
	}
	name := query.Get("name")
	desc := query.Get("desc")
	// Make item
	i := Item {
		Health: h,
		Damage: d,
		LevelNeeded: l,
		GoldNeeded: g,
		Name: name,
		Desc: desc,
	}
	// Save item
	// Create key for datastore
	key := datastore.NewIncompleteKey(c, "Item", nil)
	// Store in datastore
	if token, err := datastore.Put(c, key, &i); err != nil {
		c.Infof("Error setting up item")
		return
	} else {
		c.Infof("Success setting up item")
		i.ID = strconv.FormatInt(int64(token.IntID()), 10)
		key := datastore.NewKey(c, "Item", "", token.IntID(), nil)
		// Save to datastore
		if _, err := datastore.Put(c, key, &i); err != nil {
			c.Infof("Error updating item")
			return
		} else {
			c.Infof("item updated")
		}
	}

}

func getItem(itemID string, c appengine.Context) Item {
	id, err := strconv.ParseInt(itemID, 10, 64)
	if err != nil {
		c.Infof("Error getting item")
	}
	k := datastore.NewKey(c, "Item", "", id, nil)
	var item Item
	datastore.Get(c, k, &item)
	return item
}

// adroit-chemist-144605.appspot.com/AddUserItem/
// ?userID=<id>&itemID=<itemID>
func addUserItem(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// Identify item owner and item
	query := r.URL.Query()
	userID := query.Get("userID")
	itemID := query.Get("itemID")
	// Get inventory 
	q := datastore.NewQuery("Inventory").Filter("UserID =", userID)
	var results []Inventory
	if _, err := q.GetAll(c, &results); err != nil {
		c.Infof("Issue getting inventory")
	}

	exists := false
	if ((len(results))> 0) {
		exists = true
	} 


	var inv Inventory
	// If inventory exist
	if exists {
		c.Infof("Existing inventory")
		inv = results[0]
		// Check if user has item
		for _, value := range inv.Items {
			if value.ID == itemID {
				c.Infof("Already has item")
				http.Error(w, "Already has item", 500)
				return
			}
		}
		c.Infof("User does not have item")

		// Can user afford item
		u := getUser(userID, c)
		i := getItem(itemID, c)
		if (i.Name == "") {
			http.Error(w, "Item does not exist", 500)
			return
		}
		if u.Gold < i.GoldNeeded {
			c.Infof("Too poor")
			http.Error(w, "Too poor", 500)
			return
		}
		// Add item to inventory
		inv.Items = append(inv.Items, i)

		// Save inventory
		// Create key for datastore
		key := datastore.NewIncompleteKey(c, "Inventory", nil)
		// Store in datastore
		if _, err := datastore.Put(c, key, &inv); err != nil {
			c.Infof("Error adding to inventory")
			http.Error(w, err.Error(), 500)
			return
		} else {
			c.Infof("Success adding to inventory")
		}
		loseGold(i.GoldNeeded, userID, c)
	} else {
		// New inventory
		c.Infof("New Inventory")
		container := Inventory {
			UserID: userID,
		}
		// Get item
		i := getItem(itemID, c)
		// Set item
		container.Items = append(container.Items, i)

		// Save inventory
		// Create key for datastore
		key := datastore.NewIncompleteKey(c, "Inventory", nil)
		// Store in datastore
		if _, err := datastore.Put(c, key, &container); err != nil {
			c.Infof("Error setting up inventory")
			http.Error(w, err.Error(), 500)
			return
		} else {
			c.Infof("Success setting up inventory")
		}
	}
}

func loseGold(cost int64, userID string, c appengine.Context) {
	u := getUser(userID, c)
	u.Gold -= cost
	updateUser(u, userID, c)
	c.Infof("Cost paid!")
}

func getInventory(userID string, c appengine.Context) Inventory {
	q := datastore.NewQuery("Inventory").Filter("UserID =", userID)
	var inventory []Inventory
	if _, err := q.GetAll(c, &inventory); err != nil {
		c.Infof("Error occured getting all inventory")
	}

	return inventory[0]
}

// adroit-chemist-144605.appspot.com/GetUserItems/
// ?userID=<id>
func getUserItems(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// Identify inventory owner
	query := r.URL.Query()
	userID := query.Get("userID")
	inv := getInventory(userID, c)
	// Convert to json
	js, err := json.Marshal(inv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return json
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleRequests() {
	go http.HandleFunc("/", inital)
	go http.HandleFunc("/Region/", returnRegionQuest)
	go http.HandleFunc("/CreateUser/", createUser)
	go http.HandleFunc("/LoginUser/", loginUser)
	go http.HandleFunc("/ReturnUser/", returnUser)
	go http.HandleFunc("/RegionAttack/", regionAttack)
	go http.HandleFunc("/Shop/", getShop)
	go http.HandleFunc("/Shop/Create/", createItem)
	go http.HandleFunc("/AddUserItem/", addUserItem)
	go http.HandleFunc("/GetUserItems/", getUserItems)
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

