package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (

	//server      *gin.Engine
	//ctx         context.Context
	//collection  *mongo.Collection
	client *mongo.Client
	users  *mongo.Collection
	//err        error

)

type Data struct {
	Name     string  `json:"name" binding:"required"` // validation safhasında mongo şemasında required yapmak gibi?
	Language string  `json:"language"`
	ID       string  `json:"id" binding:"required"` //json outputta ID yerine id şeklinde yer alır
	Bio      string  `json:"bio"`
	Version  float64 `json:"version"`
}

type User struct {
	FullName string `json:"name" bson:"fullName"`
	Age      int    `json:"age" bson:"age"`
}

type Params struct {
	VersionStartValue float64 `json:"versionStartValue" `
	VersionEndValue   float64 `json:"versionEndValue" `
	Lang1             string  `json:"lang1"`
	Lang2             string  `json:"lang2"`
}

var dataset = []Data{
	{
		Name:     "Adeel Solangi",
		Language: "Sindhi",
		ID:       "V59OF92YF627HFY0",
		Bio:      "Donec lobortis eleifend condimentum. Cras dictum dolor lacinia lectus vehicula rutrum. Maecenas quis nisi nunc. Nam tristique feugiat est vitae mollis. Maecenas quis nisi nunc.",
		Version:  6.1,
	},
	{
		Name:     "Afzal Ghaffar",
		Language: "Sindhi",
		ID:       "ENTOCR13RSCLZ6KU",
		Bio:      "Aliquam sollicitudin ante ligula, eget malesuada nibh efficitur et. Pellentesque massa sem, scelerisque sit amet odio id, cursus tempor urna. Etiam congue dignissim volutpat. Vestibulum pharetra libero et velit gravida euismod.",
		Version:  1.88,
	},
	{
		Name:     "Aamir Solangi",
		Language: "Sindhi",
		ID:       "IAKPO3R4761JDRVG",
		Bio:      "Vestibulum pharetra libero et velit gravida euismod. Quisque mauris ligula, efficitur porttitor sodales ac, lacinia non ex. Fusce eu ultrices elit, vel posuere neque.",
		Version:  7.27,
	},
	{
		Name:     "Abla Dilmurat",
		Language: "Uyghur",
		ID:       "5ZVOEPMJUI4MB4EN",
		Bio:      "Donec lobortis eleifend condimentum. Morbi ac tellus erat.",
		Version:  2.53,
	},
	{
		Name:     "Adil Eli",
		Language: "Uyghur",
		ID:       "6VTI8X6LL0MMPJCC",
		Bio:      "Vivamus id faucibus velit, id posuere leo. Morbi vitae nisi lacinia, laoreet lorem nec, egestas orci. Suspendisse potenti.",
		Version:  6.49,
	},
	{
		Name:     "Adile Qadir",
		Language: "Uyghur",
		ID:       "F2KEU5L7EHYSYFTT",
		Bio:      "Duis commodo orci ut dolor iaculis facilisis. Morbi ultricies consequat ligula posuere eleifend. Aenean finibus in tortor vel aliquet. Fusce eu ultrices elit, vel posuere neque.",
		Version:  1.9,
	},
	{
		Name:     "Abdukerim Ibrahim",
		Language: "Uyghur",
		ID:       "LO6DVTZLRK68528I",
		Bio:      "Vivamus id faucibus velit, id posuere leo. Nunc aliquet sodales nunc a pulvinar. Nunc aliquet sodales nunc a pulvinar. Ut viverra quis eros eu tincidunt.",
		Version:  5.9,
	},
	{
		Name:     "Adil Abro",
		Language: "Sindhi",
		ID:       "LJRIULRNJFCNZJAJ",
		Bio:      "Etiam malesuada blandit erat, nec ultricies leo maximus sed. Fusce congue aliquam elit ut luctus. Etiam malesuada blandit erat, nec ultricies leo maximus sed. Cras dictum dolor lacinia lectus vehicula rutrum. Integer vehicula, arcu sit amet egestas efficitur, orci justo interdum massa, eget ullamcorper risus ligula tristique libero.",
		Version:  9.32,
	},
	{
		Name:     "Afonso Vilarchán",
		Language: "Galician",
		ID:       "JMCL0CXNXHPL1GBC",
		Bio:      "Fusce eu ultrices elit, vel posuere neque. Morbi ac tellus erat. Nunc tincidunt laoreet laoreet.",
		Version:  5.21,
	},
	{
		Name:     "Mark Schembri",
		Language: "Maltese",
		ID:       "KU4T500C830697CW",
		Bio:      "Nam laoreet, nunc non suscipit interdum, justo turpis vestibulum massa, non vulputate ex urna at purus. Morbi ultricies consequat ligula posuere eleifend. Vivamus id faucibus velit, id posuere leo. Sed laoreet posuere sapien, ut feugiat nibh gravida at. Ut maximus, libero nec facilisis fringilla, ex sem sollicitudin leo, non congue tortor ligula in eros.",
		Version:  3.17,
	},
	{
		Name:     "Antía Sixirei",
		Language: "Galician",
		ID:       "XOF91ZR7MHV1TXRS",
		Bio:      "Pellentesque massa sem, scelerisque sit amet odio id, cursus tempor urna. Phasellus massa ligula, hendrerit eget efficitur eget, tincidunt in ligula. Morbi finibus dui sed est fringilla ornare. Duis pellentesque ultrices convallis. Morbi ultricies consequat ligula posuere eleifend.",
		Version:  6.44,
	},
	{
		Name:     "Aygul Mutellip",
		Language: "Uyghur",
		ID:       "FTSNV411G5MKLPDT",
		Bio:      "Duis commodo orci ut dolor iaculis facilisis. Nam semper gravida nunc, sit amet elementum ipsum. Donec pellentesque ultrices mi, non consectetur eros luctus non. Pellentesque massa sem, scelerisque sit amet odio id, cursus tempor urna.",
		Version:  9.1,
	},
	{
		Name:     "Awais Shaikh",
		Language: "Sindhi",
		ID:       "OJMWMEEQWMLDU29P",
		Bio:      "Nunc aliquet sodales nunc a pulvinar. Ut dictum, ligula eget sagittis maximus, tellus mi varius ex, a accumsan justo tellus vitae leo. Donec pellentesque ultrices mi, non consectetur eros luctus non. Nulla finibus massa at viverra facilisis. Nunc tincidunt laoreet laoreet.",
		Version:  1.59,
	},
	{
		Name:     "Ambreen Ahmed",
		Language: "Sindhi",
		ID:       "5G646V7E6TJW8X2M",
		Bio:      "Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Etiam consequat enim lorem, at tincidunt velit ultricies et. Ut maximus, libero nec facilisis fringilla, ex sem sollicitudin leo, non congue tortor ligula in eros.",
		Version:  2.35,
	},
	{
		Name:     "Celtia Anes",
		Language: "Galician",
		ID:       "Z53AJY7WUYPLAWC9",
		Bio:      "Nullam ac sodales dolor, eu facilisis dui. Maecenas non arcu nulla. Ut viverra quis eros eu tincidunt. Curabitur quis commodo quam.",
		Version:  8.34,
	},
	{
		Name:     "George Mifsud",
		Language: "Maltese",
		ID:       "N1AS6UFULO6WGTLB",
		Bio:      "Phasellus tincidunt sollicitudin posuere. Ut accumsan, est vel fringilla varius, purus augue blandit nisl, eu rhoncus ligula purus vel dolor. Donec congue sapien vel euismod interdum. Cras dictum dolor lacinia lectus vehicula rutrum. Phasellus massa ligula, hendrerit eget efficitur eget, tincidunt in ligula.",
		Version:  7.47,
	},
	{
		Name:     "Aytürk Qasim",
		Language: "Uyghur",
		ID:       "70RODUVRD95CLOJL",
		Bio:      "Curabitur ultricies id urna nec ultrices. Aliquam scelerisque pretium tellus, sed accumsan est ultrices id. Duis commodo orci ut dolor iaculis facilisis.",
		Version:  1.32,
	},
	{
		Name:     "Dialè Meso",
		Language: "Sesotho sa Leboa",
		ID:       "VBLI24FKF7VV6BWE",
		Bio:      "Maecenas non arcu nulla. Vivamus id faucibus velit, id posuere leo. Nullam sodales convallis mauris, sit amet lobortis magna auctor sit amet.",
		Version:  6.29,
	},
	{
		Name:     "Breixo Galáns",
		Language: "Galician",
		ID:       "4VRLON0GPEZYFCVL",
		Bio:      "Integer vehicula, arcu sit amet egestas efficitur, orci justo interdum massa, eget ullamcorper risus ligula tristique libero. Morbi ac tellus erat. In id elit malesuada, pulvinar mi eu, imperdiet nulla. Vestibulum pharetra libero et velit gravida euismod. Cras dictum dolor lacinia lectus vehicula rutrum.",
		Version:  1.62,
	},
	{
		Name:     "Bieito Lorme",
		Language: "Galician",
		ID:       "5DRDI1QLRGLP29RC",
		Bio:      "Ut viverra quis eros eu tincidunt. Morbi vitae nisi lacinia, laoreet lorem nec, egestas orci. Curabitur quis commodo quam. Morbi ac tellus erat.",
		Version:  4.45,
	},
}

func getDatas(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, dataset)
}

func getDatasMongo(c *gin.Context) {

	if client == nil {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
		fmt.Println("client gelmemiş")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the MongoDB collection if not done already
	if users == nil {
		users = client.Database(viper.GetString("config.dbname")).Collection("users")
	}
	opts := options.Find().SetLimit(5)
	filter := bson.M{}
	// = Model.find({})
	cursor, err := users.Find(context.Background(), filter, opts)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cursor.Close(context.Background())

	var results []Data
	for cursor.Next(context.Background()) {
		var data Data
		err := cursor.Decode(&data)
		if err != nil {
			fmt.Println(err)
			// Handle error as needed
		}
		results = append(results, data)
	}

	// Send the results as JSON
	//c.JSON(http.StatusOK, results)
	c.HTML(http.StatusOK, "users.html", gin.H{
		"results": results,
	})
}

func addData(c *gin.Context) {
	var newData Data
	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&newData); err != nil {
		return
	}
	dataset = append(dataset, newData)
	c.IndentedJSON(http.StatusCreated, newData)
}

func updateData(c *gin.Context) {
	var changedUser Data

	if err := c.ShouldBindJSON(&changedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if client == nil {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
		fmt.Println("client gelmemiş")
		if err != nil {
			log.Fatal(err)
		}
	}
	// Initialize the MongoDB collection if not done already
	if users == nil {
		users = client.Database(viper.GetString("config.dbname")).Collection("users")
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.D{{Key: "id", Value: changedUser.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: changedUser.Name},
		{Key: "language", Value: changedUser.Language},
		{Key: "bio", Value: changedUser.Bio},
		{Key: "version", Value: changedUser.Version},
	}}}
	// Set a timeout for the operation. --sorgu gereğinden fazla uzun sürerse servisin yanıtsız kalmasını engeller
	// 5 saniye boyunca bulamazsa cancel fonksiyonu çalışır ?
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result Data
	users.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	// Return the result as JSON.
	c.JSON(http.StatusOK, result)
}

func addDataMongo(c *gin.Context) {
	var newUser Data

	if client == nil {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
		fmt.Println("client gelmemiş")
		if err != nil {
			log.Fatal(err)
		}
	}

	users := client.Database(viper.GetString("config.dbname")).Collection("users")
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set a timeout for the operation.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := users.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	//fmt.Println(result)
	// yukarıkdaki result, insertedUser'ın id'sini döndürüyormuş
	var insertedUser Data
	err = users.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&insertedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated document"})
		return
	}

	// Return a success response.
	c.JSON(http.StatusOK, insertedUser)
}

func getDataByID(c *gin.Context) {
	id := c.Param("id")
	for _, a := range dataset {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func getDataByIDMongo(c *gin.Context) {
	if client == nil {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
		fmt.Println("client gelmemiş")
		if err != nil {
			log.Fatal(err)
		}
	}
	// Initialize the MongoDB collection if not done already
	if users == nil {
		users = client.Database(viper.GetString("config.dbname")).Collection("users")
	}
	id := c.Param("id")
	filter := bson.D{{Key: "id", Value: id}}
	var result Data
	// Set a timeout for the operation. --sorgu gereğinden fazla uzun sürerse servisin yanıtsız kalmasını engeller
	// 5 saniye boyunca bulamazsa cancel fonksiyonu çalışır ?
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := users.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.HTML(http.StatusOK, "user-info.html", gin.H{
		"id":       result.ID,
		"name":     result.Name,
		"language": result.Language,
		"bio":      result.Bio,
		"version":  result.Version,
	})
}

func getDatasBetweenGivenIdValuesMongo(c *gin.Context) {

	if client == nil {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
		fmt.Println("client gelmemiş")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the MongoDB collection if not done already
	if users == nil {
		users = client.Database(viper.GetString("config.dbname")).Collection("users")
	}
	startValue := c.Param("startValue")
	endValue := c.Param("endValue")

	opts := options.Find().SetSort(bson.D{{"id", 1}})
	filter := bson.M{"id": bson.M{"$gte": startValue, "$lt": endValue}}

	cursor, err := users.Find(context.Background(), filter, opts)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cursor.Close(context.Background())

	var results []Data
	for cursor.Next(context.Background()) {
		var data Data
		err := cursor.Decode(&data)
		if err != nil {
			fmt.Println(err)
			// Handle error as needed
		}
		results = append(results, data)
	}

	// Send the results as JSON
	c.JSON(http.StatusOK, results)
}

func getDatasBetweenGivenVersionValuesMongo(c *gin.Context) {

	if client == nil {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
		fmt.Println("client gelmemiş")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the MongoDB collection if not done already
	if users == nil {
		users = client.Database(viper.GetString("config.dbname")).Collection("users")
	}
	startValue, _ := strconv.ParseFloat(c.Query("startValue"), 64)
	endValue, _ := strconv.ParseFloat(c.Query("endValue"), 64)

	opts := options.Find().SetSort(bson.D{{"version", 1}})
	filter := bson.M{"version": bson.M{"$gte": startValue, "$lt": endValue}}

	cursor, err := users.Find(context.Background(), filter, opts)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cursor.Close(context.Background())

	var results []Data
	for cursor.Next(context.Background()) {
		var data Data
		err := cursor.Decode(&data)
		if err != nil {
			fmt.Println(err)
			// Handle error as needed
		}
		results = append(results, data)
	}

	// Send the results as JSON
	c.JSON(http.StatusOK, results)
}

func getDatasBetweenGivenVersionAndLangValuesMongo(c *gin.Context) {

	if client == nil {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
		fmt.Println("client gelmemiş")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the MongoDB collection if not done already
	if users == nil {
		users = client.Database(viper.GetString("config.dbname")).Collection("users")
	}
	sv, _ := strconv.ParseFloat(c.Query("sv"), 64)
	ev, _ := strconv.ParseFloat(c.Query("ev"), 64)

	opts := options.Find().SetSort(bson.D{{Key: "version", Value: 1}})
	filter := bson.M{"version": bson.M{"$gte": sv, "$lt": ev}, "language": c.Query("lang")}
	fmt.Println(filter)

	cursor, err := users.Find(context.Background(), filter, opts)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cursor.Close(context.Background())

	var results []Data
	for cursor.Next(context.Background()) {
		var data Data
		err := cursor.Decode(&data)
		if err != nil {
			fmt.Println(err)
			// Handle error as needed
		}
		results = append(results, data)
	}

	// Send the results as JSON
	c.JSON(http.StatusOK, results)
}

func getUsersWhoSpeakTwoLanguages(c *gin.Context) {

	if client == nil {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
		fmt.Println("client gelmemiş")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the MongoDB collection if not done already
	if users == nil {
		users = client.Database(viper.GetString("config.dbname")).Collection("users")
	}

	opts := options.Find().SetSort(bson.D{{Key: "version", Value: 1}}) // version'a göre artacak şekilde
	filter := bson.D{
		{Key: "$or",
			Value: bson.A{
				bson.D{{Key: "language", Value: c.Query("lang1")}},
				bson.D{{Key: "language", Value: c.Query("lang2")}},
			}},
	}

	cursor, err := users.Find(context.Background(), filter, opts)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cursor.Close(context.Background())

	var results []Data
	for cursor.Next(context.Background()) {
		var data Data
		err := cursor.Decode(&data)
		if err != nil {
			fmt.Println(err)
			// Handle error as needed
		}
		results = append(results, data)
	}

	// Send the results as JSON
	c.JSON(http.StatusOK, results)
}

/*
func setIDIntWrongWay(c *gin.Context) {
	for i, a := range dataset {
		fmt.Println(i)
		a.ID = strconv.Itoa(i)
		// a değeri dataset'in ilgili değerinden bağımsız yeni bir değişken olduğu için update ederken kullanılamaz.
		fmt.Println(a.ID)
	}
	c.IndentedJSON(http.StatusOK, dataset)
}
*/

func setIDIntCorrectWay(c *gin.Context) {
	for i, _ := range dataset {
		dataset[i].ID = strconv.Itoa(i)
	}
	c.IndentedJSON(http.StatusOK, dataset)
}

func convertIDsMongo(c *gin.Context) {

	if client == nil {
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
		fmt.Println("client gelmemiş")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the MongoDB collection if not done already
	if users == nil {
		users = client.Database(viper.GetString("config.dbname")).Collection("users")
	}

	// Set a timeout for the operation.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find all documents in the users.
	cursor, err := users.Find(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer cursor.Close(ctx)

	// Update each document with a new countable ID.
	counter := 0
	var results []Data
	for cursor.Next(ctx) {
		var user Data
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		// Update the document's ID.
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "id", Value: strconv.Itoa(counter)},
			}},
		}

		/*
			_, err := collection.UpdateOne(ctx, bson.M{"id": user.ID}, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}

			// ayrı bir sorgu ile result slice'ına eklemek gerekiyor sanırım ??
			// veya res.send('ok') diyip geçmek gerekiyor.
			var updatedUser Data
			err = collection.FindOne(ctx, bson.M{"id": strconv.Itoa(counter)}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated document"})
				return
			}
		*/

		// Configure the options for FindOneAndUpdate.
		options := options.FindOneAndUpdate().SetReturnDocument(options.After)

		// Find, update, and get the document.
		var updatedUser Data
		err := users.FindOneAndUpdate(ctx, bson.M{"id": user.ID}, update, options).Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		results = append(results, updatedUser)
		counter++

	}

	c.IndentedJSON(http.StatusOK, results)
}

/*
func setIDIntWrongPointerWay(c *gin.Context) {
	for i, a := range dataset {
		pr := &a
		(*pr).ID = strconv.Itoa(i)
		fmt.Println(a.ID)
	}
	c.IndentedJSON(http.StatusOK, dataset)
}
*/

func setIDIntCorrectPointerWay(c *gin.Context) {
	for i, _ := range dataset {
		pr := &dataset[i]
		(*pr).ID = strconv.Itoa(i)
	}
	c.IndentedJSON(http.StatusOK, dataset)
}

/*
func test(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
*/

// Middlewares
func randomMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("MW çalıştı.\n")
		c.Next()
	}
}

/* func randomFunc(c *gin.Context) {
	// Perform a MongoDB query to get BSON data
	result := collection.FindOne(context.Background(), bson.M{"fullName": "User 1"})
	if result.Err() != nil {
		fmt.Println("Error fetching data from MongoDB:", result.Err())
		return
	}
	fmt.Println(result)

	// Decode BSON data into a Go struct
	var user User
	err = result.Decode(&user)
	if err != nil {
		fmt.Println("Error decoding BSON data:", err)
		return
	}
	c.JSON(http.StatusOK, user)
} */

func init() {
	defer fmt.Printf("Initializing complete.\n")
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.Set("mongoDBURL", viper.GetString("mongoDBURL")+viper.GetString("dbname"))
	// executed only once

	// var err error
	// initte yapılıp global variable'lar arasında eklenmeli ??

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	// global variable'lar collection'lar ile doldurulur
	users = client.Database(viper.GetString("config.dbname")).Collection("users")

	filter := bson.M{}
	count, err := users.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	fmt.Println("COUNT NUMBER: ", count)
	if count <= 0 {
		// create fake data DB
		for _, data := range dataset {
			_, err := users.InsertOne(context.Background(), data)
			if err != nil {
				fmt.Println(err)
				// Handle errors
			}
		}
	}
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
		"loop": func(from, to int) <-chan int {
			ch := make(chan int)
			go func() {
				for i := from; i <= to; i++ {
					ch <- i
				}
				close(ch)
			}()
			return ch
		},
	})
	//router.Static("/assets", "./assets")
	router.LoadHTMLGlob("./assets/templates/*.html")

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://" + viper.GetString("config.appUrl") + viper.GetString("config.port"),
			"https://" + viper.GetString("config.appUrl") + viper.GetString("config.port"),
		},
		AllowMethods:     []string{"GET", "POST", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12,
	}))

	// proje 500 ile crash yerse yeniden ayağa kaldırıyor?
	// bütün projeyi try-catch bloğuna sarmak gibi?
	router.Use(gin.Recovery())

	/* GET routes */
	router.GET("/dataset", getDatas)
	router.GET("/mongodb", randomMW(), getDatasMongo)
	router.GET("/mongodb/either-languages", getUsersWhoSpeakTwoLanguages)
	router.GET("/mongodb/between-version-values", getDatasBetweenGivenVersionValuesMongo)
	router.GET("/mongodb/between-version-and-lang-values", getDatasBetweenGivenVersionAndLangValuesMongo)
	router.GET("/dataset/:id", getDataByID)
	router.GET("/mongodb/:id", getDataByIDMongo)
	router.GET("/mongodb/between-id-values/:startValue/:endValue", getDatasBetweenGivenIdValuesMongo)

	//router.GET("/test", test)
	router.POST("/test", testFunc)

	/* POST routes */
	router.POST("/add-data", addData)
	router.POST("/update-user", updateData)
	router.POST("/add-data-mongo", addDataMongo)
	//router.POST("/convert-ids", setIDIntWrongWay)
	router.PATCH("/convert-ids", setIDIntCorrectWay)
	router.PATCH("/convert-mongo-ids", convertIDsMongo)
	//router.PATCH("/convert-ids-with-pointers", setIDIntWrongPointerWay)
	router.PATCH("/convert-ids-with-pointers", setIDIntCorrectPointerWay)

	router.Run(viper.GetString("config.appUrl") + viper.GetString("config.port"))
	fmt.Printf("Project has started.")
}
