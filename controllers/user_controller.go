package controllers

import (
    "server/configs"
    "server/models"
    "server/responses"
	"server/utils"
    "net/http"
    "time"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/net/context"
    "go.mongodb.org/mongo-driver/bson"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var cityCollection2 *mongo.Collection = configs.GetCollection(configs.DB, "cities")
var validate = validator.New()

func CreateUser(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var user models.User
    defer cancel()

    //validate the request body
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(&user); validationErr != nil {
        return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
    }

	birthDate, _ := utils.ParseDate(user.BirthDateString)

    newUser := models.User{
        Id:       primitive.NewObjectID(),
        DNI:     user.DNI,
        FirstName: user.FirstName,
        LastName:    user.LastName,
        CivilStatus:     user.CivilStatus,
        BirthDate: birthDate,
        CityID:    user.CityID,
        Email: user.Email,
        Phone:    user.Phone,
    }

    result, err := userCollection.InsertOne(ctx, newUser)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
}

func GetAUser(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    dni := c.Param("dni")
    var users []models.User
    defer cancel()

    results, err := userCollection.Find(ctx, bson.M{"dni": dni})

    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleUser models.User
        if err = results.Decode(&singleUser); err != nil {
            return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
        }

        users = append(users, singleUser)
    }

    return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func EditAUser(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Param("userId")
    var user models.User
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(userId)

    //validate the request body
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(&user); validationErr != nil {
        return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
    }

	birthDate, _ := utils.ParseDate(user.BirthDateString)

    update := bson.M{
        "dni": user.DNI, 
        "firstName": user.FirstName, 
        "lastName": user.LastName, 
        "civilStatus": user.CivilStatus, 
        "birthDate": birthDate, 
        "cityId": user.CityID, 
        "email": user.Email, 
        "phone": user.Phone,
    }

    result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //get updated user details
    var updatedUser models.User
    if result.MatchedCount == 1 {
        err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

        if err != nil {
            return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
        }

        var cityValue models.City
        objId, _ := primitive.ObjectIDFromHex(updatedUser.CityID)
        err2 := cityCollection2.FindOne(ctx, bson.M{"key": objId}).Decode(&cityValue)
    
        if err2 != nil {
            return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err2.Error()}})
        }
        updatedUser.City = cityValue.Value
    }

    return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedUser}})
}

func DeleteAUser(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Param("userId")
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(userId)

    result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})

    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    if result.DeletedCount < 1 {
        return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": "User with specified ID not found!"}})
    }

    return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "User successfully deleted!"}})
}

func GetAllUsers(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var users []models.User
    defer cancel()

    results, err := userCollection.Find(ctx, bson.M{})

    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleUser models.User
        if err = results.Decode(&singleUser); err != nil {
            return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
        }
        
        var cityValue models.City
        objId, _ := primitive.ObjectIDFromHex(singleUser.CityID)
        err2 := cityCollection2.FindOne(ctx, bson.M{"key": objId}).Decode(&cityValue)
    
        if err2 != nil {
            return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err2.Error()}})
        }
        singleUser.City = cityValue.Value

        users = append(users, singleUser)
    }

    return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}