package handlers

import (
	"context"
	"curd-service/common"
	"curd-service/constants"
	"curd-service/database"
	"curd-service/helpers"
	"curd-service/logger"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func Index(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Printf("Home Page is call")
	w.Write([]byte("HOME PUBLIC INDEX PAGE"))
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	if !constants.IsAccess(r.Header.Get("Role"), constants.READ) {
		msg := fmt.Sprintf("CreateProduct: Error: %s is not authorized to %s", r.Header.Get("Role"), constants.READ)
		logger.ErrorLogger.Printf(msg)
		// w.Write([]byte("Not authorized to Create Product"))
		// return
		res := common.APIResponse{
			StatusCode: 401,
			Message:    "Not authorized to Create Product",
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}
	logger.InfoLogger.Printf("Create Product Request start")
	var reqProdParms common.Product
	rand := strconv.FormatInt(time.Now().Unix(), 10)
	reqProdParms.ProductId = "p" + rand
	reqProdParms.Name = r.FormValue("name")
	reqProdParms.Description = r.FormValue("description")
	reqProdParms.Price = r.FormValue("price")

	v := validator.New()
	err := v.Struct(reqProdParms)
	if err != nil {
		msg := fmt.Sprintf("validation Error: %s", err.Error())
		logger.ErrorLogger.Printf(msg)

		res := common.APIResponse{
			StatusCode: 400,
			Message:    err.Error(),
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// the Header and the size of the file
	file, handlerFile, err := r.FormFile("uploadFile")
	if err != nil {
		msg := fmt.Sprintf("CreateProduct: File upload error: %s", err.Error())
		logger.ErrorLogger.Printf(msg)
		res := common.APIResponse{
			StatusCode: 400,
			Message:    err.Error(),
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	defer file.Close()
	splitFilename := strings.Split(handlerFile.Filename, ".")
	fileExt := splitFilename[len(splitFilename)-1]

	if handlerFile.Size > 2000000 {
		msg := fmt.Sprintf("CreateProduct: File upload error: Uploaded file not be allowed more than 2 mb")
		logger.ErrorLogger.Printf(msg)

		res := common.APIResponse{
			StatusCode: 400,
			Message:    "Uploaded file not be allowed more than 2 mb",
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("CreateProduct: Read file error: %s", err.Error())
		logger.ErrorLogger.Printf(msg)

		res := common.APIResponse{
			StatusCode: 400,
			Message:    err.Error(),
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	dbConn := database.Connection()
	usersession := dbConn.Database("productcatalog").Collection("products")
	defer database.CloseClientDB(dbConn)

	reqProdParms.File = fileBytes
	reqProdParms.FileType = fileExt

	result, err := usersession.InsertOne(context.TODO(), &reqProdParms)
	if err != nil {
		mesg := fmt.Sprintf("Inseration failed with error %s", err.Error())
		logger.ErrorLogger.Printf("CreateProduct: " + mesg)
		fmt.Println(mesg)
		res := common.APIResponse{
			StatusCode: 500,
			Message:    mesg,
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	res := common.APIResponse{
		StatusCode: 201,
		Message:    "Product saved Sucessfully!!",
		Result:     result,
	}
	msg := fmt.Sprintf("CreateProduct: Product saved sucessfully: %s", result)
	logger.InfoLogger.Printf(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
	return
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {

	if !constants.IsAccess(r.Header.Get("Role"), constants.READ) {
		msg := fmt.Sprintf("GetProducts: Error: %s is not authorized to %s", r.Header.Get("Role"), constants.READ)
		logger.ErrorLogger.Printf(msg)
		// w.Write([]byte("Not authorized."))
		res := common.APIResponse{
			StatusCode: 401,
			Message:    "Not authorized to Get Product",
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}
	logger.InfoLogger.Printf("Get ALL Product Request start")
	var rProductArray []common.Product

	dbConn := database.Connection()
	usersession := dbConn.Database("productcatalog").Collection("products")
	defer database.CloseClientDB(dbConn)

	result, err := usersession.Find(context.TODO(), bson.M{})
	if err != nil {
		mesg := fmt.Sprintf("GetProducts : Error : %s", err.Error())
		logger.ErrorLogger.Printf(mesg)
		fmt.Println(mesg)
		res := common.APIResponse{
			StatusCode: 500,
			Message:    mesg,
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	if err = result.All(context.TODO(), &rProductArray); err != nil {
		mesg := fmt.Sprintf("Error While getting products %s", err.Error())
		msg := fmt.Sprintf("GetProducts: Error: %s", err.Error())
		logger.ErrorLogger.Printf(msg)

		fmt.Println(mesg)
		res := common.APIResponse{
			StatusCode: 500,
			Message:    mesg,
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	var respArray []common.RespProduct
	for _, rProd := range rProductArray {
		var prod common.RespProduct
		prod.ProductId = rProd.ProductId
		prod.Name = rProd.Name
		prod.Price = rProd.Price
		prod.Description = rProd.Description
		prod.File = common.GetHost(r) + "/download/" + rProd.ProductId
		respArray = append(respArray, prod)
	}

	res := common.APIResponse{
		StatusCode: 200,
		Message:    "Get all products Sucessfully!!",
		Result:     respArray,
	}
	msg := fmt.Sprintf("GetProducts: Get all products sucessfully")
	logger.InfoLogger.Printf(msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	if !constants.IsAccess(r.Header.Get("Role"), constants.READ) {
		msg := fmt.Sprintf("GetProduct: Error: %s is not authorized to %s", r.Header.Get("Role"), constants.READ)
		logger.ErrorLogger.Printf(msg)
		// w.Write([]byte("Not authorized to Create Product"))
		// return
		res := common.APIResponse{
			StatusCode: 401,
			Message:    "Not authorized to Get Product",
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}
	logger.InfoLogger.Printf("Get Product by ID Request start")
	var rProduct common.Product
	productID := mux.Vars(r)["id"]
	if productID == "" {
		msg := fmt.Sprintf("GetProductByID: Error: productID cannot be blank")
		logger.ErrorLogger.Printf(msg)

		mesg := fmt.Sprintf("productID is required")
		//logger.Error(mesg)
		fmt.Println(mesg)
		res := common.APIResponse{
			StatusCode: 500,
			Message:    mesg,
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	dbConn := database.Connection()
	usersession := dbConn.Database("productcatalog").Collection("products")
	defer database.CloseClientDB(dbConn)
	err := usersession.FindOne(context.TODO(), bson.M{"productid": productID}).Decode(&rProduct)
	if err != nil {
		msg := fmt.Sprintf("GetProductById: Error: %s", err.Error())
		logger.ErrorLogger.Printf(msg)

		mesg := fmt.Sprintf("Error While getting product %s", err.Error())
		//logger.Error(mesg)
		fmt.Println(mesg)
		res := common.APIResponse{
			StatusCode: 500,
			Message:    mesg,
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	var prod common.RespProduct
	prod.ProductId = rProduct.ProductId
	prod.Name = rProduct.Name
	prod.Price = rProduct.Price
	prod.Description = rProduct.Description
	prod.File = common.GetHost(r) + "/download/" + rProduct.ProductId

	res := common.APIResponse{
		StatusCode: 200,
		Message:    "Get product Sucessfully!!",
		Result:     prod,
	}
	logger.InfoLogger.Printf("GetProductByID: Get product sucessfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	productID := mux.Vars(r)["id"]
	fmt.Println(productID)

	msg := fmt.Sprintf("DowloadFile: Request to Download file %s", productID)
	logger.InfoLogger.Printf(msg)

	dbConn := database.Connection()
	usersession := dbConn.Database("productcatalog").Collection("products")
	defer database.CloseClientDB(dbConn)

	var rProduct common.Product
	_ = usersession.FindOne(context.TODO(), bson.M{"productid": productID}).Decode(&rProduct)
	if rProduct.ProductId != "" {
		rand := strconv.FormatInt(time.Now().Unix(), 10)
		fmt.Println("writing file")
		newFileName := "download" + rand + "." + rProduct.FileType
		w.Header().Set("Content-Disposition", "attachment; filename="+newFileName)
		_, err := io.Copy(w, strings.NewReader(string(rProduct.File)))
		if err != nil {
			msg := fmt.Sprintf("DowloadFile: unable to download file %s", err.Error())
			logger.ErrorLogger.Printf(msg)
			http.Error(w, "Remote server error", 503)
			return
		}
		return

	}

}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {

	if !constants.IsAccess(r.Header.Get("Role"), constants.READ) {
		msg := fmt.Sprintf("UpdateProduct: Error: %s is not authorized to %s", r.Header.Get("Role"), constants.READ)
		logger.ErrorLogger.Printf(msg)
		// w.Write([]byte("Not authorized to Create Product"))
		// return
		res := common.APIResponse{
			StatusCode: 401,
			Message:    "Not authorized to Update Product",
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}
	logger.InfoLogger.Printf("Update Product Request start")
	productID := mux.Vars(r)["id"]

	if productID == "" {
		msg := fmt.Sprintf("UpdateProduct: productID cannot be empty")
		logger.ErrorLogger.Printf(msg)
	}

	dbConn := database.Connection()
	usersession := dbConn.Database("productcatalog").Collection("products")
	defer database.CloseClientDB(dbConn)

	var rProduct common.Product
	_ = usersession.FindOne(context.TODO(), bson.M{"productid": productID}).Decode(&rProduct)

	if r.FormValue("productId") != "" {
		rProduct.ProductId = r.FormValue("productId")
	}
	if r.FormValue("name") != "" {
		rProduct.Name = r.FormValue("name")
	}
	if r.FormValue("description") != "" {
		rProduct.Description = r.FormValue("description")
	}
	if r.FormValue("price") != "" {
		rProduct.Price = r.FormValue("price")
	}

	v := validator.New()
	err := v.Struct(rProduct)
	if err != nil {
		res := common.APIResponse{
			StatusCode: 400,
			Message:    err.Error(),
			IsError:    true,
		}
		msg := fmt.Sprintf("UpdateProduct: Error: %s", err.Error())
		logger.ErrorLogger.Printf(msg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	// the Header and the size of the file
	file, handlerFile, err := r.FormFile("uploadFile")
	if file != nil {
		if err != nil {
			res := common.APIResponse{
				StatusCode: 400,
				Message:    err.Error(),
				IsError:    true,
			}
			msg := fmt.Sprintf("UpdateProduct: Upload new file Error: %s", err.Error())
			logger.ErrorLogger.Printf(msg)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		defer file.Close()
		splitFilename := strings.Split(handlerFile.Filename, ".")
		fileExt := splitFilename[len(splitFilename)-1]

		if handlerFile.Size > 2000000 {
			res := common.APIResponse{
				StatusCode: 400,
				Message:    "Uploaded file not be allowed more than 2 mb",
				IsError:    true,
			}
			msg := fmt.Sprintf("UpdateProduct: Error: Uploaded file not be allowed more than 2 mb")
			logger.ErrorLogger.Printf(msg)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			res := common.APIResponse{
				StatusCode: 400,
				Message:    err.Error(),
				IsError:    true,
			}
			msg := fmt.Sprintf("UpdateProduct: ReadFile Error: %s", err.Error())
			logger.ErrorLogger.Printf(msg)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		rProduct.File = fileBytes
		rProduct.FileType = fileExt
	}

	rProduct.Id = ""
	_, err = usersession.UpdateOne(context.TODO(), bson.M{"productid": productID}, bson.M{"$set": rProduct})
	if err != nil {
		res := common.APIResponse{
			StatusCode: 500,
			Message:    err.Error(),
			IsError:    true,
		}
		msg := fmt.Sprintf("UpdateProduct: Updating Error: %s", err.Error())
		logger.ErrorLogger.Printf(msg)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	var prod common.RespProduct
	prod.ProductId = rProduct.ProductId
	prod.Name = rProduct.Name
	prod.Price = rProduct.Price
	prod.Description = rProduct.Description
	prod.File = common.GetHost(r) + "/download/" + rProduct.ProductId

	res := common.APIResponse{
		StatusCode: 200,
		Message:    "Request updated sucessfully!!",
		IsError:    false,
		Result:     prod,
	}
	msg := fmt.Sprintf("UpdateProduct: update request sucessfull")
	logger.InfoLogger.Printf(msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if !constants.IsAccess(r.Header.Get("Role"), constants.READ) {
		msg := fmt.Sprintf("DeleteProduct: Error: %s is not authorized to %s", r.Header.Get("Role"), constants.READ)
		logger.ErrorLogger.Printf(msg)
		// w.Write([]byte("Not authorized to Create Product"))
		// return
		res := common.APIResponse{
			StatusCode: 401,
			Message:    "Not authorized to Delete Product",
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}

	productID := mux.Vars(r)["id"]

	if productID == "" {
		msg := fmt.Sprintf("DeleteProduct: Product ID cannot be blank")
		logger.ErrorLogger.Printf(msg)
	}
	fmt.Println(productID)
	var rProduct common.Product
	dbConn := database.Connection()
	usersession := dbConn.Database("productcatalog").Collection("products")
	defer database.CloseClientDB(dbConn)
	err := usersession.FindOne(context.TODO(), bson.M{"productid": productID}).Decode(&rProduct)
	if err != nil {
		mesg := fmt.Sprintf("Error While getting product detail %s", err.Error())
		msg := fmt.Sprintf("DeleteProduct:" + mesg)
		logger.ErrorLogger.Printf(msg)
		fmt.Println(mesg)
		res := common.APIResponse{
			StatusCode: 500,
			Message:    mesg,
			IsError:    true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	idPrimitive, err := primitive.ObjectIDFromHex(rProduct.Id)

	if err != nil {
		mesg := fmt.Sprintf("Error While deleting product %s", err.Error())
		res := common.APIResponse{
			StatusCode: 500,
			Message:    mesg,
			IsError:    true,
		}
		msg := fmt.Sprintf("DeleteProduct: " + mesg)
		logger.ErrorLogger.Printf(msg)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	_, err = usersession.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	if err != nil {
		mesg := fmt.Sprintf("Error While deleting product %s", err.Error())
		res := common.APIResponse{
			StatusCode: 500,
			Message:    mesg,
			IsError:    true,
		}
		msg := fmt.Sprintf("DeleteProduct: " + mesg)
		logger.ErrorLogger.Printf(msg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	res := common.APIResponse{
		StatusCode: 200,
		Message:    "Product Deleted successfully",
		IsError:    false,
	}
	msg := fmt.Sprintf("DeleteProduct: Product delete sucessfully")
	logger.InfoLogger.Printf(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

//check whether user is authorized or not
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer")
		if len(splitToken) != 2 {
			var err common.Error
			logger.ErrorLogger.Printf("Auth: No Token Found")
			err = helpers.SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}
		reqToken = strings.TrimSpace(splitToken[1])
		dbConn := database.Connection()
		var respUsertoken common.Token
		usersession := dbConn.Database("usercatalog").Collection("usertoken")
		_ = usersession.FindOne(context.TODO(), bson.M{"tokenstring": reqToken}).Decode(&respUsertoken)
		if respUsertoken.TokenString == "" {
			var err common.Error
			logger.ErrorLogger.Printf("Auth: Your Token is not vaild")
			err = helpers.SetError(err, "Your Token is not vaild.")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte(constants.SECRETKEY)

		token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.ErrorLogger.Printf("Auth: There was an error in parsing token.")
				return nil, fmt.Errorf("There was an error in parsing token.")
			}
			return mySigningKey, nil
		})
		if err != nil {
			var err common.Error
			logger.ErrorLogger.Printf("Auth: Your Token has been expired")
			err = helpers.SetError(err, "Your Token has been expired.")
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == constants.ADMIN {
				r.Header.Set("Role", constants.ADMIN)
				handler.ServeHTTP(w, r)
				return
			} else if claims["role"] == constants.USER {
				r.Header.Set("Role", constants.USER)
				handler.ServeHTTP(w, r)
				return
			}
		}
		var reserr common.Error
		logger.ErrorLogger.Printf("Auth: Not Authorized")
		reserr = helpers.SetError(reserr, "Not Authorized.")
		json.NewEncoder(w).Encode(err)
	}
}
