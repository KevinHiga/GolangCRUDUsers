package controllers

import (
    "server/configs"
    "server/models"
    "server/responses"
    "net/http"
    "time"
	"strings"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/net/context"
    "go.mongodb.org/mongo-driver/bson"
)

var cityCollection *mongo.Collection = configs.GetCollection(configs.DB, "cities")
var validateCity = validator.New()

func CreateCity(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var user models.City
    defer cancel()

    //validate the request body
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validateCity.Struct(&user); validationErr != nil {
        return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
    }

    newCity := models.City{
        Key:       primitive.NewObjectID(),
        Value:     user.Value,
    }

    result, err := cityCollection.InsertOne(ctx, newCity)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
}

func GetACity(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    cityId := c.Param("cityId")
    var user models.City
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(cityId)

    err := cityCollection.FindOne(ctx, bson.M{"key": objId}).Decode(&user)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
}

func EditACity(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    cityId := c.Param("cityId")
    var user models.City
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(cityId)

    //validate the request body
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validateCity.Struct(&user); validationErr != nil {
        return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
    }

    update := bson.M{
        "value": user.Value, 
    }

    result, err := cityCollection.UpdateOne(ctx, bson.M{"key": objId}, bson.M{"$set": update})

    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //get updated user details
    var updatedCity models.City
    if result.MatchedCount == 1 {
        err := cityCollection.FindOne(ctx, bson.M{"key": objId}).Decode(&updatedCity)

        if err != nil {
            return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
        }
    }

    return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedCity}})
}

func DeleteACity(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    cityId := c.Param("cityId")
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(cityId)

    result, err := cityCollection.DeleteOne(ctx, bson.M{"key": objId})

    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    if result.DeletedCount < 1 {
        return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": "City with specified ID not found!"}})
    }

    return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "City successfully deleted!"}})
}

func GetAllCitys(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var users []models.City
    defer cancel()

    results, err := cityCollection.Find(ctx, bson.M{})

    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleCity models.City
        if err = results.Decode(&singleCity); err != nil {
            return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
        }

        users = append(users, singleCity)
    }

    return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func GetAllCitysByKey(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    city := strings.ToLower(c.Param("city"))
    var cities []models.City
    defer cancel()

    filter := bson.M{"value": primitive.Regex{Pattern: "^" + city + "$", Options: "i"}}
    results, err := cityCollection.Find(ctx, filter)

    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
    }

    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleCity models.City
        if err = results.Decode(&singleCity); err != nil {
            return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
        }

        cities = append(cities, singleCity)
    }

    return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": cities}})
}