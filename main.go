package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

/*
var (

	server      *gin.Engine
	ctx         context.Context
	collection  *mongo.Collection
	mongoclient *mongo.Client
	err         error

)
*/

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

func addData(c *gin.Context) {
	var newData Data
	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&newData); err != nil {
		return
	}
	dataset = append(dataset, newData)
	c.IndentedJSON(http.StatusCreated, newData)
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
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.Set("mongoDBURL", viper.GetString("mongoDBURL")+viper.GetString("dbname"))
}

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	usersCollection := client.Database(viper.GetString("config.dbname")).Collection("users")
	fmt.Println(usersCollection)

	// insert a single document into a collection
	// create a bson.D object
	user := bson.D{{"fullName", "User 1"}, {"age", 30}}
	// insert the bson object using InsertOne()
	result, err := usersCollection.InsertOne(context.TODO(), user)
	// check for errors in the insertion
	if err != nil {
		panic(err)
	}
	// display the id of the newly inserted object
	fmt.Println(result.InsertedID)

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://" + viper.GetString("config.appUrl") + viper.GetString("config.port"),
			"https://" + viper.GetString("config.appUrl") + viper.GetString("config.port"),
		},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12,
	}))
	// proje 500 ile crash yerse yeniden ayağa kaldırıyor?
	// bütün projeyi try-catch bloğuna sarmak gibi?
	router.Use(gin.Recovery())

	/* GET routes */
	router.GET("/dataset", randomMW(), getDatas)
	router.GET("/dataset/:id", getDataByID)
	//router.GET("/dataset", randomFunc)

	/* POST routes */
	router.POST("/add-data", addData)
	//router.POST("/convert-ids", setIDIntWrongWay)
	router.POST("/convert-ids", setIDIntCorrectWay)
	//router.POST("/convert-ids-with-pointers", setIDIntWrongPointerWay)
	router.POST("/convert-ids-with-pointers", setIDIntCorrectPointerWay)

	fmt.Printf("Project has started.")
	router.Run(viper.GetString("config.appUrl") + viper.GetString("config.port"))
}
